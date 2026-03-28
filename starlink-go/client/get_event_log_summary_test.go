package client

import (
	"context"
	"errors"
	"testing"

	pb "github.com/Romerito007/starlink-adapter/starlink-go/proto/gen/spacex/api/device"
)

func TestGetEventLogSummary_RequestAndMapping(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_WifiGetHistory{
					WifiGetHistory: &pb.WifiGetHistoryResponse{},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	got, err := c.GetEventLogSummary(context.Background())
	if err != nil {
		t.Fatalf("GetEventLogSummary() unexpected error: %v", err)
	}

	if tr.lastRequest == nil || tr.lastRequest.GetGetHistory() == nil {
		t.Fatalf("expected get_history request, got: %#v", tr.lastRequest)
	}
	if got == nil {
		t.Fatalf("expected non-nil summary")
	}
	if got.StartTimestampNs != 0 || got.CurrentTimestampNs != 0 {
		t.Fatalf("expected zero timestamps with current protobuf support, got %+v", got)
	}
	if got.Events == nil || len(got.Events) != 0 {
		t.Fatalf("expected empty non-nil events, got %+v", got.Events)
	}
}

func TestGetEventLogSummary_TransportError(t *testing.T) {
	transportErr := errors.New("boom")
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return nil, transportErr
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetEventLogSummary(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetEventLogSummary_UnexpectedResponseType(t *testing.T) {
	tr := &fakeTransport{
		handleFn: func(ctx context.Context, req *pb.Request) (*pb.Response, error) {
			return &pb.Response{
				Response: &pb.Response_DishGetHistory{
					DishGetHistory: &pb.DishGetHistoryResponse{},
				},
			}, nil
		},
	}

	c := newTestClient(t, tr)
	_, err := c.GetEventLogSummary(context.Background())
	if err == nil {
		t.Fatalf("expected error for unexpected response type")
	}
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("expected ErrUnsupported, got %v", err)
	}
}
