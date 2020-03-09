package cloudstore

import (
	"context"
	"encoding/json"
	"fmt"

	dgo "github.com/dgraph-io/dgo/v2"

	_ "github.com/graymeta/stow/google"
	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"
)

type FileInfo struct {
	ID        string `json:"uid"`
	Path      string `json:"path"`
	URL       string `json:"url"`
	MediaType string `json:"content_type"`
}

func FindFileImpl(ctx context.Context, tx *dgo.Txn, id string) (*FileInfo, error) {
	filter := "eq(path, $id)"
	if dgraph.IsUID(id) {
		filter = "uid($id)"
	}

	query := fmt.Sprintf(`query file($id: string) {
		files(func: %s) {
			uid
			path
		}
	  }`, filter)
	resp, err := tx.QueryWithVars(ctx, query, map[string]string{
		"$id": id,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		return nil, err
	}

	var result struct {
		Files []FileInfo `json:"files"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return nil, err
	}

	if len(result.Files) == 0 {
		return nil, nil
	}

	if len(result.Files) > 1 {
		return nil, fmt.Errorf("inconsistent db state: found multiple file nodes")
	}

	file := result.Files[0]
	return &file, nil
}

func FindFile(ctx context.Context, id string) (*FileInfo, error) {
	dg, close, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}
	defer close()

	tx := dg.NewTxn()
	defer dgraph.Discard(ctx, tx)

	return FindFileImpl(ctx, tx, id)
}

func AddFile(ctx context.Context, tx *dgo.Txn, file *FileInfo) (map[string]interface{}, error) {
	in := make(utils.OrderedJSON)
	id := file.ID
	if len(id) > 0 {
		in["uid"] = id
	}

	in["content_type"] = file.MediaType

	if len(file.URL) > 0 {
		in["url"] = file.URL
	}

	if len(file.Path) > 0 {
		in["path"] = file.Path
	}

	if len(file.URL) == 0 && len(file.Path) > 0 {
		baseURL := config.ServerURL()
		in["url"] = fmt.Sprintf("%s/api/file/%s", baseURL, file.Path)
	}

	dispose1 := noop
	dispose2 := noop
	if tx == nil {
		dg, close, err := dgraph.NewClient()
		if err != nil {
			return nil, err
		}

		dispose1 = close
		tx = dg.NewTxn()
		dispose2 = func() {
			dgraph.Discard(ctx, tx)
		}
	}

	defer dispose1()
	defer dispose2()

	results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: dgraph.NodeLabel("file"),
		ID:        id,
	})
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected mutation results: %v", results)
	}

	return results[0], nil
}
