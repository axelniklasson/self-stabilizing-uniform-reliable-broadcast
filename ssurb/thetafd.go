package ssurb

import (
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// ThetafdModule models a theta failure detector
type ThetafdModule struct {
	ID       int
	P        []int
	Resolver IResolver

	Vector []int
}

// Init initializes the thetafd module
func (m *ThetafdModule) Init() {
	for i := 0; i < len(m.P); i++ {
		m.Vector = append(m.Vector, 0)
	}
}

// Trusted returns the set of processor IDs that are below the threshold ThetafdW
func (m *ThetafdModule) Trusted() []int {
	trusted := []int{}
	for idx, x := range m.Vector {
		if x < constants.ThetafdW && idx != m.ID {
			trusted = append(trusted, idx)
		}
	}

	return trusted
}

// DoForever starts the algorithm and runs forever
func (m *ThetafdModule) DoForever() {
	for {
		for _, id := range m.P {
			if id != m.ID {
				m.sendHeartbeat(id)
			}
		}

		time.Sleep(time.Second * constants.ModuleRunSleepSeconds)
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
	message := models.Message{Type: models.THETAheartbeat, Sender: m.ID, Data: nil}
	go SendToProcessor(receiverID, &message)
}
