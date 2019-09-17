package modules

import (
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
		txObsS = append(rxObsS, -1)
	}
	r := MockResolver{Modules: make(map[ModuleType]interface{})}
	urbModule := UrbModule{ID: 0, P: P, Resolver: &r, Seq: seq, Buffer: buffer, RxObsS: rxObsS, TxObsS: txObsS}
	thetaModule := ThetafdModule{ID: 0, P: P, Resolver: &r, Vector: zeroedSlice}
	hbfdModule := HbfdModule{ID: 0, P: P, Resolver: &r, Hb: zeroedSlice}

	r.Modules[URB] = urbModule
	r.Modules[THETAFD] = thetaModule
	r.Modules[HBFD] = hbfdModule

	return &urbModule, &r
}

func TestFlushBufferIfStaleInfo(t *testing.T) {
	// TODO
}

func TestCheckTransmitWindow(t *testing.T) {
	// TODO
}

func TestCheckReceivingWindow(t *testing.T) {
	// TODO
}

func TestUpdateReceiverCounters(t *testing.T) {
	// TODO
}

func TestTrimBuffer(t *testing.T) {
	// TODO
}

func TestProcessMessages(t *testing.T) {
	// TODO
}

func TestGossip(t *testing.T) {
	// TODO
}

func TestHasObsoleteRecord(t *testing.T) {
	mod, resolver := bootstrap()
	mod.TxObsS[1] = 1
	resolver.TrustedRet = []int{0, 1}

	// buffer empty, should return nil
	assert.Assert(t, mod.hasObsoleteRecord() == nil)

	// add one record to buffer that has Delivered = false, i.e. not obsolete
	mod.Buffer.Add(BufferRecord{Delivered: false})
	assert.Assert(t, mod.hasObsoleteRecord() == nil)

	// add another record to buffer that is obsolete, make sure that record is returned by hasObsoleteRecord
	r := BufferRecord{Delivered: true, Identifier: Identifier{ID: 1, Seq: 2}, RecBy: map[int]bool{0: true, 1: true}}
	mod.Buffer.Add(r)
	assert.Assert(t, mod.hasObsoleteRecord() != nil)
	assert.Equal(t, mod.hasObsoleteRecord().Identifier, Identifier{ID: 1, Seq: 2})
}
