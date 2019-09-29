package ssurb

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

type clientMetrics struct {
	ConnError      string
	PackError      string
	WriteError     string
	FatalSendError string

	ErrorCount *prometheus.CounterVec
	MsgCount   *prometheus.CounterVec
}

var metrics = &clientMetrics{
	ConnError:      "conn_error",
	PackError:      "pack_error",
	WriteError:     "write_error",
	FatalSendError: "fatal_send_error",

	ErrorCount: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "udp_client_error_count",
		Help: "The amount of errors emitted by the udp client",
	}, []string{"error_type", "receiver_id"}),
	MsgCount: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "udp_client_msg_count",
		Help: "The amount messages sent by this client",
	}, []string{"receiver_id"}),
}

// SendToProcessor is a wrapper around send, intended to be called from modules
func SendToProcessor(receiverID int, msg *models.Message) {
	// don't send messages during unit testing
	if helpers.IsUnitTesting() {
		return
	}

	addr := net.UDPAddr{IP: helpers.Processors[receiverID].IP, Port: 4000 + receiverID}

	// try for a maximum of ten times to send packet
	tries := 0
	sent := false
	for !sent && tries < 10 {
		tries++
		err := send(&addr, msg, receiverID)
		if err != nil {
			log.Printf("Got error when sending %v to %d: %v, retrying..", msg, receiverID, err)
			time.Sleep(time.Millisecond * 10)
		} else {
			metrics.MsgCount.WithLabelValues(strconv.Itoa(receiverID)).Inc()
			sent = true
		}
	}

	if !sent {
		log.Printf("Fatal error when sending %v to %d, not re-trying..", msg, receiverID)
		metrics.ErrorCount.WithLabelValues(metrics.FatalSendError, strconv.Itoa(receiverID)).Inc()
	}
}

// Send is used to send payload over UDP to a destIP:destPort
func send(addr *net.UDPAddr, msg *models.Message, receiverID int) error {
	// construct connection to server
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		metrics.ErrorCount.WithLabelValues(metrics.ConnError, strconv.Itoa(receiverID)).Inc()
		return err
	}

	// prepare payload
	payload, err := helpers.Pack(msg)
	if err != nil {
		metrics.ErrorCount.WithLabelValues(metrics.PackError, strconv.Itoa(receiverID)).Inc()
		return err
	}

	// write payload over socket
	_, err = conn.Write(payload)
	if err != nil {
		metrics.ErrorCount.WithLabelValues(metrics.WriteError, strconv.Itoa(receiverID)).Inc()
		return err
	}
	conn.Close()

	return nil
}
