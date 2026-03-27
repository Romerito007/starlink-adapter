package client

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

type fakeTransport struct {
	handleFn        func(ctx context.Context, req *pb.Request) (*pb.Response, error)
	handleCallCount int
	lastRequest     *pb.Request
}

func (f *fakeTransport) Host() string                        { return "fake-host" }
func (f *fakeTransport) Reconnect(ctx context.Context) error { return nil }
func (f *fakeTransport) Close() error                        { return nil }

func (f *fakeTransport) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	f.handleCallCount++
	f.lastRequest = req
	if f.handleFn != nil {
		return f.handleFn(ctx, req)
	}
	return nil, nil
}

func newTestClient(t *testing.T, tr transport) *grpcClient {
	t.Helper()
	return newGRPCClient(tr, Config{Timeout: time.Second})
}

func TestGetConnectedClients_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetClients{
					WifiGetClients: &pb.WifiGetClientsResponse{
						Clients: []*pb.WifiClient{
							{
								ClientId:           2,
								MacAddress:         "BB:00:00:00:00:02",
								IpAddress:          "192.168.1.20",
								Ipv6Addresses:      []string{"2001::b", "2001::a"},
								Name:               "phone-b",
								GivenName:          "Phone B",
								Domain:             "lan",
								Iface:              pb.WifiClient_RF_5GHZ,
								IfaceName:          "wlan1",
								Role:               pb.WifiClient_REPEATER,
								UpstreamMacAddress: "AA:AA:AA:AA:AA:AA",
								AssociatedTimeS:    22,
								SignalStrength:     -58.5,
								Snr:                30.5,
								ChannelWidth:       40,
								ModeStr:            "11ac",
								Blocked:            true,
								DhcpLeaseActive:    true,
								DhcpLeaseRenewed:   false,
								NoDataIdleS:        120,
								RxStats: &pb.WifiClient_RxStats{
									RateMbps:         180,
									RateMbpsLast_15S: 150.5,
								},
								TxStats: &pb.WifiClient_TxStats{
									RateMbps:         90,
									RateMbpsLast_15S: 75.5,
								},
							},
							{
								ClientId:           1,
								MacAddress:         "AA:00:00:00:00:01",
								IpAddress:          "192.168.1.10",
								Ipv6Addresses:      []string{"2001::2", "2001::1"},
								Name:               "phone-a",
								GivenName:          "Phone A",
								Domain:             "lan",
								Iface:              pb.WifiClient_ETH,
								IfaceName:          "eth1",
								Role:               pb.WifiClient_CLIENT,
								UpstreamMacAddress: "CC:CC:CC:CC:CC:CC",
								AssociatedTimeS:    11,
								SignalStrength:     -48.5,
								Snr:                40.5,
								ChannelWidth:       20,
								ModeStr:            "eth",
								Blocked:            false,
								DhcpLeaseActive:    true,
								DhcpLeaseRenewed:   true,
								NoDataIdleS:        10,
								RxStats: &pb.WifiClient_RxStats{
									RateMbps:         950,
									RateMbpsLast_15S: 900.5,
								},
								TxStats: &pb.WifiClient_TxStats{
									RateMbps:         400,
									RateMbpsLast_15S: 350.5,
								},
							},
						},
					},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetConnectedClients(context.Background())
	if err != nil {
		t.Fatalf("GetConnectedClients() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetWifiGetClients() == nil {
		t.Fatalf("expected wifi_get_clients request, got: %#v", tr.lastRequest)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 clients, got %d", len(got))
	}

	// Output is deterministic by interface + mac address, so ETH entry comes first.
	first := got[0]
	if first.MacAddress != "AA:00:00:00:00:01" {
		t.Fatalf("unexpected MacAddress: %q", first.MacAddress)
	}
	if first.IpAddress != "192.168.1.10" {
		t.Fatalf("unexpected IpAddress: %q", first.IpAddress)
	}
	if first.Name != "phone-a" || first.GivenName != "Phone A" || first.Domain != "lan" {
		t.Fatalf("unexpected name fields: %+v", first)
	}
	if first.Interface != "ETH" {
		t.Fatalf("unexpected Interface: %q", first.Interface)
	}
	if first.InterfaceName != "eth1" || first.Role != "CLIENT" {
		t.Fatalf("unexpected interface/role mapping: %+v", first)
	}
	if first.AssociatedTimeSeconds != 11 {
		t.Fatalf("unexpected AssociatedTimeSeconds: %d", first.AssociatedTimeSeconds)
	}
	if first.SignalStrength != -48.5 {
		t.Fatalf("unexpected SignalStrength: %v", first.SignalStrength)
	}
	if first.Snr != 40.5 || first.ChannelWidth != 20 || first.Mode != "eth" {
		t.Fatalf("unexpected radio mode mapping: %+v", first)
	}
	if !first.DhcpLeaseActive || !first.DhcpLeaseRenewed || first.Blocked {
		t.Fatalf("unexpected dhcp/blocked mapping: %+v", first)
	}
	if first.NoDataIdleSeconds != 10 {
		t.Fatalf("unexpected NoDataIdleSeconds: %d", first.NoDataIdleSeconds)
	}
	if first.RxRateMbps != 950 || first.TxRateMbps != 400 {
		t.Fatalf("unexpected rx/tx rate mapping: %+v", first)
	}
	if first.RxRateMbpsLast15s != 900.5 || first.TxRateMbpsLast15s != 350.5 {
		t.Fatalf("unexpected last15s rate mapping: %+v", first)
	}
	if len(first.Ipv6Addresses) != 2 || first.Ipv6Addresses[0] != "2001::1" || first.Ipv6Addresses[1] != "2001::2" {
		t.Fatalf("unexpected sorted Ipv6Addresses: %#v", first.Ipv6Addresses)
	}
}

func TestGetConnectedClients_EmptyResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetClients{
					WifiGetClients: nil,
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetConnectedClients(context.Background())
	if err != nil {
		t.Fatalf("GetConnectedClients() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestGetConnectedClients_OptionalFieldsMissing(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetClients{
					WifiGetClients: &pb.WifiGetClientsResponse{
						Clients: []*pb.WifiClient{
							{},
						},
					},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetConnectedClients(context.Background())
	if err != nil {
		t.Fatalf("GetConnectedClients() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 client, got %d", len(got))
	}

	entry := got[0]
	if entry.MacAddress != "" || entry.IpAddress != "" || entry.Name != "" || entry.GivenName != "" || entry.Domain != "" {
		t.Fatalf("expected empty optional strings, got %+v", entry)
	}
	if entry.Interface != "UNKNOWN" {
		t.Fatalf("expected unknown enum string, got interface=%q", entry.Interface)
	}
	if entry.Role != "ROLE_UNKNOWN" {
		t.Fatalf("expected unknown role enum string, got role=%q", entry.Role)
	}
	if entry.Ipv6Addresses == nil {
		t.Fatalf("expected non-nil Ipv6Addresses slice")
	}
	if entry.RxRateMbps != 0 || entry.TxRateMbps != 0 || entry.RxRateMbpsLast15s != 0 || entry.TxRateMbpsLast15s != 0 {
		t.Fatalf("expected zeroed throughput metrics, got %+v", entry)
	}
}

func TestGetConnectedClients_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetConnectedClients(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetConnectedClients_UnexpectedResponseType(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_DishGetStatus{DishGetStatus: &pb.DishGetStatusResponse{}},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetConnectedClients(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
