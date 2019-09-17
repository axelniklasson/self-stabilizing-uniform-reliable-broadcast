package comm

import (
	"log"
	"net"
)

// Server models a server that listens on IP:Port for UDP packets
type Server struct {
	Port int
	IP   net.IP

	Count int
}

// Start starts the server which then listens for incoming UDP connections
func (s Server) Start() {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{IP: s.IP, Port: s.Port})
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, addr, _ := conn.ReadFromUDP(buf)
		s.Count++
		log.Printf("Received %s from %v. Got %d messages.", string(buf[0:n]), addr, s.Count)
	}

}
