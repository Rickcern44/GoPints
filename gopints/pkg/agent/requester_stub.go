//go:build !linux

package agent

import (
	"errors"
	"time"
)

type stubRequester struct{}

func (r *stubRequester) RequestLine(_ string, _ int, _ Edge, _ time.Duration, _ func(LineEvent)) (Line, error) {
	return nil, errors.New("gpiocdev is only supported on linux")
}

func DefaultRequester() LineRequester {
	return &stubRequester{}
}
