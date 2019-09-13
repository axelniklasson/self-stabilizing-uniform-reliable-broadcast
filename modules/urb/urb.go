package urb

import (
	"log"
	"time"

	"self-stabilizing-uniform-reliable-broadcast/modules/thetafd"
)

// bufferUnitSize is used to control the number of messages allowed to be in the buffer for a processor
const bufferUnitSize = 10

func (m *Module) obsolete(r BufferRecord) bool {
	// return false if trusted is not subset of r.RecBy
	for _, id := range thetafd.Trusted() {
		if _, exists := r.RecBy[id]; !exists {
			return false
		}
	}

	return m.RxObsS[r.Identifier.ID]+1 == r.Identifier.Seq && r.Delivered
}

func (m *Module) maxSeq(k int) int {
	max := -1
	for _, record := range m.Buffer.Records {
		if record.Identifier.ID == k && record.Identifier.Seq > max {
			max = record.Identifier.Seq
		}
	}

	return max
}

func (m *Module) minTxObsS() int {
	trusted := thetafd.Trusted()
	min := -1
	for _, id := range trusted {
		if min == -1 || m.TxObsS[id] < min {
			min = m.TxObsS[id]
		}
	}
	return min
}

func (m *Module) update(msg *Message, j int, s int, k int) {
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

func (m *Module) urbBroadcast(msg *Message) {
	// wait(seq − min{seqMin[k]}k∈trusted < bufferUnitSize);
	m.Seq++
	m.update(msg, m.ID, m.Seq, m.ID)
}

func (m *Module) urbDeliver(msg *Message) {
	// TODO implement
}

// DoForever starts the algorithm and runs forever
func (m *Module) DoForever() {
	for {
		time.Sleep(time.Second * 1)
		log.Printf("One iteration of doForever() done")
	}
}
