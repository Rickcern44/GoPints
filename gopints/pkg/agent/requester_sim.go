package agent

import (
	"fmt"
	"sync"
	"time"
)

// SimulatorRequester is a LineRequester that fires synthetic pulse events
// without real GPIO hardware. Safe for use on any platform.
//
// Call Trigger to inject a burst of pulses into a registered line.
type SimulatorRequester struct {
	mu       sync.Mutex
	handlers map[string]func(LineEvent)
}

// NewSimulatorRequester creates a SimulatorRequester ready for use.
func NewSimulatorRequester() *SimulatorRequester {
	return &SimulatorRequester{
		handlers: make(map[string]func(LineEvent)),
	}
}

// RequestLine registers the handler for the given chip/pin and returns a simLine.
// The returned line's Close method is a no-op.
func (r *SimulatorRequester) RequestLine(chip string, offset int, edge Edge, debounce time.Duration, handler func(LineEvent)) (Line, error) {
	key := lineKey(chip, offset)
	r.mu.Lock()
	r.handlers[key] = handler
	r.mu.Unlock()
	return &simLine{key: key, sim: r}, nil
}

// Trigger fires totalPulses synthetic rising-edge events into the line identified
// by chip and pin, spacing them evenly at hz pulses per second.
// It returns immediately; pulses are fired in a background goroutine.
func (r *SimulatorRequester) Trigger(chip string, pin int, hz float64, totalPulses int) {
	key := lineKey(chip, pin)
	r.mu.Lock()
	handler := r.handlers[key]
	r.mu.Unlock()
	if handler == nil {
		return
	}
	go func() {
		interval := time.Duration(float64(time.Second) / hz)
		for i := 0; i < totalPulses; i++ {
			handler(LineEvent{
				Timestamp: uint64(time.Now().UnixNano()),
				Type:      1, // rising edge
			})
			time.Sleep(interval)
		}
	}()
}

type simLine struct {
	key string
	sim *SimulatorRequester
}

func (l *simLine) Close() error {
	l.sim.mu.Lock()
	delete(l.sim.handlers, l.key)
	l.sim.mu.Unlock()
	return nil
}

func lineKey(chip string, pin int) string {
	return fmt.Sprintf("%s:%d", chip, pin)
}
