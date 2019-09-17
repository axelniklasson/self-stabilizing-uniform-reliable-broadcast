package modules

import (
	"reflect"
	"testing"

	"gotest.tools/assert"
)

type MockResolver struct {
	Modules    map[ModuleType]interface{}
	TrustedRet []int
	HbRet      []int
}

func (r *MockResolver) hb() []int      { return r.HbRet }
func (r *MockResolver) trusted() []int { return r.TrustedRet }

func bootstrap() (*UrbModule, *MockResolver) {
	P := []int{0, 1, 2, 3, 4, 5}
	seq := 0
	buffer := Buffer{}
	rxObsS := []int{}
	txObsS := []int{}
	zeroedSlice := []int{}
	for i := 0; i < len(P); i++ {
		rxObsS = append(rxObsS, -1)
		txObsS = append(txObsS, -1)
	}
	r := MockResolver{Modules: make(map[ModuleType]interface{})}
	urbModule := UrbModule{ID: 0, P: P, Resolver: &r, Seq: seq, Buffer: &buffer, RxObsS: rxObsS, TxObsS: txObsS}
	thetaModule := ThetafdModule{ID: 0, P: P, Resolver: &r, Vector: zeroedSlice}
	hbfdModule := HbfdModule{ID: 0, P: P, Resolver: &r, Hb: zeroedSlice}

	r.Modules[URB] = urbModule
	r.Modules[THETAFD] = thetaModule
	r.Modules[HBFD] = hbfdModule

	return &urbModule, &r
}

func TestInit(t *testing.T) {
	mod := UrbModule{ID: 0, P: []int{0, 1, 2}}
	mod.Init()
	assert.Equal(t, mod.Seq, 0)
	assert.Equal(t, len(mod.Buffer.Records), 0)
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{-1, -1, -1}))
	assert.Assert(t, reflect.DeepEqual(mod.RxObsS, []int{-1, -1, -1}))
}

func TestFlushBufferIfStaleInfo(t *testing.T) {
	mod, _ := bootstrap()

	// add two records without stale info, should not be flushed
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 1, Seq: 0}})
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 2, Seq: 0}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.flushBufferIfStaleInfo()
	assert.Equal(t, len(mod.Buffer.Records), 2)

	// add one record with empty msg, buffer should be flushed
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 1}})
	assert.Equal(t, len(mod.Buffer.Records), 3)
	mod.flushBufferIfStaleInfo()
	assert.Equal(t, len(mod.Buffer.Records), 0)

	// add two non-stale records and one with a duplicate identifier, buffer should be flushed
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 1, Seq: 0}})
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 2, Seq: 0}})
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 2, Seq: 0}})
	assert.Equal(t, len(mod.Buffer.Records), 3)
	mod.flushBufferIfStaleInfo()
	assert.Equal(t, len(mod.Buffer.Records), 0)

}

func TestCheckTransmitWindow(t *testing.T) {
	mod, resolver := bootstrap()
	resolver.TrustedRet = []int{0, 1, 2}
	initTx := []int{10, 5, 2, -1, -1, 1}
	mod.TxObsS = initTx

	// mod.TxObsS should get mod.Seq in all values since mod.Seq (1) < minTxObs (2)
	mod.Seq = 1
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, initTx))
	mod.checkTransmitWindow()
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{1, 1, 1, 1, 1, 1}))

	// mod.TxObsS should get mod.Seq in all values since mod.Seq (20) > minTxObs (2) + bufferUnitSize (10) = 12
	mod.Seq = 12
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{1, 1, 1, 1, 1, 1}))
	mod.checkTransmitWindow()
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{12, 12, 12, 12, 12, 12}))

	// mod.TxObsS should get mod.Seq in all values since set of allowed seqnums not subset of msg seqnums sent by this processor
	// set of allowed seqnums = {3, 4, 5}
	// set of msg seqnums from this processor = {4, 5}
	mod.TxObsS = initTx
	mod.Seq = 5
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 4}})
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 5}})
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, initTx))
	mod.checkTransmitWindow()
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{5, 5, 5, 5, 5, 5}))

	// mod.TxObsS should remain unchanged since all conditions hold up
	// set of allowed seqnums = {3, 4, 5}
	// set of msg seqnums from this processor = {3, 4, 5}
	mod.TxObsS = initTx
	mod.Seq = 5
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 3}})
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 4}})
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 5}})
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, initTx))
	mod.checkTransmitWindow()
	assert.Assert(t, reflect.DeepEqual(mod.TxObsS, []int{5, 5, 5, 5, 5, 5}))

}

func TestCheckReceivingWindow(t *testing.T) {
	mod, _ := bootstrap()
	mod.RxObsS[1] = 0
	mod.RxObsS[2] = 10

	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 1, Seq: 20}})
	mod.Buffer.Add(&BufferRecord{Msg: &Message{}, Identifier: Identifier{ID: 2, Seq: 0}})
	assert.Equal(t, mod.RxObsS[1], 0)
	assert.Equal(t, mod.RxObsS[2], 10)
	// should choose 20 (Seq) - 10 (bufferUnitSize) for RxObsS[1] and 10 (RxObsS[2]) for RxObsS[2]
	mod.checkReceivingWindow()
	assert.Equal(t, mod.RxObsS[1], 10)
	assert.Equal(t, mod.RxObsS[2], 10)
}

func TestUpdateReceiverCounters(t *testing.T) {
	mod, resolver := bootstrap()
	mod.TxObsS[1] = 0
	resolver.TrustedRet = []int{0, 1}

	// add one obselete and one non-obsolete record to buffer
	mod.Buffer.Add(&BufferRecord{Delivered: true, Identifier: Identifier{ID: 1, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})
	mod.Buffer.Add(&BufferRecord{Delivered: false, Identifier: Identifier{ID: 2, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})

	// should increment mod.RxObs[1] by one since that is the only obsolete record in buffer
	assert.Equal(t, mod.RxObsS[1], -1)
	assert.Equal(t, mod.RxObsS[2], -1)
	mod.updateReceiverCounters()
	assert.Equal(t, mod.RxObsS[1], 0)
	assert.Equal(t, mod.RxObsS[2], -1)
}

func TestTrimBuffer(t *testing.T) {
	mod, resolver := bootstrap()
	resolver.TrustedRet = []int{0, 1, 2}
	mod.TxObsS[0] = 5
	mod.TxObsS[1] = 7
	mod.TxObsS[2] = 2

	// add two records constructed at this processor, one with seq < minTxObsS() and one not
	// should only keep the first record since its minTxObs < its seqnum
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 13}, RecBy: map[int]bool{0: true, 1: true}})
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.trimBuffer()
	assert.Equal(t, len(mod.Buffer.Records), 1)
	assert.Equal(t, mod.Buffer.Records[0].Identifier, Identifier{ID: 0, Seq: 13})

	// add record with processor not part of P, should be removed
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 20, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.trimBuffer()
	assert.Equal(t, len(mod.Buffer.Records), 1)

	// add record with seqnum not > mod.rxObs[k], k = record.id
	mod.RxObsS[1] = 5
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.trimBuffer()
	assert.Equal(t, len(mod.Buffer.Records), 1)

	// add record which seqnum is < maxSeq(k) - bufferUnitsize
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 2}, RecBy: map[int]bool{0: true, 1: true}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.trimBuffer()
	assert.Equal(t, len(mod.Buffer.Records), 1)

	// add record from other processor which should be kept in buffer
	mod.Buffer.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 8}, RecBy: map[int]bool{0: true, 1: true}})
	assert.Equal(t, len(mod.Buffer.Records), 2)
	mod.trimBuffer()
	assert.Equal(t, len(mod.Buffer.Records), 2)

}

func TestProcessMessages(t *testing.T) {
	mod, resolver := bootstrap()
	resolver.TrustedRet = []int{0, 1}

	// add one record to deliver and one record that should not be delivered due to trusted not subset of recby
	buf := mod.Buffer
	buf.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 0}, RecBy: map[int]bool{0: true}})
	// mod.Buffer.Add()
	buf.Add(&BufferRecord{Identifier: Identifier{ID: 1, Seq: 1}, RecBy: map[int]bool{0: true, 1: true, 2: true}})
	assert.Assert(t, !mod.Buffer.Records[0].Delivered)
	assert.Assert(t, !mod.Buffer.Records[1].Delivered)
	mod.processMessages()
	assert.Assert(t, !mod.Buffer.Records[0].Delivered)
	assert.Assert(t, mod.Buffer.Records[1].Delivered)

	// adding processor 1 to recBy makes trusted subset of recBy, should now be delivered
	mod.Buffer.Records[0].RecBy[1] = true
	mod.processMessages()
	assert.Assert(t, mod.Buffer.Records[0].Delivered)
}

func TestGossip(t *testing.T) {
	// TODO
}

func TestHasObsoleteRecord(t *testing.T) {
	mod, resolver := bootstrap()
	mod.RxObsS[1] = 1
	resolver.TrustedRet = []int{0, 1}

	// buffer empty, should return nil
	assert.Assert(t, mod.hasObsoleteRecord() == nil)

	// add one record to buffer that has Delivered = false, i.e. not obsolete
	mod.Buffer.Add(&BufferRecord{Delivered: false})
	assert.Assert(t, mod.hasObsoleteRecord() == nil)

	// add another record to buffer that is obsolete, make sure that record is returned by hasObsoleteRecord
	r := BufferRecord{Delivered: true, Identifier: Identifier{ID: 1, Seq: 2}, RecBy: map[int]bool{0: true, 1: true}}
	mod.Buffer.Add(&r)
	assert.Assert(t, mod.hasObsoleteRecord() != nil)
	assert.Equal(t, mod.hasObsoleteRecord().Identifier, Identifier{ID: 1, Seq: 2})
}
