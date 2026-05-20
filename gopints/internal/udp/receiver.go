package udp

import (
	"fmt"
	"net"

	"github.com/rickcern44/gopints/pkg/protocol"
)

// Receiver listens on a UDP address and decodes incoming PulseMessages.
type Receiver struct {
	conn *net.UDPConn
}

// NewReceiver binds to the given UDP address and returns a Receiver.
func NewReceiver(addr string) (*Receiver, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("udp receiver: resolve %s: %w", addr, err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("udp receiver: listen %s: %w", addr, err)
	}
	return &Receiver{conn: conn}, nil
}

// Receive blocks until a valid PulseMessage arrives.
func (r *Receiver) Receive() (protocol.PulseMessage, error) {
	buf := make([]byte, protocol.PulseMsgSize)
	n, err := r.conn.Read(buf)
	if err != nil {
		return protocol.PulseMessage{}, fmt.Errorf("udp receiver: read: %w", err)
	}
	return protocol.Unmarshal(buf[:n])
}

// Close releases the underlying UDP socket, unblocking any pending Receive.
func (r *Receiver) Close() error {
	return r.conn.Close()
}
