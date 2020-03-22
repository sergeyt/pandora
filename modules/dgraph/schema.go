package dgraph

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/dgraph-io/dgo/v2/protos/api"
	log "github.com/sirupsen/logrus"
)

func InitSchema() {
	// TODO configurable path to schemas
	for i := 1; i <= 100; i++ {
		err := initSchema("./schema.txt")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	for i := 1; i <= 100; i++ {
		err := initGraphqlSchema("./schema.gql")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}

func initSchema(path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("read schema fail: %v", err)
		return err
	}

	ctx := context.Background()

	dg, close, err := NewClient(ctx)
	if err != nil {
		log.Errorf("cannot init dgraph schema: %v", err)
		// TODO retry after few seconds
		return err
	}
	defer close()

	ctx = WithAuthToken(ctx)

	err = dg.Alter(ctx, &api.Operation{
		Schema: string(schema),
	})
	if err != nil {
		log.Errorf("init %s fail: %v", path, err)
		return err
	}

	log.Infof("schema %s initialized", path)
	return nil
}

func initGraphqlSchema(path string) error {
	schema, err := os.Open(path)
	if err != nil {
		log.Errorf("read schema fail: %v", err)
		return err
	}

	// init graphql schema via HTTP call for now
	rc := NewRestClient()
	err = rc.PostData("/admin/schema", "text/plain", schema, nil)
	if err != nil {
		log.Errorf("init of graphql schema failed: %v", err)
		return err
	}

	log.Info("graphql schema initialized")
	return nil
}
