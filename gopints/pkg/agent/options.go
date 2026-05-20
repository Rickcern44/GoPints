package agent

import "time"

// Edge defines the type of GPIO edge to detect.
type Edge int

const (
	// EdgeRising detects a low-to-high transition.
	EdgeRising Edge = iota
	// EdgeFalling detects a high-to-low transition.
	EdgeFalling
	// EdgeBoth detects any transition.
	EdgeBoth
)

// Option is a functional option for configuring an Observer.
type Option func(*Observer)

// WithEdge sets the edge detection type.
func WithEdge(edge Edge) Option {
	return func(o *Observer) {
		o.edge = edge
	}
}

// WithDebounce sets the debounce duration for the GPIO line.
func WithDebounce(d time.Duration) Option {
	return func(o *Observer) {
		o.debounce = d
	}
}

// WithRequester overrides the LineRequester used by the Observer.
// Useful for injecting a SimulatorRequester or a custom implementation.
func WithRequester(r LineRequester) Option {
	return func(o *Observer) {
		o.requester = r
	}
}
