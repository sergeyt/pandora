package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/sergeyt/pandora/modules/config"
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
