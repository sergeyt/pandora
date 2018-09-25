package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo"
)

func readNode(ctx context.Context, tx *dgo.Txn, id string) (map[string]interface{}, error) {
	query := fmt.Sprintf(`{
  node(func: uid(%s)) {
		expand(_all_) {
			expand(_all_)
		}
  }
}`, id)

	resp, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []map[string]interface{} `json:"node"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)

	if err != nil {
		return nil, err
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}

	d := result.Results[0]
	d["uid"] = id

	return d, nil
}

func queryKeys(ctx context.Context, tx *dgo.Txn, id string) ([]string, error) {
	query := fmt.Sprintf(`{
  keys(func: uid(%s)) {
    _predicate_
  }
}`, id)

	resp, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []struct {
			Keys []string `json:"_predicate_"`
		} `json:"keys"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return result.Results[0].Keys, nil
}
