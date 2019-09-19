package ssurb

import (
	"log"
	"net"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// SendToProcessor is a wrapper around send, intended to be called from modules
func SendToProcessor(receiverID int, msg *models.Message) {
	// don't send messages during unit testing
	if helpers.IsUnitTesting() {
		return
	}

	addr := net.UDPAddr{IP: helpers.Processors[receiverID].IP, Port: 4000 + receiverID}
	err := send(&addr, msg)
	if err != nil {
		log.Printf("Got error when sending %v to %d: %v", msg, receiverID, err)
	}
}

// Send is used to send payload over UDP to a destIP:destPort
func send(addr *net.UDPAddr, msg *models.Message) error {
	// construct connection to server
	conn, err := net.DialUDP("udp", nil, addr)

	defer conn.Close()
	if err != nil {
		return err
	}

	// prepare payload
	payload, err := helpers.Pack(msg)
	if err != nil {
		return err
	}

	// write payload over socket
	_, err = conn.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
