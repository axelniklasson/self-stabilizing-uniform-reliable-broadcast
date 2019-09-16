package modules

import (
	"log"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/constants"
)

// bufferUnitSize is used to control the number of messages allowed to be in the buffer for a processor
const bufferUnitSize = 10

// UrbModule models the URB algorithm in the paper
type UrbModule struct {
	ID       int
	P        []int
	Resolver IResolver
	Seq      int
	Buffer   Buffer
	RxObsS   []int
	TxObsS   []int
}

// MessageType indicates the type of message
type MessageType int

const (
	// MSG represents a broadcasted message
	MSG MessageType = 0
	// MSGack represents an acknowledgement of a broadcasted message
	MSGack MessageType = 1
	// GOSSIP represents messages used by processors to update each other
	GOSSIP MessageType = 2
)

func (m *UrbModule) obsolete(r BufferRecord) bool {
	// return false if trusted is not subset of r.RecBy
	for _, id := range m.Resolver.trusted() {
		if _, exists := r.RecBy[id]; !exists {
			return false
		}
	}

	return m.RxObsS[r.Identifier.ID]+1 == r.Identifier.Seq && r.Delivered
}

func (m *UrbModule) maxSeq(k int) int {
	max := -1
	for _, record := range m.Buffer.Records {
		if record.Identifier.ID == k && record.Identifier.Seq > max {
			max = record.Identifier.Seq
		}
	}

	return max
}

func (m *UrbModule) minTxObsS() int {
	trusted := m.Resolver.trusted()
	min := -1
	for _, id := range trusted {
		if min == -1 || m.TxObsS[id] < min {
			min = m.TxObsS[id]
		}
	}
	return min
}

func (m *UrbModule) update(msg *Message, j int, s int, k int) {
	if s <= m.RxObsS[j] {
		return
	}

	id := Identifier{ID: j, Seq: s}
	r := m.Buffer.Get(id)

	// add record to buffer if new id and message is not nil
	if r == nil && msg != nil {
		recBy := map[int]bool{j: true, k: true}
		prevHB := []int{}
		for i := 0; i < len(m.P); i++ {
			prevHB = append(prevHB, -1)
		}

		newRecord := BufferRecord{Msg: msg, Identifier: id, Delivered: false, RecBy: recBy, PrevHB: prevHB}
		m.Buffer.Add(newRecord)
	} else if r != nil {
		r.RecBy[j] = true
		r.RecBy[k] = true
	}
}

// TODO figure out if this really should spawn a goroutine or rather be wrapped entirely in a goroutine
func (m *UrbModule) urbBroadcast(msg *Message) {
	go func(m *UrbModule) {
		for m.Seq < m.minTxObsS()+bufferUnitSize {
		}

		m.Seq++
		m.update(msg, m.ID, m.Seq, m.ID)
	}(m)
}

func (m *UrbModule) urbDeliver(msg *Message) {
	// TODO something more sofisticated with msg than just logging it
	log.Printf("Delivering message %v\n", msg)
}

// DoForever starts the algorithm and runs forever
func (m *UrbModule) DoForever() {
	for {
		time.Sleep(time.Second * constants.ModuleRunSleepSeconds)
		log.Printf("One iteration of doForever() done")
	}
}
