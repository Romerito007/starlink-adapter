package client

import (
	"context"
	"errors"
	"testing"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetDhcpLeases_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{
						DhcpServers: []*pb.DhcpServer{
							{
								Domain: "guest",
								Leases: []*pb.DhcpLease{
									{
										IpAddress:   "192.168.2.20",
										MacAddress:  "BB:00:00:00:00:02",
										Hostname:    "tablet",
										ExpiresTime: "2026-03-28T16:00:00Z",
										Active:      true,
										ClientId:    22,
									},
								},
							},
							{
								Domain: "lan",
								Leases: []*pb.DhcpLease{
									{
										IpAddress:   "192.168.1.20",
										MacAddress:  "CC:00:00:00:00:03",
										Hostname:    "phone-c",
										ExpiresTime: "2026-03-28T15:00:00Z",
										Active:      true,
										ClientId:    3,
									},
									{
										IpAddress:   "192.168.1.10",
										MacAddress:  "AA:00:00:00:00:01",
										Hostname:    "phone-a",
										ExpiresTime: "2026-03-28T14:00:00Z",
										Active:      false,
										ClientId:    1,
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
	got, err := c.GetDhcpLeases(context.Background())
	if err != nil {
		t.Fatalf("GetDhcpLeases() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetGetStatus() == nil {
		t.Fatalf("expected get_status request, got: %#v", tr.lastRequest)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 leases, got %d", len(got))
	}

	// Stable sort: domain + ip + mac
	if got[0].Domain != "guest" || got[0].IpAddress != "192.168.2.20" || got[0].MacAddress != "BB:00:00:00:00:02" {
		t.Fatalf("unexpected first lease ordering: %+v", got[0])
	}
	if got[1].Domain != "lan" || got[1].IpAddress != "192.168.1.10" || got[1].MacAddress != "AA:00:00:00:00:01" {
		t.Fatalf("unexpected second lease ordering: %+v", got[1])
	}
	if got[2].Domain != "lan" || got[2].IpAddress != "192.168.1.20" || got[2].MacAddress != "CC:00:00:00:00:03" {
		t.Fatalf("unexpected third lease ordering: %+v", got[2])
	}
	if got[1].Hostname != "phone-a" || got[1].ExpiresTime != "2026-03-28T14:00:00Z" || got[1].Active || got[1].ClientID != 1 {
		t.Fatalf("unexpected lease mapping: %+v", got[1])
	}
}

func TestGetDhcpLeases_EmptyWhenNoServers(t *testing.T) {
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
	got, err := c.GetDhcpLeases(context.Background())
	if err != nil {
		t.Fatalf("GetDhcpLeases() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestGetDhcpLeases_EmptyWhenWifiStatusNil(t *testing.T) {
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
	got, err := c.GetDhcpLeases(context.Background())
	if err != nil {
		t.Fatalf("GetDhcpLeases() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestGetDhcpLeases_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetDhcpLeases(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetDhcpLeases_UnexpectedResponseType(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_DishGetStatus{DishGetStatus: &pb.DishGetStatusResponse{}},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetDhcpLeases(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
