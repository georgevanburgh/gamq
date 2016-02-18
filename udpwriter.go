package gamq

import (
	"net"
)

type UdpWriter struct {
	address    *net.UDPAddr
	connection *net.UDPConn
}

func NewUdpWriter(givenConnection *net.UDPConn, givenAddress *net.UDPAddr) *UdpWriter {
	return &UdpWriter{address: givenAddress, connection: givenConnection}
}

func (udpWriter *UdpWriter) Write(givenBytes []byte) (int, error) {
	return udpWriter.connection.WriteToUDP(givenBytes, udpWriter.address)
}
