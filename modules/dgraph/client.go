package dgraph

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/sergeyt/pandora/modules/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type CloseFunc func()

// NewClient creates new dgraph client
func NewClient() (*dgo.Dgraph, CloseFunc, error) {
	// TODO configurable timeout
	conn, err := grpc.Dial(config.DB.Addr, grpc.WithInsecure(), grpc.WithTimeout(30*time.Second))
	if err != nil {
		log.Errorf("grpc.Dial fail: %v", err)
		return nil, nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	close := func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Error while closing connection: %v", err)
		}
	}

	return dg, close, nil
}

func WithAuthToken(ctx context.Context) context.Context {
	token := os.Getenv("DGRAPH_TOKEN")
	if len(token) > 0 {
		md := metadata.New(nil)
		md.Append("auth-token", token)
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
