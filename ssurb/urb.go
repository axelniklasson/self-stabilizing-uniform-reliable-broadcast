package ssurb

import (
	"log"
	"sync"
	"time"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/helpers"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var mux sync.Mutex

// UrbMessage is the type of the actual message that is sent from the app
type UrbMessage struct {
	Text string
}

type urbMetrics struct {
	// General
	BroadcastedMessagesCount prometheus.Counter
	DeliveredMessagesCount   prometheus.Counter
	DeliveredByteCount       prometheus.Counter

	// Throughput
	MessageLatency prometheus.Histogram
}

// UrbModule models the URB algorithm in the paper
type UrbModule struct {
	ID       int
	P        []int
	Resolver IResolver

	Seq    int
	Buffer *Buffer
	RxObsS []int
	TxObsS []int

	// Metrics stuff
	Metrics         *urbMetrics
	PendingMessages map[*UrbMessage]time.Time
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

	// init metrics
	m.Metrics = &urbMetrics{
		BroadcastedMessagesCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "urb_broadcasted_messages_count",
			Help: "The total number of broadcasted messages",
		}),
		DeliveredMessagesCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "urb_delivered_messages_count",
			Help: "The total number of delivered messages",
		}),
		DeliveredByteCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "urb_delivered_bytes_count",
			Help: "The total number of delivered bytes",
		}),
		MessageLatency: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "urb_message_latency",
				Help:    "Message latency (ms)",
				Buckets: []float64{50, 100, 250, 500, 1000, 10000},
			},
		),
	}
	m.PendingMessages = map[*UrbMessage]time.Time{}
	prometheus.MustRegister(m.Metrics.MessageLatency)
}

// obsolete is used to determine what records in the buffer are considered to be obsolete
// a record r is obsolete if it is delivered, its seqnum is the next in line to be obsolete and
// all trusted processors have acked the record
func (m *UrbModule) obsolete(r *BufferRecord) bool {
	// return false if trusted is not subset of r.RecBy
	for _, id := range m.Resolver.Trusted() {
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

// BlockUntilAvailableSpace busy-waits until flow control mechanism ensures enough space on all trusted receivers
func (m *UrbModule) BlockUntilAvailableSpace() {
	for m.Seq >= m.minTxObsS()+helpers.GetBufferUnitSize() {
		time.Sleep(time.Millisecond * 1)
	}
}

// minTxObsS returns the smallest obsolete sequence number that pi had received from a trusted receiver
func (m *UrbModule) minTxObsS() int {
	trusted := m.Resolver.Trusted()
	min := -1
	for _, id := range trusted {
		if min == -1 || m.TxObsS[id] < min {
			min = m.TxObsS[id]
		}
	}

	if m.TxObsS[m.ID] < min {
		min = m.TxObsS[m.ID]
	}

	return min
}

// update processes a message through creating a unique operation index and adding it to buffer if it's a new message.
// Otherwise it adds processors j and k to recBy of the existing record
func (m *UrbModule) update(msg *UrbMessage, j int, s int, k int) {
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

		newRecord := &BufferRecord{Msg: msg, Identifier: id, Delivered: false, RecBy: recBy, PrevHB: prevHB}
		m.Buffer.Add(newRecord)
	} else if r != nil {
		r.RecBy[j] = true
		r.RecBy[k] = true
	}
}

// UrbBroadcast is called from the application layer to broadcast a message
// NOTE call this in a separate goroutine
func (m *UrbModule) UrbBroadcast(msg *UrbMessage) {

	// grab lock
	mux.Lock()

	for m.Seq >= m.minTxObsS()+helpers.GetBufferUnitSize() {
		mux.Unlock()
		time.Sleep(time.Microsecond * 100)
		mux.Lock()
	}

	m.Seq++
	m.update(msg, m.ID, m.Seq, m.ID)

	// TODO use NTP time
	// ts := helpers.GetNTPTime().UnixNano()

	// ts := time.Now().UnixNano()
	m.PendingMessages[msg] = time.Now()

	// release lock
	mux.Unlock()

	// emit metric that msg was broadcasted
	m.Metrics.BroadcastedMessagesCount.Inc()
}

// UrbDeliver delivers a message to the application layer
func (m *UrbModule) UrbDeliver(msg *UrbMessage, id int) {
	if !helpers.IsUnitTesting() && m.Metrics != nil {
		if t1, exists := m.PendingMessages[msg]; exists {
			// TODO use NTP time
			latency := time.Now().Sub(t1)
			m.Metrics.MessageLatency.Observe(float64(latency.Milliseconds()))
		}

		m.Metrics.DeliveredMessagesCount.Inc()
		m.Metrics.DeliveredByteCount.Add(float64(len(msg.Text)))
	}
}

// DoForever starts the algorithm and runs forever
func (m *UrbModule) DoForever() {
	for {
		// retrieve lock
		mux.Lock()

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

		// release lock
		mux.Unlock()

		time.Sleep(helpers.GetModRunSleepDuration())
	}
}

// flushBufferIfStaleInfo flushes the buffer whenever records with msg == nil or two (or more) records with same msg identifier
func (m *UrbModule) flushBufferIfStaleInfo() {
	identifiers := map[Identifier]bool{}
	flush := false

	// lines 18-19
	for _, r := range m.Buffer.Records {

		// if empty message found, abort and flush
		if r.Msg == nil {
			log.Println("flushing buffer due to empty msg found")
			flush = true
			break
		}

		// if multiple record identifiers are found, abort and flush
		if _, exists := identifiers[r.Identifier]; exists {
			log.Printf("flushing buffer due to duplicate message identifier %v", r.Identifier)
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
	seqBound := mS <= m.Seq && m.Seq <= (mS+helpers.GetBufferUnitSize())
	subSet := isSubset(s, s2)
	if !(seqBound && subSet) {
		if !seqBound {
			log.Printf("setting all values in TxObsS to %d due to m.seq not being between mS (%d) and mS+bufferUnitSize (%d)", m.Seq, mS, mS+helpers.GetBufferUnitSize())
		} else if !subSet {
			log.Printf("setting all values in TxObsS to %d due to seqnums ms+1..seq not in buffer", m.Seq)
		}

		for idx := range m.TxObsS {
			m.TxObsS[idx] = m.Seq
		}
	}
}

// checkReceivingWindow makes sure the gap between the largest obsolete record and largest buffered sequence number
// is not larger than bufferUnitSize
func (m *UrbModule) checkReceivingWindow() {
	for _, k := range m.P {
		m.RxObsS[k] = max(m.RxObsS[k], m.maxSeq(k)-helpers.GetBufferUnitSize())
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
			if contains(m.P, k) && m.RxObsS[k] < s && m.maxSeq(k)-helpers.GetBufferUnitSize() <= s {
				newBuffer.Add(r)
			}
		}
	}

	m.Buffer = &newBuffer
}

// processMessages delivers messages when acks from all trusted processors are present before sampling hb fd (used for re-transmission)
func (m *UrbModule) processMessages() {
	trusted := listToMap(m.Resolver.Trusted())
	for _, r := range m.Buffer.Records {
		if !r.Delivered && isSubset(trusted, r.RecBy) {
			m.UrbDeliver(r.Msg, r.Identifier.ID)
		}
		r.Delivered = r.Delivered || isSubset(trusted, r.RecBy)

		u := m.Resolver.Hb()
		for _, k := range m.P {
			if _, exists := r.RecBy[k]; !exists || (r.Identifier.ID == m.ID && r.Identifier.Seq == m.TxObsS[k]+1) && r.PrevHB[k] < u[k] {
				r.PrevHB = u
				m.sendMSG(k, r.Msg, r.Identifier.ID, r.Identifier.Seq)
			}
		}
	}
}

// gossip sends control info about max seq that pi stores for pk as well as info about max obsolete record for pk
func (m *UrbModule) gossip() {
	for _, k := range m.P {
		m.sendGOSSIP(k, m.maxSeq(k), m.RxObsS[k], m.TxObsS[k])
	}
}

// --- communication methods ---

func (m *UrbModule) sendMSG(receiverID int, msg *UrbMessage, j int, s int) {
	data := map[string]interface{}{
		"msgText": msg.Text,
		"j":       j,
		"s":       s,
	}

	message := models.Message{Type: models.MSG, Sender: m.ID, Data: data}
	go SendToProcessor(receiverID, &message)
}

func (m *UrbModule) sendMSGack(receiverID int, j int, s int) {
	data := map[string]interface{}{
		"j": j,
		"s": s,
	}

	message := models.Message{Type: models.MSGack, Sender: m.ID, Data: data}
	go SendToProcessor(receiverID, &message)
}

func (m *UrbModule) sendGOSSIP(receiverID int, seqJ int, txObsSJ int, rxObsSJ int) {
	data := map[string]interface{}{
		"seqJ":    float64(seqJ),
		"txObsSJ": float64(txObsSJ),
		"rxObsSJ": float64(rxObsSJ),
	}

	message := models.Message{Type: models.GOSSIP, Sender: m.ID, Data: data}
	if receiverID == m.ID {
		// deliver directly if sending to self
		go m.onGOSSIP(&message)
	} else {
		go SendToProcessor(receiverID, &message)
	}
}

func (m *UrbModule) onMSG(msg *models.Message) {
	k := msg.Sender
	message := UrbMessage{Text: msg.Data["msgText"].(string)}
	j := int(msg.Data["j"].(float64))
	s := int(msg.Data["s"].(float64))

	mux.Lock()
	m.update(&message, j, s, k)
	mux.Unlock()

	m.sendMSGack(k, j, s)
}

func (m *UrbModule) onMSGack(msg *models.Message) {
	k := msg.Sender
	j := int(msg.Data["j"].(float64))
	s := int(msg.Data["s"].(float64))

	mux.Lock()
	m.update(nil, j, s, k)
	mux.Unlock()
}

func (m *UrbModule) onGOSSIP(msg *models.Message) {
	j := msg.Sender
	seqJ := int(msg.Data["seqJ"].(float64))
	txObsSJ := int(msg.Data["txObsSJ"].(float64))
	rxObsSJ := int(msg.Data["rxObsSJ"].(float64))

	mux.Lock()
	m.Seq = max(seqJ, m.Seq)
	m.TxObsS[j] = max(txObsSJ, m.TxObsS[j])
	m.RxObsS[j] = max(rxObsSJ, m.RxObsS[j])
	mux.Unlock()
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

// GetData returns a dict with urb mod data
func (m *UrbModule) GetData() map[string]interface{} {
	return map[string]interface{}{
		"seq":           m.Seq,
		"bufferRecords": m.Buffer.Records,
		"rxObsS":        m.RxObsS,
		"txObsS":        m.TxObsS,
	}
}
