package modules

import (
	"log"
	"self-stabilizing-uniform-reliable-broadcast/constants"
	"time"
)

// ThetafdModule models a theta failure detector
type ThetafdModule struct {
	ID       int
	P        []int
	Resolver IResolver
	Vector   []int
}

// DoForever starts the algorithm and runs forever
func (m *ThetafdModule) DoForever() {
	for {
		for _, id := range m.P {
			if id != m.ID {
				m.sendHeartbeat(id)
			}
		}

		time.Sleep(time.Second * constants.MODULE_RUN_SLEEP_SECONDS)
		log.Printf("One iteration of doForever() done")
	}
}

// onHearbeat is called by the resolver when a new heartbeat message was received from another processor
func (m *ThetafdModule) onHeartbeat(senderID int) {
	m.Vector[senderID] = 0
	for idx := range m.Vector {
		if idx == senderID || idx == m.ID {
			m.Vector[idx] = 0
		} else {
			m.Vector[idx]++
		}
	}
}

// sendHeartbeat sends a heartbeat to another processor to indicate that this processor is alive
func (m *ThetafdModule) sendHeartbeat(receiverID int) {
	return
}
