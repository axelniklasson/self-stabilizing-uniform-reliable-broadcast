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
		txObsS = append(rxObsS, -1)
	}
	r := MockResolver{Modules: make(map[ModuleType]interface{})}
	urbModule := UrbModule{ID: 0, P: P, Resolver: &r, Seq: seq, Buffer: buffer, RxObsS: rxObsS, TxObsS: txObsS}
	thetaModule := ThetafdModule{ID: 0, P: P, Resolver: &r, Vector: zeroedSlice}

	r.Modules[URB] = urbModule
	r.Modules[THETAFD] = thetaModule

	return &urbModule, &r
}

func TestObsolete(t *testing.T) {
	mod, resolver := bootstrap()
	mod.RxObsS = []int{0, 1, 0, 0, 0, 0}

	// construct record that is considered to be obsolete
	resolver.TrustedRet = []int{0, 1, 2}
	r := BufferRecord{Identifier: Identifier{ID: 1, Seq: 2}, Delivered: true, RecBy: map[int]bool{0: true, 1: true, 2: true, 3: true}}
	assert.Assert(t, mod.obsolete(r))

	// testing delivered == false returns not obsolete
	r.Delivered = false
	assert.Assert(t, !mod.obsolete(r))
	r.Delivered = true

	// testing that trusted not subset of recBy returns not obsolete
	r.RecBy = map[int]bool{0: true, 1: true}
	assert.Assert(t, !mod.obsolete(r))
	r.RecBy = map[int]bool{0: true, 1: true, 2: true, 3: true}

	// testing that r.seqnum != last highest obsolete seqnum for r.id returns not obsolete
	r.Identifier.Seq = 5 // 5 != 1 + 1
	assert.Assert(t, mod.obsolete(r))
}

func TestMaxSeq(t *testing.T) {
	mod, _ := bootstrap()
	k := 1

	// no record with id = k in buffer since it's empty, should return -1
	assert.Assert(t, mod.maxSeq(k) == -1)

	// add a few records to the buffer, should return highest seq num for id == k
	mod.Buffer.Add(BufferRecord{Msg: nil, Identifier: Identifier{ID: k, Seq: 0}})
	mod.Buffer.Add(BufferRecord{Msg: nil, Identifier: Identifier{ID: k, Seq: 1}})
	mod.Buffer.Add(BufferRecord{Msg: nil, Identifier: Identifier{ID: k, Seq: 2}})
	mod.Buffer.Add(BufferRecord{Msg: nil, Identifier: Identifier{ID: k + 1, Seq: 3}})

	assert.Equal(t, mod.maxSeq(k), 2)
	assert.Equal(t, mod.maxSeq(k+1), 3)
}

func TestMinTxObsS(t *testing.T) {
	mod, resolver := bootstrap()
	mod.TxObsS = []int{1, 2, 5, 10, 0, 50}
	resolver.TrustedRet = []int{1, 3, 5}

	// should return 2, since mod.TxObsS[1] is smallest value for x, x part of resolver.TrustedRet
	assert.Equal(t, mod.minTxObsS(), 2)
}

func TestUpdate(t *testing.T) {
	mod, _ := bootstrap()

	// populate buffer with a few records
	mod.Buffer.Add(BufferRecord{Identifier: Identifier{ID: 1, Seq: 0}, RecBy: map[int]bool{0: true, 1: true}})
	mod.Buffer.Add(BufferRecord{Identifier: Identifier{ID: 2, Seq: 0}, RecBy: map[int]bool{0: true, 2: true}})
	mod.Buffer.Add(BufferRecord{Identifier: Identifier{ID: 2, Seq: 1}, RecBy: map[int]bool{0: true, 2: true}})

	assert.Equal(t, len(mod.Buffer.Records), 3)
	// trying to update with a nil message with new identifier should result in buffer being unchanged
	mod.update(nil, 3, 1, 3)
	assert.Equal(t, len(mod.Buffer.Records), 3)
	// updating with proper message with new identifier should add to buffer
	mod.update(&Message{Contents: []byte("asd")}, 3, 0, 3)
	assert.Equal(t, len(mod.Buffer.Records), 4)

	// trying to update with a message whose identifier already exists should simply add j and k to recBy
	mod.update(&Message{Contents: []byte("asd")}, 1, 0, 5)
	m := mod.Buffer.Get(Identifier{ID: 1, Seq: 0}).RecBy
	m2 := map[int]bool{0: true, 1: true, 5: true}
	assert.Assert(t, reflect.DeepEqual(m, m2))
}

func TestUrbBroadcast(t *testing.T) {

}

func TestUrbDeliver(t *testing.T) {

}
