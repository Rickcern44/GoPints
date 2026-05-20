package agent

// PulseHandler is the interface that must be implemented to handle GPIO pulses.
type PulseHandler interface {
	HandlePulse()
}

// LineEvent represents a single GPIO edge event.
type LineEvent struct {
	// Timestamp is the event time in nanoseconds.
	Timestamp uint64
	// Type is the event type (Rising, Falling).
	Type int
}
