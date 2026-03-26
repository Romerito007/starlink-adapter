package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Eitol/starlink-client/starlink-go/internal/transport/localgrpc"
	pb "github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
)

const DefaultDishAddress = "192.168.100.1:9200"

// StarlinkClient is the minimal monitoring and basic-ops API.
type StarlinkClient interface {
	GetStatus(ctx context.Context) (*Status, error)
	GetStats(ctx context.Context) (*Stats, error)
	GetLocation(ctx context.Context) (*Location, error)
	Reboot(ctx context.Context) error
}

type grpcClient struct {
	transport transport
}

func NewGRPCClient(transport transport) *grpcClient {
	return &grpcClient{transport: transport}
}

func Dial(ctx context.Context, address string) (*grpcClient, error) {
	if address == "" {
		address = DefaultDishAddress
	}

	t, err := localgrpc.Dial(ctx, address)
	if err != nil {
		return nil, err
	}

	return NewGRPCClient(t), nil
}

var _ StarlinkClient = (*grpcClient)(nil)

func (c *grpcClient) Close() error {
	if c == nil || c.transport == nil {
		return nil
	}
	return c.transport.Close()
}

func (c *grpcClient) GetStatus(ctx context.Context) (*Status, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.send(ctx, &pb.Request{Request: &pb.Request_GetStatus{GetStatus: &pb.GetStatusRequest{}}})
	if err != nil {
		return nil, err
	}

	statusResp, ok := resp.Response.(*pb.Response_DishGetStatus)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", resp.Response)
	}

	return mapStatus(statusResp.DishGetStatus), nil
}

func (c *grpcClient) GetStats(ctx context.Context) (*Stats, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.send(ctx, &pb.Request{Request: &pb.Request_GetHistory{GetHistory: &pb.GetHistoryRequest{}}})
	if err != nil {
		return nil, err
	}

	historyResp, ok := resp.Response.(*pb.Response_DishGetHistory)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", resp.Response)
	}

	return mapStats(historyResp.DishGetHistory), nil
}

func (c *grpcClient) GetLocation(ctx context.Context) (*Location, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.send(ctx, &pb.Request{Request: &pb.Request_GetLocation{GetLocation: &pb.GetLocationRequest{}}})
	if err != nil {
		return nil, err
	}

	locationResp, ok := resp.Response.(*pb.Response_GetLocation)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", resp.Response)
	}

	return mapLocation(locationResp.GetLocation), nil
}

func (c *grpcClient) Reboot(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.send(ctx, &pb.Request{Request: &pb.Request_Reboot{Reboot: &pb.RebootRequest{}}})
	if err != nil {
		return err
	}

	if _, ok := resp.Response.(*pb.Response_Reboot); !ok {
		return fmt.Errorf("unexpected response type %T", resp.Response)
	}

	return nil
}

func (c *grpcClient) send(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("transport is not configured")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	return c.transport.Handle(ctx, req)
}
