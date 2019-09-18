package modules

import (
	"log"
	"net"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// Send is used to send payload over UDP to a destIP:destPort
func Send(receiverID int, msg *models.Message) {
	// don't send messages during unit testing
	if helpers.IsUnitTesting() {
		return
	}

	// construct connection to server
	addr := net.UDPAddr{IP: helpers.Processors[receiverID].IP, Port: 4000 + receiverID}
	conn, err := net.DialUDP("udp", nil, &addr)

	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	// prepare payload
	payload, err := helpers.Pack(msg)
	if err != nil {
		log.Fatal(err)
	}

	// write payload over socket
	_, err = conn.Write(payload)
	if err != nil {
		log.Fatal(err)
	}
}
