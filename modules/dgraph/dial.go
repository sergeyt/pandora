package dgraph

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
)

type grpcCon struct {
	conn   *grpc.ClientConn
	cancel context.CancelFunc
}

func (c *grpcCon) Close() error {
	c.cancel()
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	c, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &grpcCon{c, cancel}, nil
}
