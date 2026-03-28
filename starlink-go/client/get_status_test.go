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

func TestGetStatusDetailed_AcceptsWifiStatusResponse(t *testing.T) {
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
						DeviceState:         &pb.DeviceState{UptimeS: 222},
						Ipv4WanAddress:      "100.64.0.2",
						Ipv6WanAddresses:    []string{"2001::b", "2001::a"},
						PingDropRate:        0.2,
						PingDropRate_5M:     0.3,
						PingLatencyMs:       45,
						DishPingDropRate:    0.4,
						DishPingDropRate_5M: 0.5,
						DishPingLatencyMs:   55,
						PopPingDropRate:     0.6,
						PopPingDropRate_5M:  0.7,
						PopPingLatencyMs:    65,
						DishId:              "dish-123",
						UtcNs:               1700000000,
						PoeStats:            &pb.PoeStats{PoeState: pb.PoeState_POE_STATE_ON, PoePower: 12.5},
						SoftwareUpdateStats: &pb.WifiSoftwareUpdateStats{State: pb.WifiSoftwareUpdateState_FLASHING, SecondsSinceGetTargetVersions: 33},
					},
				},
			}, nil
		},
	}

	c := newGRPCClient(tr, Config{Timeout: time.Second})
	got, err := c.GetStatusDetailed(context.Background())
	if err != nil {
		t.Fatalf("GetStatusDetailed() unexpected error: %v", err)
	}

	if got.DeviceID != "router-1" || got.UptimeSeconds != 222 || got.Ipv4WanAddress != "100.64.0.2" {
		t.Fatalf("unexpected detailed root mapping: %+v", got)
	}
	if len(got.Ipv6WanAddresses) != 2 || got.Ipv6WanAddresses[0] != "2001::a" || got.Ipv6WanAddresses[1] != "2001::b" {
		t.Fatalf("unexpected ipv6 ordering mapping: %+v", got.Ipv6WanAddresses)
	}
	if got.PingDropRate != 0.2 || got.PingDropRate5m != 0.3 || got.PingLatencyMs != 45 {
		t.Fatalf("unexpected ping mapping: %+v", got)
	}
	if got.DishPingDropRate != 0.4 || got.DishPingDropRate5m != 0.5 || got.DishPingLatencyMs != 55 {
		t.Fatalf("unexpected dish ping mapping: %+v", got)
	}
	if got.PopPingDropRate != 0.6 || got.PopPingDropRate5m != 0.7 || got.PopPingLatencyMs != 65 {
		t.Fatalf("unexpected pop ping mapping: %+v", got)
	}
	if got.DishID != "dish-123" || got.UtcNs != 1700000000 {
		t.Fatalf("unexpected dish/utc mapping: %+v", got)
	}
	if got.PoeState != "poe_state_on" || got.PoePower != 12.5 {
		t.Fatalf("unexpected poe mapping: %+v", got)
	}
	if got.SoftwareUpdateState != "flashing" || got.SoftwareUpdateSecondsSinceGetTargetVersions != 33 {
		t.Fatalf("unexpected software update mapping: %+v", got)
	}
	// unsupported by current protobuf on wifi_get_status path
	if got.PopIpv6PingLatencyMs != 0 || got.PopIpv6PingDropRate != 0 || got.PopIpv6PingDropRate5m != 0 ||
		got.SecsSinceLastPublicIpv4Change != 0 || got.DishDisablementCode != "" || got.CalibrationPartitionsState != "" ||
		got.SetupRequirementState != "" || got.SoftwareUpdateRunningVersion != "" || got.PoeVin != 0 {
		t.Fatalf("unexpected defaults for unavailable fields: %+v", got)
	}
}

func TestGetStatusDetailed_AcceptsDishStatusResponse(t *testing.T) {
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
						DeviceState:         &pb.DeviceState{UptimeS: 111},
						PopPingDropRate:     0.5,
						PopPingLatencyMs:    48,
						SoftwareUpdateState: pb.SoftwareUpdateState_REBOOT_REQUIRED,
					},
				},
			}, nil
		},
	}

	c := newGRPCClient(tr, Config{Timeout: time.Second})
	got, err := c.GetStatusDetailed(context.Background())
	if err != nil {
		t.Fatalf("GetStatusDetailed() unexpected error: %v", err)
	}
	if got.DeviceID != "dish-1" || got.UptimeSeconds != 111 {
		t.Fatalf("unexpected detailed dish root mapping: %+v", got)
	}
	if got.PopPingDropRate != 0.5 || got.PopPingLatencyMs != 48 {
		t.Fatalf("unexpected detailed dish pop mapping: %+v", got)
	}
	if got.SoftwareUpdateState == "" {
		t.Fatalf("expected mapped software update enum from dish path, got %+v", got)
	}
}
