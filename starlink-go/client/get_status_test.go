package client

import (
	"context"
	"testing"
	"time"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetStatus_AcceptsDishStatusResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_DishGetStatus{
					DishGetStatus: &pb.DishGetStatusResponse{
						DeviceInfo: &pb.DeviceInfo{
							Id:              "dish-1",
							HardwareVersion: "hw",
							SoftwareVersion: "sw",
						},
						DeviceState:           &pb.DeviceState{UptimeS: 111},
						PopPingDropRate:       0.5,
						PopPingLatencyMs:      48,
						UplinkThroughputBps:   100,
						DownlinkThroughputBps: 200,
					},
				},
			}, nil
		},
	}

	c := newGRPCClient(tr, Config{Timeout: time.Second})
	got, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus() unexpected error: %v", err)
	}

	if got.DeviceID != "dish-1" || got.UptimeSeconds != 111 {
		t.Fatalf("unexpected status mapping: %+v", got)
	}
	if got.PopPingDropRate != 0.5 || got.PopPingLatencyMs != 48 {
		t.Fatalf("unexpected ping mapping: %+v", got)
	}
}

func TestGetStatus_AcceptsWifiStatusResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetStatus{
					WifiGetStatus: &pb.WifiGetStatusResponse{
						DeviceInfo: &pb.DeviceInfo{
							Id:              "router-1",
							HardwareVersion: "hw-r",
							SoftwareVersion: "sw-r",
						},
						DeviceState:      &pb.DeviceState{UptimeS: 222},
						PopPingDropRate:  1.5,
						PopPingLatencyMs: 55,
					},
				},
			}, nil
		},
	}

	c := newGRPCClient(tr, Config{Timeout: time.Second})
	got, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus() unexpected error: %v", err)
	}

	if got.DeviceID != "router-1" || got.UptimeSeconds != 222 {
		t.Fatalf("unexpected status mapping: %+v", got)
	}
	if got.PopPingDropRate != 1.5 || got.PopPingLatencyMs != 55 {
		t.Fatalf("unexpected ping mapping: %+v", got)
	}
}
