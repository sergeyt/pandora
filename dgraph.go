package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

func newDgraphClient() (*dgo.Dgraph, error) {
	d, err := grpc.Dial(config.DB.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), nil
}

func initSchema() {
	c, err := newDgraphClient()
	if err != nil {
		log.Fatal(err)
	}

	// TODO configurable path to schema
	schema, err := ioutil.ReadFile("./schema.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = c.Alter(context.Background(), &api.Operation{
		Schema: string(schema),
	})
	if err != nil {
		log.Fatal(err)
	}
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
