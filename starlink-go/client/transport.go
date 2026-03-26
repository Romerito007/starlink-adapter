package client

import (
	"context"

	pb "github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
)

// transport is an internal protocol layer abstraction.
type transport interface {
	Handle(ctx context.Context, req *pb.Request) (*pb.Response, error)
	Close() error
}
