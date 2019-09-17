package comm

import "net"

func Send(destIP net.IP, destPort int, payload []byte) {
	conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: destIP, Port: destPort})
	defer conn.Close()

	conn.Write(payload)
}
