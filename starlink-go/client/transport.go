package client

import (
	"context"

	"github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
)

// Transport abstracts how requests are sent to the local dish endpoint.
type Transport interface {
	Handle(ctx context.Context, req *device.Request) (*device.Response, error)
	Close() error
}
