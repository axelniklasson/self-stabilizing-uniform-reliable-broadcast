package modules

import (
	"log"
	"net"
	"os"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// Send is used to send payload over UDP to a destIP:destPort
func Send(receiverID int, msg *models.Message) {
	if _, isSet := os.LookupEnv("TESTING"); isSet {
		return
	}

	addr := net.UDPAddr{IP: helpers.Processors[receiverID].IP, Port: 4000 + receiverID}
	conn, err := net.DialUDP("udp", nil, &addr)

	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}

	payload, err := helpers.Pack(msg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(payload)
	if err != nil {
		log.Fatal(err)
	}
}
