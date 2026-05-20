package udp

import (
	"fmt"
	"net"
	"time"

	"github.com/rickcern44/gopints/pkg/protocol"
)

// Sender sends PulseMessages to a UDP address.
type Sender struct {
	conn *net.UDPConn
}

// NewSender dials the given UDP address and returns a Sender.
func NewSender(addr string) (*Sender, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("udp sender: resolve %s: %w", addr, err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("udp sender: dial %s: %w", addr, err)
	}
	return &Sender{conn: conn}, nil
}

// Send encodes and transmits a PulseMessage. Errors are non-fatal for the
// caller — a dropped UDP datagram does not interrupt the GPIO stream.
func (s *Sender) Send(msg protocol.PulseMessage) error {
	buf := msg.Marshal()
	_, err := s.conn.Write(buf[:])
	return err
}

// Close releases the underlying UDP socket.
func (s *Sender) Close() error {
	return s.conn.Close()
}

// Handler implements the agent.PulseHandler interface by encoding each pulse
// as a PulseMessage and sending it over UDP. Errors are silently dropped to
// keep the GPIO event path non-blocking.
type Handler struct {
	TapID  uint8
	Sender *Sender
}

// HandlePulse satisfies agent.PulseHandler.
func (h *Handler) HandlePulse() {
	_ = h.Sender.Send(protocol.PulseMessage{
		TapID:     h.TapID,
		Timestamp: uint64(time.Now().UnixNano()),
	})
}
