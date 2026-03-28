package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Romerito007/starlink-adapter/starlink-go/internal/transport/localgrpc"
	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

// StarlinkClient is the minimal monitoring and basic-ops API.
type StarlinkClient interface {
	GetStatus(ctx context.Context) (*Status, error)
	GetStats(ctx context.Context) (*Stats, error)
	GetLocation(ctx context.Context) (*Location, error)
	GetConnectedClients(ctx context.Context) ([]ClientDevice, error)
	GetDhcpLeases(ctx context.Context) ([]DhcpLease, error)
	GetWifiConfig(ctx context.Context) (*WifiConfigSnapshot, error)
	GetNetworkInterfaces(ctx context.Context) ([]NetworkInterfaceSnapshot, error)
	GetRadioStats(ctx context.Context) ([]RadioStat, error)
	Reboot(ctx context.Context) error
	Close() error
}

type grpcClient struct {
	transport transport
	cfg       Config
	logger    *slog.Logger
}

var dialTransport = func(ctx context.Context, address string) (transport, error) {
	return localgrpc.Dial(ctx, address)
}

func newGRPCClient(transport transport, cfg Config) *grpcClient {
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultConfig().Timeout
	}

	return &grpcClient{
		transport: transport,
		cfg:       cfg,
		logger:    cfg.Logger,
	}
}

func NewClient(ctx context.Context, cfg Config) (StarlinkClient, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrUnavailable)
	}

	if cfg.Host == "" {
		cfg.Host = defaultConfig().Host
	}
	if cfg.Port <= 0 {
		cfg.Port = defaultConfig().Port
	}

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return dialWithConfig(ctx, address, cfg)
}

func dialWithConfig(ctx context.Context, address string, cfg Config) (*grpcClient, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultConfig().Timeout
	}

	dialCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	t, err := dialTransport(dialCtx, address)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf("initial dial to starlink local endpoint %q via LAN/VPN exceeded timeout %s: %w", address, cfg.Timeout, err)
		}
		return nil, normalizeError(err)
	}

	return newGRPCClient(t, cfg), nil
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

	switch statusResp := resp.Response.(type) {
	case *pb.Response_DishGetStatus:
		return mapStatus(statusResp.DishGetStatus), nil
	case *pb.Response_WifiGetStatus:
		return mapStatusFromWifi(statusResp.WifiGetStatus), nil
	default:
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}
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

func (c *grpcClient) GetConnectedClients(ctx context.Context) ([]ClientDevice, error) {
	resp, err := c.sendWithRetry(ctx, "GetConnectedClients", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_GetStatus{GetStatus: &pb.GetStatusRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	switch clientsResp := resp.Response.(type) {
	case *pb.Response_WifiGetStatus:
		if clientsResp.WifiGetStatus == nil {
			return []ClientDevice{}, nil
		}
		return mapConnectedClients(clientsResp.WifiGetStatus.GetClients()), nil
	case *pb.Response_DishGetStatus:
		// Fallback path for endpoints that do not expose wifi_get_status clients
		// in get_status responses.
		return c.getConnectedClientsViaWifiGetClients(ctx)
	default:
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}
}

func (c *grpcClient) GetDhcpLeases(ctx context.Context) ([]DhcpLease, error) {
	resp, err := c.sendWithRetry(ctx, "GetDhcpLeases", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_GetStatus{GetStatus: &pb.GetStatusRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	wifiResp, ok := resp.Response.(*pb.Response_WifiGetStatus)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}
	if wifiResp.WifiGetStatus == nil {
		return []DhcpLease{}, nil
	}

	return mapDhcpLeases(wifiResp.WifiGetStatus.GetDhcpServers()), nil
}

func (c *grpcClient) GetWifiConfig(ctx context.Context) (*WifiConfigSnapshot, error) {
	resp, err := c.sendWithRetry(ctx, "GetWifiConfig", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_WifiGetConfig{WifiGetConfig: &pb.WifiGetConfigRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	wifiResp, ok := resp.Response.(*pb.Response_WifiGetConfig)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}

	return mapWifiConfigSnapshot(wifiResp.WifiGetConfig), nil
}

func (c *grpcClient) GetNetworkInterfaces(ctx context.Context) ([]NetworkInterfaceSnapshot, error) {
	resp, err := c.sendWithRetry(ctx, "GetNetworkInterfaces", func() *pb.Request {
		return &pb.Request{
			Request: &pb.Request_GetNetworkInterfaces{
				GetNetworkInterfaces: &pb.GetNetworkInterfacesRequest{},
			},
		}
	})
	if err != nil {
		return nil, err
	}

	netResp, ok := resp.Response.(*pb.Response_GetNetworkInterfaces)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}
	if netResp.GetNetworkInterfaces == nil {
		return []NetworkInterfaceSnapshot{}, nil
	}

	return mapNetworkInterfaces(netResp.GetNetworkInterfaces.GetNetworkInterfaces()), nil
}

func (c *grpcClient) GetRadioStats(ctx context.Context) ([]RadioStat, error) {
	resp, err := c.sendWithRetry(ctx, "GetRadioStats", func() *pb.Request {
		return &pb.Request{
			Request: &pb.Request_GetRadioStats{
				GetRadioStats: &pb.GetRadioStatsRequest{},
			},
		}
	})
	if err != nil {
		return nil, err
	}

	radioResp, ok := resp.Response.(*pb.Response_GetRadioStats)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected response type %T", ErrUnsupported, resp.Response)
	}
	if radioResp.GetRadioStats == nil {
		return []RadioStat{}, nil
	}

	return mapRadioStats(radioResp.GetRadioStats.GetRadioStats()), nil
}

func (c *grpcClient) getConnectedClientsViaWifiGetClients(ctx context.Context) ([]ClientDevice, error) {
	resp, err := c.sendWithRetry(ctx, "GetConnectedClients.WifiGetClientsFallback", func() *pb.Request {
		return &pb.Request{Request: &pb.Request_WifiGetClients{WifiGetClients: &pb.WifiGetClientsRequest{}}}
	})
	if err != nil {
		return nil, err
	}

	clientsResp, ok := resp.Response.(*pb.Response_WifiGetClients)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected fallback response type %T", ErrUnsupported, resp.Response)
	}
	if clientsResp.WifiGetClients == nil {
		return []ClientDevice{}, nil
	}

	return mapConnectedClients(clientsResp.WifiGetClients.GetClients()), nil
}

func (c *grpcClient) sendWithRetry(ctx context.Context, operation string, reqFn func() *pb.Request) (*pb.Response, error) {
	if c == nil || c.transport == nil {
		return nil, fmt.Errorf("%w: transport is not configured", ErrUnavailable)
	}
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrUnavailable)
	}

	const retryMax = 3
	const baseBackoff = 200 * time.Millisecond

	var lastErr error
	for attempt := 1; attempt <= retryMax; attempt++ {
		started := time.Now()
		attemptCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
		resp, err := c.transport.Handle(attemptCtx, reqFn())
		cancel()
		latency := time.Since(started)

		if err == nil {
			if c.logger != nil {
				c.logger.Info("starlink operation success",
					"host", c.transport.Host(),
					"operation", operation,
					"attempt", attempt,
					"latency_ms", latency.Milliseconds(),
				)
			}
			return resp, nil
		}

		nerr := normalizeError(err)
		lastErr = nerr
		if c.logger != nil {
			c.logger.Warn("starlink operation failed",
				"host", c.transport.Host(),
				"operation", operation,
				"attempt", attempt,
				"latency_ms", latency.Milliseconds(),
				"error", nerr.Error(),
			)
		}

		if !isTransientError(nerr) || attempt == retryMax {
			return nil, nerr
		}

		_ = c.transport.Reconnect(ctx)

		backoff := baseBackoff * time.Duration(1<<(attempt-1))
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
