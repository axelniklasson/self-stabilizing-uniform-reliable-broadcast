package comm

import "net"

// Send is used to send payload over UDP to a destIP:destPort
func Send(destIP net.IP, destPort int, payload []byte) {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: destIP, Port: destPort})
	defer conn.Close()

	conn.Write(payload)
}
