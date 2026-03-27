package status

import (
	"errors"

	"google.golang.org/grpc/codes"
)

type Status struct {
	code codes.Code
	err  error
}

func (s *Status) Code() codes.Code {
	if s == nil {
		return codes.OK
	}
	return s.code
}

type coder interface {
	GRPCCode() codes.Code
}

func FromError(err error) (*Status, bool) {
	if err == nil {
		return &Status{code: codes.OK}, true
	}
	var c coder
	if errors.As(err, &c) {
		return &Status{code: c.GRPCCode(), err: err}, true
	}
	return nil, false
}
