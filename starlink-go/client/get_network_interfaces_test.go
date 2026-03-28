package client

import (
	"context"
	"errors"
	"testing"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetNetworkInterfaces_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_GetNetworkInterfaces{
					GetNetworkInterfaces: &pb.GetNetworkInterfacesResponse{
						NetworkInterfaces: []*pb.NetworkInterface{
							{
								Name:          "wan0",
								Up:            true,
								MacAddress:    "AA:AA:AA:AA:AA:AA",
								Ipv4Addresses: []string{"10.0.0.2"},
								Ipv6Addresses: []string{"2001::2", "2001::1"},
								RxStats: &pb.NetworkInterface_RxStats{
									Bytes:   1000,
									Packets: 100,
								},
								TxStats: &pb.NetworkInterface_TxStats{
									Bytes:   2000,
									Packets: 200,
								},
								Interface: &pb.NetworkInterface_Ethernet{
									Ethernet: &pb.EthernetNetworkInterface{
										LinkDetected:      true,
										SpeedMbps:         1000,
										AutonegotiationOn: true,
										Duplex:            pb.EthernetNetworkInterface_FULL,
									},
								},
							},
							{
								Name:          "ra0",
								Up:            true,
								MacAddress:    "BB:BB:BB:BB:BB:BB",
								Ipv4Addresses: []string{"192.168.1.1"},
								Ipv6Addresses: []string{"fd00::1"},
								RxStats: &pb.NetworkInterface_RxStats{
									Bytes:   3000,
									Packets: 300,
								},
								TxStats: &pb.NetworkInterface_TxStats{
									Bytes:   4000,
									Packets: 400,
								},
								Interface: &pb.NetworkInterface_Wifi{
									Wifi: &pb.WifiNetworkInterface{
										Channel:     44,
										LinkQuality: 0.93,
									},
								},
							},
							{
								Name:          "bridge0",
								Up:            true,
								MacAddress:    "CC:CC:CC:CC:CC:CC",
								Ipv4Addresses: []string{"172.16.0.1"},
								Ipv6Addresses: []string{},
								Interface: &pb.NetworkInterface_Bridge{
									Bridge: &pb.BridgeNetworkInterface{
										MemberNames: []string{"lan1", "lan0"},
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
	got, err := c.GetNetworkInterfaces(context.Background())
	if err != nil {
		t.Fatalf("GetNetworkInterfaces() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetGetNetworkInterfaces() == nil {
		t.Fatalf("expected get_network_interfaces request, got: %#v", tr.lastRequest)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 interfaces, got %d", len(got))
	}

	// sorted by name
	if got[0].Name != "bridge0" || got[1].Name != "ra0" || got[2].Name != "wan0" {
		t.Fatalf("unexpected sort by name: %+v", got)
	}

	bridge := got[0]
	if bridge.Bridge == nil || len(bridge.Bridge.MemberNames) != 2 || bridge.Bridge.MemberNames[0] != "lan0" || bridge.Bridge.MemberNames[1] != "lan1" {
		t.Fatalf("unexpected bridge mapping: %+v", bridge)
	}

	wifi := got[1]
	if wifi.Wifi == nil || wifi.Wifi.Channel != 44 || wifi.Wifi.LinkQuality != 0.93 {
		t.Fatalf("unexpected wifi mapping: %+v", wifi)
	}
	if wifi.RxStats.Bytes != 3000 || wifi.RxStats.Packets != 300 || wifi.TxStats.Bytes != 4000 || wifi.TxStats.Packets != 400 {
		t.Fatalf("unexpected stats mapping: %+v", wifi)
	}

	eth := got[2]
	if eth.Ethernet == nil || !eth.Ethernet.LinkDetected || eth.Ethernet.SpeedMbps != 1000 || !eth.Ethernet.AutonegotiationOn || eth.Ethernet.Duplex != "full" {
		t.Fatalf("unexpected ethernet mapping: %+v", eth)
	}
	if len(eth.Ipv6Addresses) != 2 || eth.Ipv6Addresses[0] != "2001::1" || eth.Ipv6Addresses[1] != "2001::2" {
		t.Fatalf("unexpected ipv6 sorting: %+v", eth.Ipv6Addresses)
	}
}

func TestGetNetworkInterfaces_EmptyResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_GetNetworkInterfaces{
					GetNetworkInterfaces: nil,
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetNetworkInterfaces(context.Background())
	if err != nil {
		t.Fatalf("GetNetworkInterfaces() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestGetNetworkInterfaces_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetNetworkInterfaces(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetNetworkInterfaces_UnexpectedResponseType(t *testing.T) {
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
	_, err := c.GetNetworkInterfaces(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
