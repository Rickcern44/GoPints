package flow

import (
	"sync"
	"time"
)

const defaultIdleTimeout = 2 * time.Second

// Meter accumulates pulse events for a single tap and emits PourEvents.
// It is safe for concurrent use.
//
// Meter.HandlePulse satisfies the agent.PulseHandler interface via structural
// typing — no import of pkg/agent is required.
type Meter struct {
	tapID          uint8
	pulsesPerLiter float64
	idleTimeout    time.Duration
	events         chan<- PourEvent

	mu        sync.Mutex
	pulses    int
	pouring   bool
	startedAt time.Time
	idleTimer *time.Timer
}

// NewMeter creates a Meter for the given tap.
//
//   - pulsesPerLiter is the flow sensor calibration constant (e.g. 450).
//   - idleTimeout is how long after the last pulse before a pour is declared
//     finished (pass 0 to use the 2-second default).
//   - events is the channel that PourEvents are sent to; it must be non-nil.
func NewMeter(tapID uint8, pulsesPerLiter float64, idleTimeout time.Duration, events chan<- PourEvent) *Meter {
	if idleTimeout <= 0 {
		idleTimeout = defaultIdleTimeout
	}
	return &Meter{
		tapID:          tapID,
		pulsesPerLiter: pulsesPerLiter,
		idleTimeout:    idleTimeout,
		events:         events,
	}
}

// HandlePulse records a single pulse and updates pour state.
// It satisfies agent.PulseHandler via structural typing.
func (m *Meter) HandlePulse() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	m.pulses++

	if !m.pouring {
		m.pouring = true
		m.startedAt = now
		m.emit(PourStarted, now)
	} else {
		m.emit(PourUpdated, now)
	}

	if m.idleTimer != nil {
		m.idleTimer.Reset(m.idleTimeout)
	} else {
		m.idleTimer = time.AfterFunc(m.idleTimeout, m.onIdle)
	}
}

// onIdle is called by the idle timer when no pulses have arrived.
func (m *Meter) onIdle() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emit(PourEnded, time.Now())
	m.pouring = false
	m.pulses = 0
	m.idleTimer = nil
}

func (m *Meter) volumeMl() float64 {
	if m.pulsesPerLiter <= 0 {
		return 0
	}
	return float64(m.pulses) / m.pulsesPerLiter * 1000
}

func (m *Meter) emit(t EventType, now time.Time) {
	e := PourEvent{
		Type:      t,
		TapID:     m.tapID,
		VolumeMl:  m.volumeMl(),
		StartedAt: m.startedAt,
	}
	if t == PourEnded {
		e.EndedAt = now
	}
	// Non-blocking send: a slow consumer should not stall GPIO processing.
	select {
	case m.events <- e:
	default:
	}
}
