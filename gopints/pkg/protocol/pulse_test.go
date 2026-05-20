package protocol

import (
	"math"
	"testing"
)

func TestPulseMessage_RoundTrip(t *testing.T) {
	msg := PulseMessage{TapID: 3, Timestamp: 1_234_567_890_000_000_000}
	buf := msg.Marshal()
	got, err := Unmarshal(buf[:])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TapID != msg.TapID {
		t.Errorf("TapID: want %d, got %d", msg.TapID, got.TapID)
	}
	if got.Timestamp != msg.Timestamp {
		t.Errorf("Timestamp: want %d, got %d", msg.Timestamp, got.Timestamp)
	}
}

func TestPulseMessage_BoundaryValues(t *testing.T) {
	cases := []PulseMessage{
		{TapID: 0, Timestamp: 0},
		{TapID: 255, Timestamp: math.MaxUint64},
	}
	for _, msg := range cases {
		buf := msg.Marshal()
		got, err := Unmarshal(buf[:])
		if err != nil {
			t.Fatalf("TapID=%d: unexpected error: %v", msg.TapID, err)
		}
		if got != msg {
			t.Errorf("TapID=%d: want %+v, got %+v", msg.TapID, msg, got)
		}
	}
}

func TestUnmarshal_ShortBuffer(t *testing.T) {
	_, err := Unmarshal([]byte{0x01, 0x02})
	if err == nil {
		t.Fatal("expected error for short buffer, got nil")
	}
}

func TestUnmarshal_UnknownType(t *testing.T) {
	buf := [PulseMsgSize]byte{}
	buf[0] = 0xFF // unknown type
	_, err := Unmarshal(buf[:])
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestMarshal_SizeIsConstant(t *testing.T) {
	raw := PulseMessage{TapID: 1, Timestamp: 42}.Marshal()
	buf := raw[:]
	if len(buf) != PulseMsgSize {
		t.Errorf("want %d bytes, got %d", PulseMsgSize, len(buf))
	}
}
