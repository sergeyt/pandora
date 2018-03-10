package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
	"google.golang.org/grpc"
)

func newDgraphClient() (*client.Dgraph, error) {
	d, err := grpc.Dial(config.DB.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return client.NewDgraphClient(
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
