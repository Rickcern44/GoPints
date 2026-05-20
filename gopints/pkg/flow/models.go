package flow

import "time"

// EventType classifies a PourEvent.
type EventType int

const (
	// PourStarted fires when the first pulse arrives after an idle period.
	PourStarted EventType = iota
	// PourUpdated fires on each subsequent pulse during an active pour.
	PourUpdated
	// PourEnded fires when no pulses have arrived within the idle timeout.
	PourEnded
)

// PourEvent describes the current state of a pour on a single tap.
type PourEvent struct {
	Type      EventType
	TapID     uint8
	VolumeMl  float64
	StartedAt time.Time
	EndedAt   time.Time // zero until PourEnded
}
