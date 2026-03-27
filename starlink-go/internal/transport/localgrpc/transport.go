package localgrpc

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const deviceHandleMethod = "/SpaceX.API.Device.Device/Handle"

// Transport implements the only supported network path for this package:
// direct local gRPC access to a Starlink terminal over LAN/VPN.
type Transport struct {
	mu      sync.RWMutex
	conn    *grpc.ClientConn
	address string
}

func Dial(ctx context.Context, address string) (*Transport, error) {
	if address == "" {
		return nil, fmt.Errorf("starlink local gRPC endpoint is required (host:port)")
	}

	t := &Transport{address: address}
	if err := t.Reconnect(ctx); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Transport) Host() string {
	return t.address
}

func (t *Transport) Reconnect(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, t.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("cannot reach starlink local gRPC endpoint %q via LAN/VPN: %w", t.address, err)
	}

	t.mu.Lock()
	old := t.conn
	t.conn = conn
	t.mu.Unlock()

	if old != nil {
		_ = old.Close()
	}

	return nil
}

func (t *Transport) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	t.mu.RLock()
	conn := t.conn
	t.mu.RUnlock()

	if conn == nil {
		if err := t.Reconnect(ctx); err != nil {
			return nil, err
		}
		t.mu.RLock()
		conn = t.conn
		t.mu.RUnlock()
	}

	resp := new(pb.Response)
	if err := conn.Invoke(ctx, deviceHandleMethod, req, resp); err != nil {
		return nil, fmt.Errorf("local gRPC Device.Handle failed for %q: %w", t.address, err)
	}
	return resp, nil
}

func (t *Transport) Close() error {
	t.mu.Lock()
	conn := t.conn
	t.conn = nil
	t.mu.Unlock()

	if conn == nil {
		return nil
	}
	return conn.Close()
}
