package agent

import (
	"sync"
	"testing"
	"time"
)

type countingHandler struct {
	mu    sync.Mutex
	count int
}

func (h *countingHandler) count_() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.count
}

func (h *countingHandler) handle(evt LineEvent) {
	h.mu.Lock()
	h.count++
	h.mu.Unlock()
}

func TestSimulatorRequester_RegistersHandler(t *testing.T) {
	sim := NewSimulatorRequester()
	h := &countingHandler{}

	line, err := sim.RequestLine("gpiochip0", 17, EdgeRising, 0, h.handle)
	if err != nil {
		t.Fatalf("RequestLine failed: %v", err)
	}
	defer line.Close()

	key := lineKey("gpiochip0", 17)
	sim.mu.Lock()
	_, ok := sim.handlers[key]
	sim.mu.Unlock()
	if !ok {
		t.Error("handler was not registered in map")
	}
}

func TestSimulatorRequester_TriggerDeliversPulses(t *testing.T) {
	sim := NewSimulatorRequester()
	h := &countingHandler{}

	line, err := sim.RequestLine("gpiochip0", 17, EdgeRising, 0, h.handle)
	if err != nil {
		t.Fatalf("RequestLine failed: %v", err)
	}
	defer line.Close()

	const want = 10
	sim.Trigger("gpiochip0", 17, 500, want) // 500 Hz → 2ms between pulses

	// Wait long enough for all pulses to fire (10 pulses * 2ms = 20ms, give 5x margin).
	time.Sleep(200 * time.Millisecond)

	if got := h.count_(); got != want {
		t.Errorf("want %d pulses, got %d", want, got)
	}
}

func TestSimulatorRequester_TriggerUnknownLine_NoOp(t *testing.T) {
	sim := NewSimulatorRequester()
	// Should not panic when no handler is registered.
	sim.Trigger("gpiochip0", 99, 100, 5)
	time.Sleep(100 * time.Millisecond) // let goroutine exit cleanly
}

func TestSimulatorRequester_CloseRemovesHandler(t *testing.T) {
	sim := NewSimulatorRequester()
	h := &countingHandler{}

	line, _ := sim.RequestLine("gpiochip0", 17, EdgeRising, 0, h.handle)
	line.Close()

	sim.Trigger("gpiochip0", 17, 500, 10)
	time.Sleep(100 * time.Millisecond)

	if got := h.count_(); got != 0 {
		t.Errorf("expected 0 pulses after Close, got %d", got)
	}
}

func TestLineKey_Format(t *testing.T) {
	if got := lineKey("gpiochip0", 17); got != "gpiochip0:17" {
		t.Errorf("want gpiochip0:17, got %s", got)
	}
}
