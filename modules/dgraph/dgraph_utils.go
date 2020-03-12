package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/gocontrib/auth"
	"github.com/sergeyt/pandora/modules/orderedjson"
	"github.com/sergeyt/pandora/modules/pagination"
	log "github.com/sirupsen/logrus"
)

func NodeLabel(resourceType string) string {
	return strings.Title(strings.ToLower(resourceType))
}

type ListResult struct {
	Items []map[string]interface{} `json:"items"`
	Total int64                    `json:"total"`
}

func ReadList(ctx context.Context, tx *dgo.Txn, label string, pg pagination.Pagination) (*ListResult, error) {
	query := fmt.Sprintf(`query items($offset: int, $limit: int) {
  items(func: has(%s), offset: $offset, first: $limit) {
    uid
    expand(_all_)
  }
  total(func: has(%s)) {
    count: count(uid)
  }
}`, label, label)
	resp, err := tx.QueryWithVars(ctx, query, map[string]string{
		"$offset": fmt.Sprintf("%d", pg.Offset),
		"$limit":  fmt.Sprintf("%d", pg.Limit),
	})
	if err != nil {
		log.Errorf("dgrapg.Txn.QueryWithVars fail: %v", err)
		return nil, err
	}

	var result struct {
		Items []map[string]interface{} `json:"items"`
		Total []struct {
			Count int64 `json:"count"`
		} `json:"total"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return nil, err
	}

	return &ListResult{
		Items: result.Items,
		Total: result.Total[0].Count,
	}, nil
}

func ReadNode(ctx context.Context, tx *dgo.Txn, id string) (map[string]interface{}, error) {
	query := `query node($id: string) {
  node(func: uid($id)) {
    expand(_all_)
  }
}`

	resp, err := tx.QueryWithVars(ctx, query, map[string]string{
		"$id": id,
	})
	if err != nil {
		log.Errorf("dgrapg.Txn.QueryWithVars fail: %v", err)
		return nil, err
	}

	var result struct {
		Results []map[string]interface{} `json:"node"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}

	d := result.Results[0]
	d["uid"] = id

	return d, nil
}

func ReadNodeTx(ctx context.Context, id string) (map[string]interface{}, error) {
	dg, close, err := NewClient()
	if err != nil {
		return nil, err
	}
	defer close()

	tx := dg.NewTxn()
	defer Discard(ctx, tx)

	return ReadNode(ctx, tx, id)
}

type Mutation struct {
	Input     orderedjson.Map
	NodeLabel string
	ID        string
	By        string
	NoCommit  bool
}

func Mutate(ctx context.Context, tx *dgo.Txn, m Mutation) ([]map[string]interface{}, error) {
	id := m.ID
	isNew := len(id) == 0
	now := time.Now()

	if len(m.By) == 0 {
		user := auth.GetContextUser(ctx)
		if user == nil {
			m.By = "system"
		} else {
			m.By = user.GetID()
		}
	}

	in := m.Input
	in["modified_at"] = now
	if len(m.By) > 0 && m.By != "system" {
		in["modified_by"] = wrapUID(m.By)
	}

	if isNew {
		in[m.NodeLabel] = ""
		in["dgraph.type"] = m.NodeLabel
		in["created_at"] = now
		if len(m.By) > 0 && m.By != "system" {
			in["created_by"] = wrapUID(m.By)
		}
	} else {
		in["uid"] = id
	}

	data, err := in.ToJSON("uid", m.NodeLabel)
	if err != nil {
		log.Errorf("OrderedJSON.ToJSON fail: %v", err)
		return nil, err
	}

	resp, err := tx.Mutate(ctx, &api.Mutation{
		SetJson: data,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		return nil, err
	}

	var results []map[string]interface{}

	if isNew {
		results = make([]map[string]interface{}, len(resp.Uids))
		i := 0
		for _, uid := range resp.Uids {
			result, err := ReadNode(ctx, tx, uid)
			if err != nil {
				return nil, err
			}
			results[i] = result
			i = i + 1
			if len(results) == 1 {
				id = uid
			}
		}
	} else {
		result, err := ReadNode(ctx, tx, id)
		if err != nil {
			return nil, err
		}
		results = []map[string]interface{}{result}
	}

	if !m.NoCommit {
		err = tx.Commit(ctx)
		if err != nil {
			log.Errorf("dgraph.Txn.Commit fail: %v", err)
			return nil, err
		}
	}

	return results, nil
}

func wrapUID(uid string) map[string]string {
	return map[string]string{
		"uid": uid,
	}
}

func DeleteNode(ctx context.Context, tx *dgo.Txn, id string) (*api.Response, error) {
	// TODO unlink all connected nodes
	// del2 := fmt.Sprintf("* * <%s> .\n", id)
	del := fmt.Sprintf("<%s> * * .\n", id)
	resp, err := tx.Mutate(ctx, &api.Mutation{
		DelNquads: []byte(del),
		CommitNow: true,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		return nil, err
	}
	return resp, err
}

func IsUID(s string) bool {
	if strings.HasPrefix(s, "0x") {
		_, err := strconv.ParseInt(s[2:], 16, 64)
		if err == nil {
			return true
		}
	}
	return false
}

func Discard(ctx context.Context, txn *dgo.Txn) {
	if err := txn.Discard(ctx); err != nil {
		log.Errorf("Error while discarding transaction: %v", err)
	}
}
