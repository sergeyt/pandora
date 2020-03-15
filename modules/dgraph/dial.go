package dgraph

import (
	"net"
	"time"

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
