package ssurb

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
	Resolver IResolver
	Conn     *net.UDPConn

	Count int
}

// Start starts the server and binds it to IP:PORT
func (s *Server) Start() error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: s.IP, Port: s.Port})
	if err != nil {
		return err
	}

	log.Printf("Server listening on %s\n", conn.LocalAddr().String())
	s.Conn = conn
	s.Count = 0
	return nil
}

// Listen tells the server to start listening for packets on IP:PORT
func (s *Server) Listen() error {
	defer s.Conn.Close()
	buf := make([]byte, constants.ServerBufferSize)

	for {
		n, _, err := s.Conn.ReadFromUDP(buf)

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
