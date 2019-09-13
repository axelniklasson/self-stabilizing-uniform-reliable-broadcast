package modules

import (
	"testing"

	"gotest.tools/assert"
)

func initModule() UrbModule {
	P := []int{0, 1, 2, 3, 4, 5}
	seq := 0
	buffer := Buffer{}
	rxObsS := []int{}
	txObsS := []int{}
	for i := 0; i < len(P); i++ {
		rxObsS = append(rxObsS, -1)
		txObsS = append(rxObsS, -1)
	}
	return UrbModule{ID: 0, P: P, Seq: seq, Buffer: buffer, RxObsS: rxObsS, TxObsS: txObsS}
}

func TestObsolete(t *testing.T) {
	mod := initModule()
	mod.RxObsS = []int{0, 5, 0, 0, 0, 0}
	// r := BufferRecord{Identifier: Identifier{ID: 1, Seq: 0}, Delivered: false}
	// r2 := BufferRecord{Identifier: Identifier{ID: 1, Seq: 1}, Delivered: true}

	// TODO figure out how to test that it returns record is obsolete based on trusted not being subset of recBy

	// return false if record is not delivered and vice verse
	assert.Assert(t, !mod.obsolete(BufferRecord{Delivered: false}))

	// returns false if seqnum != last highest obsolete number for that node + 1
	assert.Assert(t, !mod.obsolete(BufferRecord{Delivered: true, Identifier: Identifier{ID: 1, Seq: 8}}))

	// make sure it returns true when applicable
	// TODO
}

func TestMaxSeq(t *testing.T) {
	mod := initModule()
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
	mod := initModule()
	mod.TxObsS = []int{1, 2, 5, 10, 0, 50}

	// TODO figure out a way to mock thetafd.Trusted() so it returns [1,3,5]
}

func TestUpdate(t *testing.T) {
	mod := initModule()

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
}

func TestUrbBroadcast(t *testing.T) {

}

func TestUrbDeliver(t *testing.T) {

}
