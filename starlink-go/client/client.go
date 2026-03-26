package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Eitol/starlink-client/starlink-go/internal/transport/localgrpc"
	"github.com/Eitol/starlink-client/starlink-go/proto/gen/spacex/api/device"
)

const DefaultDishAddress = "192.168.100.1:9200"

// Client is a deterministic Starlink local gRPC client.
type Client struct {
	transport Transport
}

func New(transport Transport) *Client {
	return &Client{transport: transport}
}

func Dial(ctx context.Context, address string) (*Client, error) {
	if address == "" {
		address = DefaultDishAddress
	}

	t, err := localgrpc.Dial(ctx, address)
	if err != nil {
		return nil, err
	}

	return New(t), nil
}

func (c *Client) Close() error {
	if c == nil || c.transport == nil {
		return nil
	}
	return c.transport.Close()
}

// Handle sends a raw Request to the local dish endpoint.
func (c *Client) Handle(ctx context.Context, req *device.Request) (*device.Response, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("transport is not configured")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	return c.transport.Handle(ctx, req)
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
