package dgraph

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc/metadata"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcCon struct {
	conn *grpc.ClientConn
}

func (c *grpcCon) Close() error {
	return c.conn.Close()
}

func (c *grpcCon) LocalAddr() net.Addr {
	panic("not implemented")
}

func (c *grpcCon) RemoteAddr() net.Addr {
	panic("not implemented")
}

func (c *grpcCon) Read(b []byte) (n int, err error) {
	panic("not implemented")
}

func (c *grpcCon) Write(b []byte) (n int, err error) {
	panic("not implemented")
}

func (c *grpcCon) SetDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *grpcCon) SetReadDeadline(t time.Time) error {
	panic("not implemented")
}

func (c *grpcCon) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
}

// Dial to gRPC service. It is used in health checks
func Dial(network, address string, timeout time.Duration) (net.Conn, error) {
	c, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &grpcCon{c}, nil
}

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

// TODO incremental update of schema
func InitSchema() {
	dg, close, err := NewClient()
	if err != nil {
		log.Fatal(err)
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

func TransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dg, close, err := NewClient()
		if err != nil {
			apiutil.SendError(w, err)
			return
		}
		defer close()

		ctx := r.Context()
		tx := dg.NewTxn()
		defer Discard(ctx, tx)

		ctx = context.WithValue(ctx, "tx", tx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RequestTransaction(r *http.Request) *dgo.Txn {
	return r.Context().Value("tx").(*dgo.Txn)
}
