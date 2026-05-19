package flow_test

import (
	"testing"
	"time"

	"github.com/rickcern44/gopints/pkg/flow"
)

const idleTimeout = 50 * time.Millisecond

// drainEvent reads one PourEvent within timeout, failing the test if none arrives.
func drainEvent(t *testing.T, ch <-chan flow.PourEvent, timeout time.Duration) flow.PourEvent {
	t.Helper()
	select {
	case e := <-ch:
		return e
	case <-time.After(timeout):
		t.Fatal("timed out waiting for PourEvent")
		return flow.PourEvent{}
	}
}

func newMeter(t *testing.T, pulsesPerLiter float64) (*flow.Meter, chan flow.PourEvent) {
	t.Helper()
	ch := make(chan flow.PourEvent, 16)
	m := flow.NewMeter(1, pulsesPerLiter, idleTimeout, ch)
	return m, ch
}

func TestMeter_PourStarted(t *testing.T) {
	m, ch := newMeter(t, 450)
	before := time.Now()
	m.HandlePulse()

	e := drainEvent(t, ch, time.Second)
	if e.Type != flow.PourStarted {
		t.Errorf("want PourStarted, got %v", e.Type)
	}
	if e.TapID != 1 {
		t.Errorf("want TapID 1, got %d", e.TapID)
	}
	if e.StartedAt.Before(before) {
		t.Error("StartedAt should not be before the test started")
	}
	if !e.EndedAt.IsZero() {
		t.Error("EndedAt should be zero for PourStarted")
	}
}

func TestMeter_PourUpdated(t *testing.T) {
	m, ch := newMeter(t, 450)
	m.HandlePulse() // PourStarted
	drainEvent(t, ch, time.Second)

	m.HandlePulse() // PourUpdated
	e := drainEvent(t, ch, time.Second)
	if e.Type != flow.PourUpdated {
		t.Errorf("want PourUpdated, got %v", e.Type)
	}
	if e.VolumeMl <= 0 {
		t.Error("VolumeMl should be positive after second pulse")
	}
}

func TestMeter_PourEnded(t *testing.T) {
	m, ch := newMeter(t, 450)
	m.HandlePulse()
	drainEvent(t, ch, time.Second) // PourStarted

	// Wait for idle timer to fire.
	e := drainEvent(t, ch, 5*idleTimeout)
	if e.Type != flow.PourEnded {
		t.Errorf("want PourEnded, got %v", e.Type)
	}
	if e.EndedAt.IsZero() {
		t.Error("EndedAt should be set on PourEnded")
	}
}

func TestMeter_VolumeCalculation(t *testing.T) {
	cases := []struct {
		pulses         int
		pulsesPerLiter float64
		wantMl         float64
	}{
		{450, 450, 1000},
		{225, 450, 500},
		{90, 450, 200},
		{0, 450, 0},
	}
	for _, tc := range cases {
		ch := make(chan flow.PourEvent, 64)
		m := flow.NewMeter(1, tc.pulsesPerLiter, idleTimeout, ch)
		for i := 0; i < tc.pulses; i++ {
			m.HandlePulse()
		}
		// Drain events until PourEnded.
		deadline := time.After(5 * idleTimeout)
		var last flow.PourEvent
		for {
			select {
			case e := <-ch:
				last = e
				if e.Type == flow.PourEnded {
					goto done
				}
			case <-deadline:
				if tc.pulses == 0 {
					goto done // no pulses → no events expected
				}
				t.Fatalf("pulses=%d: timed out waiting for PourEnded", tc.pulses)
			}
		}
	done:
		if tc.pulses > 0 && last.VolumeMl != tc.wantMl {
			t.Errorf("pulses=%d: want %.1fml, got %.1fml", tc.pulses, tc.wantMl, last.VolumeMl)
		}
	}
}

func TestMeter_ZeroPulsesPerLiter(t *testing.T) {
	m, ch := newMeter(t, 0)
	m.HandlePulse()
	e := drainEvent(t, ch, time.Second)
	if e.VolumeMl != 0 {
		t.Errorf("want 0ml with zero calibration, got %.2f", e.VolumeMl)
	}
}

func TestMeter_SequentialPours(t *testing.T) {
	m, ch := newMeter(t, 450)

	// First pour.
	m.HandlePulse()
	drainEvent(t, ch, time.Second)   // PourStarted
	drainEvent(t, ch, 5*idleTimeout) // PourEnded

	// Second pour should start fresh.
	m.HandlePulse()
	e := drainEvent(t, ch, time.Second)
	if e.Type != flow.PourStarted {
		t.Errorf("second pour: want PourStarted, got %v", e.Type)
	}
	// Volume resets: one pulse at 450 p/L = ~2.2ml, not 1000ml+.
	if e.VolumeMl >= 10 {
		t.Errorf("second pour volume should be small (reset), got %.2fml", e.VolumeMl)
	}
}

func TestMeter_FullChannelDoesNotBlock(t *testing.T) {
	ch := make(chan flow.PourEvent) // unbuffered — will be full immediately
	m := flow.NewMeter(1, 450, idleTimeout, ch)

	done := make(chan struct{})
	go func() {
		m.HandlePulse() // should not block even though channel is full
		close(done)
	}()

	select {
	case <-done:
		// pass
	case <-time.After(time.Second):
		t.Fatal("HandlePulse blocked on full channel")
	}
}
