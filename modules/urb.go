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

	Seq    int
	Buffer *Buffer
	RxObsS []int
	TxObsS []int
}

// Init initializes the urb module
func (m *UrbModule) Init() {
	m.Seq = 0
	m.Buffer = &Buffer{Records: []*BufferRecord{}}
	m.RxObsS = []int{}
	m.TxObsS = []int{}

	for i := 0; i < len(m.P); i++ {
		m.RxObsS = append(m.RxObsS, -1)
		m.TxObsS = append(m.TxObsS, -1)
	}
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

// obsolete is used to determine what records in the buffer are considered to be obsolete
// a record r is obsolete if it is delivered, its seqnum is the next in line to be obsolete and
// all trusted processors have acked the record
func (m *UrbModule) obsolete(r *BufferRecord) bool {
	// return false if trusted is not subset of r.RecBy
	for _, id := range m.Resolver.trusted() {
		if _, exists := r.RecBy[id]; !exists {
			return false
		}
	}

	return m.RxObsS[r.Identifier.ID]+1 == r.Identifier.Seq && r.Delivered
}

// maxSeq returns the highest buffered sequence number for messages sent by processor k.
// should no such exist, -1 is returned
func (m *UrbModule) maxSeq(k int) int {
	max := -1
	for _, record := range m.Buffer.Records {
		if record.Identifier.ID == k && record.Identifier.Seq > max {
			max = record.Identifier.Seq
		}
	}

	return max
}

// minTxObsS returns the smallest obsolete sequence number that pi had received from a trusted receiver
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

// update processes a message through creating a unique operation index and adding it to buffer if it's a new message.
// Otherwise it adds processors j and k to recBy of the existing record
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
		m.Buffer.Add(&newRecord)
	} else if r != nil {
		r.RecBy[j] = true
		r.RecBy[k] = true
	}
}

// UrbBroadcast is called from the application layer to broadcast a message
// TODO figure out if this really should spawn a goroutine or rather be wrapped entirely in a goroutine
func (m *UrbModule) UrbBroadcast(msg *Message) {
	go func(m *UrbModule) {
		for m.Seq < m.minTxObsS()+bufferUnitSize {
		}

		m.Seq++
		m.update(msg, m.ID, m.Seq, m.ID)
	}(m)
}

// UrbDeliver delivers a message to the application layer
func (m *UrbModule) UrbDeliver(msg *Message) {
	// TODO something more sofisticated with msg than just logging it
	log.Printf("Delivering message %v\n", msg)
}

// DoForever starts the algorithm and runs forever
func (m *UrbModule) DoForever() {
	log.Printf("DoForever() starting")

	for {
		// lines 18-19
		m.flushBufferIfStaleInfo()

		// line 20
		m.checkTransmitWindow()

		// line 21
		m.checkReceivingWindow()

		// line 22
		m.updateReceiverCounters()

		// line 23
		m.trimBuffer()

		// lines 24-28
		m.processMessages()

		// line 29
		m.gossip()

		time.Sleep(time.Second * constants.ModuleRunSleepSeconds)
	}
}

// flushBufferIfStaleInfo flushes the buffer whenever records with msg == nil or two (or more) records with same msg identifier
func (m *UrbModule) flushBufferIfStaleInfo() {
	identifiers := map[Identifier]bool{}
	flush := false

	// go through all records
	// lines 18-19
	for _, r := range m.Buffer.Records {

		// if empty message found, abort and flush
		if r.Msg == nil {
			flush = true
			break
		}

		// if multiple record identifiers are found, abort and flush
		if _, exists := identifiers[r.Identifier]; exists {
			flush = true
			break
		} else {
			identifiers[r.Identifier] = true
		}
	}

	if flush {
		m.Buffer.Records = []*BufferRecord{}
	}
}

// checkTransmitWindow checks for transient fault (not all messages seqnums are between mS+1 and seq) and
// recovers through allowing this processor to send bufferUnitSize messages without considering receiving end
func (m *UrbModule) checkTransmitWindow() {
	mS := m.minTxObsS()

	// build set of seqnums {mS+1, ..., Seq}
	s := map[int]bool{}
	for i := mS + 1; i <= m.Seq; i++ {
		s[i] = true
	}
	// build second set of seqnums {s, ..., s'} s.t. id == i
	s2 := map[int]bool{}
	for _, r := range m.Buffer.Records {
		if r.Identifier.ID == m.ID {
			s2[r.Identifier.Seq] = true
		}
	}

	// check if should allow this node to send bufferUnitSize messages without considering receivers
	if !(mS <= m.Seq && m.Seq <= mS+bufferUnitSize && isSubset(s, s2)) {
		for idx := range m.TxObsS {
			m.TxObsS[idx] = m.Seq
		}
	}
}

// checkReceivingWindow makes sure the gap between the largest obsolete record and largest buffered sequence number
// is not larger than bufferUnitSize
func (m *UrbModule) checkReceivingWindow() {
	for _, k := range m.P {
		m.RxObsS[k] = max(m.RxObsS[k], m.maxSeq(k)-bufferUnitSize)
	}
}

// updateReceiverCounters updates the receiver-side counter that stores the highest obsolete message number per sender
func (m *UrbModule) updateReceiverCounters() {
	hasObsolete := true
	for hasObsolete {
		if r := m.hasObsoleteRecord(); r != nil {
			m.RxObsS[r.Identifier.ID]++
		} else {
			hasObsolete = false
		}
	}
}

// trimBuffer makes sure buffer only contains sent messages that are not acked by all trusted or non-obsolete messages
func (m *UrbModule) trimBuffer() {
	newBuffer := Buffer{Records: []*BufferRecord{}}

	for _, r := range m.Buffer.Records {
		if r.Identifier.ID == m.ID {
			if m.minTxObsS() < r.Identifier.Seq {
				newBuffer.Add(r)
			}
		} else {
			k := r.Identifier.ID
			s := r.Identifier.Seq
			if contains(m.P, k) && m.RxObsS[k] < s && m.maxSeq(k)-bufferUnitSize <= s {
				newBuffer.Add(r)
			}
		}
	}

	m.Buffer = &newBuffer
}

// processMessages delivers messages when acks from all trusted processors are present before sampling hb fd (used for re-transmission)
func (m *UrbModule) processMessages() {
	trusted := listToMap(m.Resolver.trusted())
	for _, r := range m.Buffer.Records {
		if !r.Delivered && isSubset(trusted, r.RecBy) {
			m.UrbDeliver(r.Msg)
		}
		r.Delivered = r.Delivered || isSubset(trusted, r.RecBy)

		u := m.Resolver.hb()
		for _, k := range m.P {
			if _, exists := r.RecBy[k]; !exists || (r.Identifier.ID == m.ID && r.Identifier.Seq == m.TxObsS[k]+1) && r.PrevHB[k] < u[k] {
				r.PrevHB = u
				// TODO send MSG(m,j,s) to pk;
			}
		}
	}
}

// gossip sends control info about max seq that pi stores for pk as well as info about max obsolete record for pk
func (m *UrbModule) gossip() {
	for _, k := range m.P {
		if k != m.ID {
			// TODO send GOSSIP(maxSeq(k),rxObsS[k],txObsS[k]) to pk;
		}
	}
}

// --- helper methods ---

// hasObsoleteRecord returns the first found obsolete record, otherwise nil
func (m *UrbModule) hasObsoleteRecord() *BufferRecord {
	for _, r := range m.Buffer.Records {
		if m.obsolete(r) {
			return r
		}
	}

	return nil
}
