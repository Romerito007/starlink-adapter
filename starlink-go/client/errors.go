package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDeviceOffline = errors.New("starlink device offline")
	ErrTimeout       = errors.New("starlink request timeout")
	ErrUnavailable   = errors.New("starlink service unavailable")
	ErrUnsupported   = errors.New("starlink operation unsupported")
)

func normalizeError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", ErrTimeout, err)
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.DeadlineExceeded:
			return fmt.Errorf("%w: %v", ErrTimeout, err)
		case codes.Unavailable:
			return fmt.Errorf("%w: %v", ErrUnavailable, err)
		case codes.Unimplemented:
			return fmt.Errorf("%w: %v", ErrUnsupported, err)
		}
	}

	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "connection refused") || strings.Contains(msg, "no route to host") {
		return fmt.Errorf("%w: %v", ErrDeviceOffline, err)
	}

	return err
}

func isTransientError(err error) bool {
	return errors.Is(err, ErrTimeout) || errors.Is(err, ErrUnavailable) || errors.Is(err, ErrDeviceOffline)
}
