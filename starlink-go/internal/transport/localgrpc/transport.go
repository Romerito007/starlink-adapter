package localgrpc

import (
	"context"
	"fmt"

	"github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const deviceHandleMethod = "/SpaceX.API.Device.Device/Handle"

type Transport struct {
	conn *grpc.ClientConn
}

func Dial(ctx context.Context, address string) (*Transport, error) {
	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dial starlink local gRPC endpoint: %w", err)
	}

	return &Transport{conn: conn}, nil
}

func (t *Transport) Handle(ctx context.Context, req *device.Request) (*device.Response, error) {
	resp := new(device.Response)
	if err := t.conn.Invoke(ctx, deviceHandleMethod, req, resp); err != nil {
		return nil, fmt.Errorf("invoke Device.Handle: %w", err)
	}
	return resp, nil
}

func (t *Transport) Close() error {
	if t == nil || t.conn == nil {
		return nil
	}
	return t.conn.Close()
}
