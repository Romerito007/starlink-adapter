package grpc

import (
	"context"
	"errors"
)

type DialOption struct{}

type ClientConn struct{}

func DialContext(_ context.Context, _ string, _ ...DialOption) (*ClientConn, error) {
	return &ClientConn{}, nil
}

func WithTransportCredentials(_ any) DialOption {
	return DialOption{}
}

func WithBlock() DialOption {
	return DialOption{}
}

func (c *ClientConn) Invoke(_ context.Context, _ string, _ any, _ any, _ ...any) error {
	return errors.New("grpc transport unavailable in stub build")
}

func (c *ClientConn) Close() error {
	return nil
}
