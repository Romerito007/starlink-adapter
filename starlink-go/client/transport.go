package client

import (
	"context"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

// transport is an internal protocol layer abstraction.
type transport interface {
	Host() string
	Reconnect(ctx context.Context) error
	Handle(ctx context.Context, req *pb.Request) (*pb.Response, error)
	Close() error
}
