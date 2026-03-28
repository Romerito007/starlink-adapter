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
	requests        []*pb.Request
}

func (f *fakeTransport) Host() string                        { return "fake-host" }
func (f *fakeTransport) Reconnect(ctx context.Context) error { return nil }
func (f *fakeTransport) Close() error                        { return nil }

func (f *fakeTransport) Handle(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	f.handleCallCount++
	f.lastRequest = req
	f.requests = append(f.requests, req)
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
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{
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
								HopsFromController: 3,
								RxStats: &pb.WifiClient_RxStats{
									Bytes:            1100,
									Nss:              1,
									Mcs:              4,
									Bandwidth:        40,
									GuardNs:          800,
									PhyMode:          2,
									RateMbps:         180,
									RateMbpsLast_15S: 150.5,
								},
								TxStats: &pb.WifiClient_TxStats{
									Bytes:            2200,
									Nss:              2,
									Mcs:              5,
									Bandwidth:        80,
									GuardNs:          400,
									PhyMode:          3,
									RateMbps:         90,
									RateMbpsLast_30S: 65.2,
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
								HopsFromController: 1,
								RxStats: &pb.WifiClient_RxStats{
									Bytes:            1010,
									Nss:              4,
									Mcs:              11,
									Bandwidth:        160,
									GuardNs:          800,
									PhyMode:          9,
									RateMbps:         950,
									RateMbpsLast_15S: 900.5,
								},
								TxStats: &pb.WifiClient_TxStats{
									Bytes:            2020,
									Nss:              2,
									Mcs:              8,
									Bandwidth:        80,
									GuardNs:          400,
									PhyMode:          7,
									RateMbps:         400,
									RateMbpsLast_30S: 300.5,
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

	if tr.lastRequest == nil || tr.lastRequest.GetGetStatus() == nil {
		t.Fatalf("expected get_status request, got: %#v", tr.lastRequest)
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
	if first.Interface != "eth" {
		t.Fatalf("unexpected Interface: %q", first.Interface)
	}
	if first.InterfaceName != "eth1" || first.Role != "client" {
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
	if first.UpstreamMacAddress != "CC:CC:CC:CC:CC:CC" || first.HopsFromController != 1 || first.ClientID != 1 {
		t.Fatalf("unexpected upstream/client mapping: %+v", first)
	}
	if first.RxRateMbps != 950 || first.TxRateMbps != 400 {
		t.Fatalf("unexpected rx/tx rate mapping: %+v", first)
	}
	if first.RxRateMbpsLast15s != 900.5 || first.TxRateMbpsLast15s != 350.5 {
		t.Fatalf("unexpected last15s rate mapping: %+v", first)
	}
	if first.RxBytes != 1010 || first.TxBytes != 2020 || !first.RxStatsValid || !first.TxStatsValid {
		t.Fatalf("unexpected stats bytes/valid mapping: %+v", first)
	}
	if first.RxNss != 4 || first.TxNss != 2 || first.RxMcs != 11 || first.TxMcs != 8 {
		t.Fatalf("unexpected nss/mcs mapping: %+v", first)
	}
	if first.RxBandwidth != 160 || first.TxBandwidth != 80 || first.RxGuardNs != 800 || first.TxGuardNs != 400 {
		t.Fatalf("unexpected bandwidth/guard mapping: %+v", first)
	}
	if first.RxPhyMode != 9 || first.TxPhyMode != 7 {
		t.Fatalf("unexpected phy mode mapping: %+v", first)
	}
	if first.TxRateMbpsLast30s != 300.5 {
		t.Fatalf("unexpected tx 30s mapping: %+v", first)
	}
	if first.CaptiveClientID != 0 || first.UploadMb != 0 || first.DownloadMb != 0 ||
		first.DhcpLeaseFound || first.SecondsUntilDhcpLeaseExpires != 0 {
		t.Fatalf("unexpected defaults for unsupported fields: %+v", first)
	}
	if first.RxRateMbpsLast1mAvg != 900.5 {
		t.Fatalf("unexpected rx 1m avg mapping: %+v", first)
	}
	if len(first.Ipv6Addresses) != 2 || first.Ipv6Addresses[0] != "2001::1" || first.Ipv6Addresses[1] != "2001::2" {
		t.Fatalf("unexpected sorted Ipv6Addresses: %#v", first.Ipv6Addresses)
	}
}

func TestGetConnectedClients_EmptyResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: nil,
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
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{
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
	if entry.Interface != "unknown" {
		t.Fatalf("expected unknown enum string, got interface=%q", entry.Interface)
	}
	if entry.Role != "unknown" {
		t.Fatalf("expected unknown role enum string, got role=%q", entry.Role)
	}
	if entry.Ipv6Addresses == nil {
		t.Fatalf("expected non-nil Ipv6Addresses slice")
	}
	if entry.RxRateMbps != 0 || entry.TxRateMbps != 0 || entry.RxRateMbpsLast15s != 0 || entry.TxRateMbpsLast15s != 0 {
		t.Fatalf("expected zeroed throughput metrics, got %+v", entry)
	}
	if entry.RxStatsValid || entry.TxStatsValid || entry.RxBytes != 0 || entry.TxBytes != 0 {
		t.Fatalf("expected zeroed rx/tx stats metadata, got %+v", entry)
	}
}

func TestGetConnectedClients_FallbacksToWifiGetClients(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			if req.GetGetStatus() != nil {
				return &pb.Response{
					Response: &pb.Response_DishGetStatus{DishGetStatus: &pb.DishGetStatusResponse{}},
				}, nil
			}
			return &pb.Response{
				Response: &pb.Response_WifiGetClients{
					WifiGetClients: &pb.WifiGetClientsResponse{
						Clients: []*pb.WifiClient{
							{MacAddress: "AA:BB:CC:DD:EE:FF"},
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
	if len(got) != 1 || got[0].MacAddress != "AA:BB:CC:DD:EE:FF" {
		t.Fatalf("unexpected fallback mapping: %+v", got)
	}
	if tr.handleCallCount != 2 {
		t.Fatalf("expected 2 transport calls with fallback, got %d", tr.handleCallCount)
	}
	if tr.requests[0].GetGetStatus() == nil || tr.requests[1].GetWifiGetClients() == nil {
		t.Fatalf("unexpected request sequence for fallback: %#v", tr.requests)
	}
}

func TestGetConnectedClients_RateFallbacksWhenRecentWindowsMissing(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{
						Clients: []*pb.WifiClient{
							{
								MacAddress: "AA:BB:CC:DD:EE:FF",
								RxStats: &pb.WifiClient_RxStats{
									RateMbps: 72,
								},
								TxStats: &pb.WifiClient_TxStats{
									RateMbps:         72,
									RateMbpsLast_30S: 10.466666,
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
	if len(got) != 1 {
		t.Fatalf("expected 1 client, got %d", len(got))
	}

	client := got[0]
	if client.RxRateMbps != 72 || client.TxRateMbps != 72 {
		t.Fatalf("unexpected current rate mapping: %+v", client)
	}
	if client.RxRateMbpsLast15s != 72 {
		t.Fatalf("expected rx last15 fallback to current rate, got %+v", client)
	}
	if client.RxRateMbpsLast1mAvg != 72 {
		t.Fatalf("expected rx last1m avg fallback, got %+v", client)
	}
	if client.TxRateMbpsLast30s != 10.466666 {
		t.Fatalf("expected tx last30 direct mapping, got %+v", client)
	}
	if client.TxRateMbpsLast15s != 10.466666 {
		t.Fatalf("expected tx last15 fallback to last30, got %+v", client)
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
				Response: &pb.Response_Reboot{Reboot: &pb.RebootResponse{}},
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
