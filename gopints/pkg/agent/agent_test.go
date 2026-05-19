package agent

import (
	"context"
	"testing"
	"time"
)

type mockLine struct {
	closed bool
}

func (m *mockLine) Close() error {
	m.closed = true
	return nil
}

type mockRequester struct {
	lastChip     string
	lastOffset   int
	lastEdge     Edge
	lastDebounce time.Duration
	lastHandler  func(LineEvent)
	mockLine     *mockLine
}

func (m *mockRequester) RequestLine(chip string, offset int, edge Edge, debounce time.Duration, handler func(LineEvent)) (Line, error) {
	m.lastChip = chip
	m.lastOffset = offset
	m.lastEdge = edge
	m.lastDebounce = debounce
	m.lastHandler = handler
	return m.mockLine, nil
}

type mockPulseHandler struct {
	pulses int
}

func (m *mockPulseHandler) HandlePulse() {
	m.pulses++
}

func TestObserver_Listen(t *testing.T) {
	chip := "gpiochip0"
	pin := 17
	handler := &mockPulseHandler{}
	line := &mockLine{}
	requester := &mockRequester{mockLine: line}

	o := NewObserver(chip, pin, handler, WithRequester(requester), WithEdge(EdgeFalling), WithDebounce(10*time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- o.Listen(ctx)
	}()

	// Wait for Listen to call RequestLine
	time.Sleep(10 * time.Millisecond)

	if requester.lastChip != chip {
		t.Errorf("expected chip %s, got %s", chip, requester.lastChip)
	}
	if requester.lastOffset != pin {
		t.Errorf("expected pin %d, got %d", pin, requester.lastOffset)
	}
	if requester.lastEdge != EdgeFalling {
		t.Errorf("expected edge %v, got %v", EdgeFalling, requester.lastEdge)
	}
	if requester.lastDebounce != 10*time.Millisecond {
		t.Errorf("expected debounce %v, got %v", 10*time.Millisecond, requester.lastDebounce)
	}

	// Simulate a pulse
	if requester.lastHandler == nil {
		t.Fatal("handler was not set")
	}
	requester.lastHandler(LineEvent{})

	if handler.pulses != 1 {
		t.Errorf("expected 1 pulse, got %d", handler.pulses)
	}

	// Cancel context and check if Listen returns
	cancel()
	select {
	case err := <-errCh:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Listen did not return after context cancellation")
	}

	if !line.closed {
		t.Error("expected line to be closed")
	}
}
