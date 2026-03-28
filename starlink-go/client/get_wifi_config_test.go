package client

import (
	"context"
	"errors"
	"testing"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetWifiConfig_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetConfig{
					WifiGetConfig: &pb.WifiGetConfigResponse{
						WifiConfig: &pb.WifiConfig{
							CountryCode:     "BR",
							SetupComplete:   true,
							MacWan:          "AA:AA:AA:AA:AA:AA",
							MacLan:          "BB:BB:BB:BB:BB:BB",
							BootCount:       12,
							Incarnation:     12345,
							WanHostDscpMark: 46,
							Networks: []*pb.WifiConfig_Network{
								{
									Ipv4:                 "192.168.2.1",
									Domain:               "guest",
									Dhcpv4Start:          100,
									Dhcpv4LeaseDurationS: 7200,
									Vlan:                 20,
									BasicServiceSets: []*pb.WifiConfig_BasicServiceSet{
										{
											Bssid:     "22:22:22:22:22:22",
											Ssid:      "Guest-5G",
											Band:      pb.WifiConfig_RF_5GHZ,
											IfaceName: "wlan1",
										},
									},
								},
								{
									Ipv4:                 "192.168.1.1",
									Domain:               "lan",
									Dhcpv4Start:          10,
									Dhcpv4LeaseDurationS: 3600,
									Vlan:                 10,
									BasicServiceSets: []*pb.WifiConfig_BasicServiceSet{
										{
											Bssid:     "11:11:11:11:11:11",
											Ssid:      "Corp-2G",
											Band:      pb.WifiConfig_RF_2GHZ,
											IfaceName: "wlan0",
										},
										{
											Bssid:     "33:33:33:33:33:33",
											Ssid:      "Corp-5G",
											Band:      pb.WifiConfig_RF_5GHZ,
											IfaceName: "wlan1",
										},
									},
								},
							},
						},
					},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetWifiConfig(context.Background())
	if err != nil {
		t.Fatalf("GetWifiConfig() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetWifiGetConfig() == nil {
		t.Fatalf("expected wifi_get_config request, got: %#v", tr.lastRequest)
	}
	if got.CountryCode != "BR" || !got.SetupComplete || got.MacWan != "AA:AA:AA:AA:AA:AA" || got.MacLan != "BB:BB:BB:BB:BB:BB" {
		t.Fatalf("unexpected root wifi config mapping: %+v", got)
	}
	if got.BootCount != 12 || got.Incarnation != 12345 || got.WanHostDscpMark != 46 {
		t.Fatalf("unexpected counters/marks mapping: %+v", got)
	}
	if len(got.Networks) != 2 {
		t.Fatalf("expected 2 networks, got %d", len(got.Networks))
	}

	// sorted by domain + ipv4 + vlan
	first := got.Networks[0]
	if first.Domain != "guest" || first.Ipv4 != "192.168.2.1" || first.Vlan != 20 {
		t.Fatalf("unexpected first network ordering: %+v", first)
	}
	if first.Dhcpv4Start != 100 || first.Dhcpv4LeaseDurationSeconds != 7200 {
		t.Fatalf("unexpected first network DHCP mapping: %+v", first)
	}
	if first.Dhcpv4End != 0 {
		t.Fatalf("expected dhcpv4_end default zero with current protobuf, got %+v", first)
	}
	if len(first.BasicServiceSets) != 1 || first.BasicServiceSets[0].Band != "rf_5ghz" {
		t.Fatalf("unexpected first network BSS mapping: %+v", first)
	}

	second := got.Networks[1]
	if second.Domain != "lan" || second.Ipv4 != "192.168.1.1" || second.Vlan != 10 {
		t.Fatalf("unexpected second network ordering: %+v", second)
	}
	if len(second.BasicServiceSets) != 2 {
		t.Fatalf("expected 2 BSS in second network, got %+v", second)
	}
	if second.BasicServiceSets[0].Ssid != "Corp-2G" || second.BasicServiceSets[0].Band != "rf_2ghz" || second.BasicServiceSets[0].InterfaceName != "wlan0" {
		t.Fatalf("unexpected second network first BSS mapping: %+v", second.BasicServiceSets[0])
	}
}

func TestGetWifiConfig_EmptyWhenNil(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetConfig{
					WifiGetConfig: nil,
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetWifiConfig(context.Background())
	if err != nil {
		t.Fatalf("GetWifiConfig() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil config")
	}
	if len(got.Networks) != 0 {
		t.Fatalf("expected empty networks, got %+v", got.Networks)
	}
}

func TestGetWifiConfig_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetWifiConfig(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetWifiConfig_UnexpectedResponseType(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetWifiConfig(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
