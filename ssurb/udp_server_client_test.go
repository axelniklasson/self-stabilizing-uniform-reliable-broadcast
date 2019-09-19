package ssurb

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
	"gotest.tools/assert"
)

const PORT = 8080

var IP = []byte{0, 0, 0, 0}

var r = MockResolver{Modules: make(map[ModuleType]interface{})}

func constructClient() (*net.Conn, error) {
	ipString := ""
	for idx, x := range IP {
		if idx == len(IP)-1 {
			ipString += fmt.Sprintf("%v", x)
		} else {
			ipString += fmt.Sprintf("%v.", x)
		}
	}

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ipString, PORT))
	return &conn, err
}

func TestSend(t *testing.T) {
	// first check that server can be created and started
	server := &Server{IP: IP, Port: PORT, Resolver: &r}
	go func(t *testing.T, s *Server) {
		err := server.Start()
		assert.NilError(t, err)
		server.Listen()
	}(t, server)

	// then check that it is possible to connect to started server
	clientConn, err := constructClient()
	assert.NilError(t, err)
	assert.Assert(t, clientConn != nil)
	c := *clientConn
	c.Close()

	// finally check that it is possible to send message
	addr := &net.UDPAddr{IP: IP, Port: PORT}
	msg := models.Message{Type: models.MSG, Sender: 0, Data: map[string]interface{}{"foo": "bar"}}
	err = send(addr, &msg)
	assert.NilError(t, err)
	send(addr, &msg)
	assert.NilError(t, err)
	send(addr, &msg)
	assert.NilError(t, err)
	send(addr, &msg)
	assert.NilError(t, err)

	messagesDelivered := false
	tries := 0
	for !messagesDelivered && tries < 5 {
		messagesDelivered = server.Count == 4
		if !messagesDelivered {
			log.Println("Messages not delivered yet, sleeping 2s..")
			time.Sleep(2 * time.Second)
		}
		tries++
	}
	assert.Assert(t, messagesDelivered)
}
