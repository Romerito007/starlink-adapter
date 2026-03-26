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
	cfg       Config
}

func NewGRPCClient(transport transport, cfg Config) *grpcClient {
	if cfg.Timeout <= 0 {
		cfg = defaultConfig()
	}
	if cfg.RetryMax <= 0 {
		cfg.RetryMax = 3
	}
	if cfg.BaseBackoff <= 0 {
		cfg.BaseBackoff = 200 * time.Millisecond
	}
	if cfg.Logger == nil {
		cfg.Logger = defaultConfig().Logger
	}

	return &grpcClient{transport: transport, cfg: cfg}
}

func Dial(ctx context.Context, address string) (*grpcClient, error) {
	return DialWithConfig(ctx, address, defaultConfig())
}

func DialWithConfig(ctx context.Context, address string, cfg Config) (*grpcClient, error) {
	if address == "" {
		address = DefaultDishAddress
	}

	t, err := localgrpc.Dial(ctx, address)
	if err != nil {
		return nil, normalizeError(err)
	}

	return NewGRPCClient(t, cfg), nil
}

var _ StarlinkClient = (*grpcClient)(nil)

func (c *grpcClient) Close() error {
	if c == nil || c.transport == nil {
		return nil
	}
	return c.transport.Close()
}

func (c *grpcClient) GetStatus(ctx context.Context) (*Status, error) {
	resp, err := c.sendWithRetry(ctx, "GetStatus", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_GetStatus{GetStatus: &pb.GetStatusRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	statusResp, ok := resp.Response.(*pb.Response_DishGetStatus)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}

	return mapStatus(statusResp.DishGetStatus), nil
}

func (c *grpcClient) GetStats(ctx context.Context) (*Stats, error) {
	resp, err := c.sendWithRetry(ctx, "GetStats", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_GetHistory{GetHistory: &pb.GetHistoryRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	historyResp, ok := resp.Response.(*pb.Response_DishGetHistory)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}

	return mapStats(historyResp.DishGetHistory), nil
}

func (c *grpcClient) GetLocation(ctx context.Context) (*Location, error) {
	resp, err := c.sendWithRetry(ctx, "GetLocation", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_GetLocation{GetLocation: &pb.GetLocationRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	locationResp, ok := resp.Response.(*pb.Response_GetLocation)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}

	return mapLocation(locationResp.GetLocation), nil
}

func (c *grpcClient) Reboot(ctx context.Context) error {
	resp, err := c.sendWithRetry(ctx, "Reboot", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_Reboot{Reboot: &pb.RebootRequest{}}}
	})
	if err != nil {
		return err
	}

	if _, ok := resp.Response.(*pb.Response_Reboot); !ok {
		return fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}

	return nil
}

func (c *grpcClient) sendWithRetry(ctx context.Context, operation string, reqFn func() *pb.Request) (*pb.Response, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("%w: transport is not configured", ErrUnavailable)
	}

	var lastErr error
	for attempt := 1; attempt <= c.cfg.RetryMax; attempt++ {
		started := time.Now()
		attemptCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
		resp, err := c.transport.Handle(attemptCtx, reqFn())
		cancel()
		latency := time.Since(started)

		if err == nil {
			c.cfg.Logger.Info("starlink operation success",
				"host", c.transport.Host(),
				"operation", operation,
				"attempt", attempt,
				"latency_ms", latency.Milliseconds(),
			)
			return resp, nil
		}

		nerr := normalizeError(err)
		lastErr = nerr
		c.cfg.Logger.Warn("starlink operation failed",
			"host", c.transport.Host(),
			"operation", operation,
			"attempt", attempt,
			"latency_ms", latency.Milliseconds(),
			"error", nerr.Error(),
		)

		if !isTransientError(nerr) || attempt == c.cfg.RetryMax {
			return nil, nerr
		}

		_ = c.transport.Reconnect(ctx)

		backoff := c.cfg.BaseBackoff * time.Duration(1<<(attempt-1))
		select {
		case <-ctx.Done():
			return nil, normalizeError(ctx.Err())
		case <-time.After(backoff):
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("%w: unknown transport failure", ErrUnavailable)
}
