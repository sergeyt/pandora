package dgraph

import (
	"context"
	"io/ioutil"

	"github.com/dgraph-io/dgo/v2/protos/api"
	log "github.com/sirupsen/logrus"
)

func InitSchema() {
	dg, close, err := NewClient()
	if err != nil {
		log.Errorf("cannot init dgraph schema: %v", err)
		// TODO retry after few seconds
		return
	}
	defer close()

	// TODO configurable path to schema
	schema, err := ioutil.ReadFile("./schema.txt")
	if err != nil {
		log.Fatal(err)
	}

	ctx := WithAuthToken(context.Background())

	err = dg.Alter(ctx, &api.Operation{
		Schema: string(schema),
	})
	if err != nil {
		log.Fatal(err)
	}
}
