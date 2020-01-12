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

func Dial(network, address string, timeout time.Duration) (net.Conn, error) {
	c, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &grpcCon{c}, nil
}

func NewClient() (*dgo.Dgraph, error) {
	// TODO configurable timeout
	d, err := grpc.Dial(config.DB.Addr, grpc.WithInsecure(), grpc.WithTimeout(30*time.Second))
	if err != nil {
		log.Errorf("grpc.Dial fail: %v", err)
		return nil, err
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), nil
}

// TODO incremental update of schema
func InitSchema() {
	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// TODO configurable path to schema
	schema, err := ioutil.ReadFile("./schema.txt")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	token := os.Getenv("DGRAPH_TOKEN")
	if len(token) > 0 {
		md := metadata.New(nil)
		md.Append("auth-token", token)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	err = c.Alter(ctx, &api.Operation{
		Schema: string(schema),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := NewClient()
		if err != nil {
			apiutil.SendError(w, err)
			return
		}

		tx := c.NewTxn()
		defer tx.Discard(r.Context())

		ctx := context.WithValue(r.Context(), "tx", tx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RequestTransaction(r *http.Request) *dgo.Txn {
	return r.Context().Value("tx").(*dgo.Txn)
}
