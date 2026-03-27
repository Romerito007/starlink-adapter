package client

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewClient_UsesConfigTimeoutWithBackgroundContext(t *testing.T) {
	originalDial := dialTransport
	defer func() { dialTransport = originalDial }()

	dialTransport = func(ctx context.Context, address string) (transport, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	start := time.Now()
	_, err := NewClient(context.Background(), Config{
		Host:    "192.0.2.1",
		Port:    9200,
		Timeout: 50 * time.Millisecond,
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("dial took too long, timeout not respected: %s", elapsed)
	}
}

func TestNewClient_UsesDefaultTimeoutWhenConfigTimeoutInvalid(t *testing.T) {
	originalDial := dialTransport
	defer func() { dialTransport = originalDial }()

	dialTransport = func(ctx context.Context, address string) (transport, error) {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatalf("expected deadline in dial context")
		}

		remaining := time.Until(deadline)
		if remaining < 4*time.Second || remaining > 6*time.Second {
			t.Fatalf("expected default timeout around 5s, got %s", remaining)
		}

		return nil, context.DeadlineExceeded
	}

	_, err := NewClient(context.Background(), Config{
		Host:    "192.0.2.1",
		Port:    9200,
		Timeout: 0,
	})
	if err == nil {
		t.Fatalf("expected timeout error")
	}
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("expected ErrTimeout, got %v", err)
	}
}

func TestNewClient_RespectsParentContextCancellation(t *testing.T) {
	originalDial := dialTransport
	defer func() { dialTransport = originalDial }()

	dialTransport = func(ctx context.Context, address string) (transport, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	parentCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewClient(parentCtx, Config{
		Host:    "192.0.2.1",
		Port:    9200,
		Timeout: time.Second,
	})
	if err == nil {
		t.Fatalf("expected cancellation error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestNewClient_SucceedsWhenDialReturnsValidTransport(t *testing.T) {
	originalDial := dialTransport
	defer func() { dialTransport = originalDial }()

	fake := &fakeTransport{}
	dialTransport = func(ctx context.Context, address string) (transport, error) {
		return fake, nil
	}

	c, err := NewClient(context.Background(), Config{
		Host:    "127.0.0.1",
		Port:    9200,
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if c == nil {
		t.Fatalf("expected non-nil client")
	}

	if err := c.Close(); err != nil {
		t.Fatalf("expected close success, got error: %v", err)
	}
}
