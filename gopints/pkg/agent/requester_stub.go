//go:build !linux

package agent

import (
	"errors"
	"time"
)

type stubRequester struct{}

func (r *stubRequester) RequestLine(chip string, offset int, edge Edge, debounce time.Duration, handler func(LineEvent)) (Line, error) {
	return nil, errors.New("gpiocdev is only supported on linux")
}

func DefaultRequester() LineRequester {
	return &stubRequester{}
}
