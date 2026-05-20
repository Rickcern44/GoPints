package agent

import (
	"context"
	"fmt"
	"time"
)

// Line represents a GPIO line that can be closed.
type Line interface {
	Close() error
}

// LineRequester defines the interface for requesting a GPIO line.
type LineRequester interface {
	RequestLine(chip string, offset int, edge Edge, debounce time.Duration, handler func(LineEvent)) (Line, error)
}

// Observer monitors a GPIO pin for pulses.
type Observer struct {
	chipName string
	pin      int
	handler  PulseHandler
	edge     Edge
	debounce time.Duration

	requester LineRequester
}

// NewObserver creates a new GPIO observer.
func NewObserver(chipName string, pin int, handler PulseHandler, opts ...Option) *Observer {
	o := &Observer{
		chipName:  chipName,
		pin:       pin,
		handler:   handler,
		edge:      EdgeRising, // Default to rising edge
		requester: DefaultRequester(),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Observer) Listen(ctx context.Context) error {
	handler := func(_ LineEvent) {
		o.handler.HandlePulse()
	}

	line, err := o.requester.RequestLine(o.chipName, o.pin, o.edge, o.debounce, handler)
	if err != nil {
		return fmt.Errorf("failed to request gpio line %d on chip %s: %w", o.pin, o.chipName, err)
	}
	defer line.Close()

	<-ctx.Done()
	return ctx.Err()
}
