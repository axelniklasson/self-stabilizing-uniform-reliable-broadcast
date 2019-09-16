package modules

import (
	"log"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
)

// HbfdModule models a HB failure detector
type HbfdModule struct {
	ID       int
	P        []int
	Resolver *Resolver
	Hb       []int
}

// HB returns the current value of the hb failure detector
func (m *HbfdModule) HB() []int {
	return m.Hb
}

// DoForever starts the algorithm and runs forever
func (m *HbfdModule) DoForever() {
	for {
		for _, id := range m.P {
			if id != m.ID {
				m.sendHeartbeat(id)
			}
		}

		time.Sleep(time.Second * constants.ModuleRunSleepSeconds)
		log.Printf("One iteration of doForever() done")
	}
}

// onHearbeat is called by the resolver when a new heartbeat message was received from another processor
func (m *HbfdModule) onHeartbeat(senderID int) {
	m.Hb[senderID]++
}

// sendHeartbeat sends a heartbeat to another processor to indicate that this processor is alive
func (m *HbfdModule) sendHeartbeat(receiverID int) {
	return
}
