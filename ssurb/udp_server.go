package ssurb

import (
	"log"
	"net"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"
)

type serverMetrics struct {
	ListenError   string
	ReadError     string
	OversizeError string
	UnpackError   string

	ErrorCount *prometheus.CounterVec
}

// Server models a server that listens on IP:Port for UDP packets
type Server struct {
	Port     int
	IP       net.IP
	Resolver IResolver
	Conn     *net.UDPConn

	Metrics *serverMetrics
	Count   int
}

// Start starts the server and binds it to IP:PORT
func (s *Server) Start() error {
	s.Metrics = &serverMetrics{
		ListenError:   "listen_error",
		ReadError:     "read_error",
		OversizeError: "oversize_error",
		UnpackError:   "unpack_error",

		ErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "udp_server_error_count",
			Help: "The amount of errors emitted by the udp server",
		}, []string{"error_type"}),
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: s.IP, Port: s.Port})
	if err != nil {
		s.Metrics.ErrorCount.WithLabelValues(s.Metrics.ListenError).Inc()
		return err
	}

	log.Printf("UDP Server listening on %s\n", conn.LocalAddr().String())
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
			s.Metrics.ErrorCount.WithLabelValues(s.Metrics.ReadError).Inc()
			log.Fatal(err)
		} else if n > len(buf) {
			s.Metrics.ErrorCount.WithLabelValues(s.Metrics.OversizeError).Inc()
			log.Fatalf("Got oversized message of size %d, max is %d", n, constants.ServerBufferSize)
		}

		s.Count++
		bytes := buf[0:n]
		msg, err := helpers.Unpack(bytes)
		if err != nil {
			s.Metrics.ErrorCount.WithLabelValues(s.Metrics.UnpackError).Inc()
			log.Printf("Could not unpack message. Got error: %v\n", err)
		} else {
			s.Resolver.Dispatch(msg)
		}
	}
}
