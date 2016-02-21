package udp

import (
	"net"
)

// Writer is a simple implementation of the writer interface, for sending
// UDP messages
type Writer struct {
	address    *net.UDPAddr
	connection *net.UDPConn
}

// NewUDPWriter returns a new UDP Writer
func NewUDPWriter(givenConnection *net.UDPConn, givenAddress *net.UDPAddr) *Writer {
	return &Writer{address: givenAddress, connection: givenConnection}
}

// Write is an implementation of the standard io.Writer interface
func (udpWriter *Writer) Write(givenBytes []byte) (int, error) {
	return udpWriter.connection.WriteToUDP(givenBytes, udpWriter.address)
}
