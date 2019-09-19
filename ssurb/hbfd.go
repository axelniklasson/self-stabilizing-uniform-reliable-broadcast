package ssurb

import (
	"log"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
)

// HbfdModule models a HB failure detector
type HbfdModule struct {
	ID       int
	P        []int
	Resolver IResolver

	Hb []int
}

// Init initializes the hbfd module
func (m *HbfdModule) Init() {
	for i := 0; i < len(m.P); i++ {
		m.Hb = append(m.Hb, 0)
	}
}

// HB returns the current value of the hb failure detector
func (m *HbfdModule) HB() []int {
	return m.Hb
}

// DoForever starts the algorithm and runs forever
func (m *HbfdModule) DoForever() {
	log.Printf("DoForever() starting")

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
func (m *HbfdModule) onHeartbeat(senderID int) {
	m.Hb[senderID]++
}

// sendHeartbeat sends a heartbeat to another processor to indicate that this processor is alive
func (m *HbfdModule) sendHeartbeat(receiverID int) {
	return
}
