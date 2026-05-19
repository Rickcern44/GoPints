package protocol

import (
	"encoding/binary"
	"errors"
)

const (
	msgTypePulse = 0x01
	PulseMsgSize = 10 // type(1) + tap_id(1) + timestamp_ns(8)
)

// PulseMessage is a single GPIO pulse event sent from agent to server over UDP.
type PulseMessage struct {
	TapID     uint8
	Timestamp uint64 // nanoseconds since Unix epoch
}

// Marshal encodes the message into a fixed 10-byte array.
func (m PulseMessage) Marshal() [PulseMsgSize]byte {
	var buf [PulseMsgSize]byte
	buf[0] = msgTypePulse
	buf[1] = m.TapID
	binary.BigEndian.PutUint64(buf[2:], m.Timestamp)
	return buf
}

// Unmarshal decodes a PulseMessage from a byte slice.
func Unmarshal(b []byte) (PulseMessage, error) {
	if len(b) < PulseMsgSize {
		return PulseMessage{}, errors.New("protocol: message too short")
	}
	if b[0] != msgTypePulse {
		return PulseMessage{}, errors.New("protocol: unknown message type")
	}
	return PulseMessage{
		TapID:     b[1],
		Timestamp: binary.BigEndian.Uint64(b[2:]),
	}, nil
}
