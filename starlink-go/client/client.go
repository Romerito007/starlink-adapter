package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DefaultDishAddress = "192.168.100.1:9200"

// Client is a minimal local gRPC client for the Starlink dish API.
type Client struct {
	conn *grpc.ClientConn
}

func Dial(ctx context.Context, address string) (*Client, error) {
	if address == "" {
		address = DefaultDishAddress
	}

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dial starlink gRPC endpoint: %w", err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Handle sends a raw Request to the local dish endpoint.
func (c *Client) Handle(ctx context.Context, req *device.Request) (*device.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	resp := new(device.Response)
	if err := c.conn.Invoke(ctx, "/SpaceX.API.Device.Device/Handle", req, resp); err != nil {
		return nil, fmt.Errorf("invoke Device.Handle: %w", err)
	}

	return resp, nil
}

// GetStatus requests current dish status from the local gRPC API.
func (c *Client) GetStatus(ctx context.Context) (*device.DishGetStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.Handle(ctx, &device.Request{
		Request: &device.Request_GetStatus{
			GetStatus: &device.GetStatusRequest{},
		},
	})
	if err != nil {
		return nil, err
	}

	statusResp, ok := resp.Response.(*device.Response_DishGetStatus)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", resp.Response)
	}

	return statusResp.DishGetStatus, nil
}
