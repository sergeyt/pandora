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

// NewClient creates new dgraph client
func NewClient(ctx context.Context) (*dgo.Dgraph, context.CancelFunc, error) {
	// TODO configurable timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	conn, err := grpc.DialContext(ctx, config.DGraph.RpcURL, grpc.WithInsecure())
	if err != nil {
		log.Errorf("grpc.Dial fail: %v", err)
		return nil, nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	close := func() {
		cancel()
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
