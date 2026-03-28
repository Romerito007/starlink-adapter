package client

import (
	"context"
	"errors"
	"math"
	"testing"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetRadioStats_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_GetRadioStats{
					GetRadioStats: &pb.GetRadioStatsResponse{
						RadioStats: []*pb.RadioStats{
							{
								Band: pb.WifiConfig_RF_5GHZ,
								RxStats: &pb.NetworkInterface_RxStats{
									Packets:     1200,
									FrameErrors: 12,
								},
								TxStats: &pb.NetworkInterface_TxStats{
									Packets: 900,
								},
								ThermalStatus: &pb.RadioStats_ThermalStatus{
									Temp2:     67.5,
									DutyCycle: 55,
								},
								AntennaStatus: &pb.RadioStats_AntennaStatus{
									Rssi1: -54.5,
									Rssi2: -55.5,
									Rssi3: float32(math.NaN()),
									Rssi4: -57.5,
								},
							},
							{
								Band: pb.WifiConfig_RF_2GHZ,
								RxStats: &pb.NetworkInterface_RxStats{
									Packets:     300,
									FrameErrors: 3,
								},
								ThermalStatus: &pb.RadioStats_ThermalStatus{
									Temp2:     math.NaN(),
									DutyCycle: 22,
								},
							},
						},
					},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetRadioStats(context.Background())
	if err != nil {
		t.Fatalf("GetRadioStats() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetGetRadioStats() == nil {
		t.Fatalf("expected get_radio_stats request, got: %#v", tr.lastRequest)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 radio entries, got %d", len(got))
	}

	// Sorted by band.
	if got[0].Band != "rf_2ghz" || got[1].Band != "rf_5ghz" {
		t.Fatalf("unexpected band ordering: %+v", got)
	}

	first := got[0]
	if first.RxStats.Packets != 300 || first.RxStats.FrameErrors != 3 {
		t.Fatalf("unexpected first rx mapping: %+v", first)
	}
	if first.TxStats.Packets != 0 || first.TxStats.FrameErrors != 0 {
		t.Fatalf("expected zero tx stats when absent: %+v", first)
	}
	if first.ThermalStatus.Temp != 0 || first.ThermalStatus.DutyCycle != 22 {
		t.Fatalf("unexpected first thermal mapping: %+v", first)
	}

	second := got[1]
	if second.RxStats.Packets != 1200 || second.RxStats.FrameErrors != 12 {
		t.Fatalf("unexpected second rx mapping: %+v", second)
	}
	if second.TxStats.Packets != 900 {
		t.Fatalf("unexpected second tx mapping: %+v", second)
	}
	if second.ThermalStatus.Temp != 67.5 || second.ThermalStatus.DutyCycle != 55 {
		t.Fatalf("unexpected second thermal mapping: %+v", second)
	}
	if second.AntennaStatus.Rssi1 != -54.5 || second.AntennaStatus.Rssi2 != -55.5 || second.AntennaStatus.Rssi3 != 0 || second.AntennaStatus.Rssi4 != -57.5 {
		t.Fatalf("unexpected second antenna mapping: %+v", second)
	}
}

func TestGetRadioStats_EmptyResponse(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_GetRadioStats{
					GetRadioStats: nil,
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetRadioStats(context.Background())
	if err != nil {
		t.Fatalf("GetRadioStats() unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d", len(got))
	}
}

func TestGetRadioStats_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetRadioStats(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetRadioStats_UnexpectedResponseType(t *testing.T) {
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
	_, err := c.GetRadioStats(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
