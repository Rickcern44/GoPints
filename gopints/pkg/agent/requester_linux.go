//go:build linux

package agent

import (
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type gpiocdevRequester struct{}

func (r *gpiocdevRequester) RequestLine(chip string, offset int, edge Edge, debounce time.Duration, handler func(LineEvent)) (Line, error) {
	var gpiopts []gpiocdev.LineReqOption

	switch edge {
	case EdgeRising:
		gpiopts = append(gpiopts, gpiocdev.WithRisingEdge)
	case EdgeFalling:
		gpiopts = append(gpiopts, gpiocdev.WithFallingEdge)
	case EdgeBoth:
		gpiopts = append(gpiopts, gpiocdev.WithBothEdges)
	}

	if debounce > 0 {
		gpiopts = append(gpiopts, gpiocdev.WithDebounce(debounce))
	}

	if handler != nil {
		gpiopts = append(gpiopts, gpiocdev.WithEventHandler(func(evt gpiocdev.LineEvent) {
			handler(LineEvent{
				Timestamp: uint64(evt.Timestamp),
				Type:      int(evt.Type),
			})
		}))
	}

	return gpiocdev.RequestLine(chip, offset, gpiopts...)
}

func DefaultRequester() LineRequester {
	return &gpiocdevRequester{}
}
