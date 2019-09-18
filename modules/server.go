package modules

import (
	"log"
	"net"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"
)

// Server models a server that listens on IP:Port for UDP packets
type Server struct {
	Port     int
	IP       net.IP
	Resolver *Resolver

	Count int
}

// Start starts the server which then listens for incoming UDP connections
func (s Server) Start() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: s.IP, Port: s.Port})
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %s\n", conn.LocalAddr().String())

	buf := make([]byte, constants.ServerBufferSize)
	for {
		n, _, err := conn.ReadFromUDP(buf)

		if err != nil {
			log.Fatal(err)
		} else if n > len(buf) {
			log.Fatalf("Got oversized message of size %d, max is %d", n, constants.ServerBufferSize)
		}

		s.Count++
		bytes := buf[0:n]
		msg, err := helpers.Unpack(bytes)
		if err != nil {
			log.Printf("Could not unpack message. Got error: %v\n", err)
		} else {
			s.Resolver.Dispatch(msg)
		}
	}
}
