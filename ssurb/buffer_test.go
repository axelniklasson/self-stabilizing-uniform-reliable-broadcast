package ssurb

import (
	"testing"

	"gotest.tools/assert"
)

func TestGet(t *testing.T) {
	// init buffer and add 3 records
	buf := Buffer{Records: []*BufferRecord{}}
	r := BufferRecord{Identifier: Identifier{ID: 0, Seq: 0}}
	r2 := BufferRecord{Identifier: Identifier{ID: 0, Seq: 1}}
	r3 := BufferRecord{Identifier: Identifier{ID: 0, Seq: 2}}
	buf.Add(&r)
	buf.Add(&r2)
	buf.Add(&r3)

	// make sure that Get returns correct record for a given identifier and nil if no record with id exists in buffer
	assert.Equal(t, len(buf.Records), 3)
	assert.Equal(t, buf.Get(Identifier{ID: 0, Seq: 0}).Identifier, Identifier{ID: 0, Seq: 0})
	assert.Assert(t, buf.Get(Identifier{ID: 1, Seq: 0}) == nil)
}

func TestAdd(t *testing.T) {
	// init buffer and check that a record can be added
	buf := Buffer{Records: []*BufferRecord{}}
	assert.Equal(t, len(buf.Records), 0)
	buf.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 0}})
	assert.Equal(t, len(buf.Records), 1)

	// adding another record should (not surprisingly) add that record to the buffer
	buf.Add(&BufferRecord{Identifier: Identifier{ID: 0, Seq: 1}})
	assert.Equal(t, len(buf.Records), 2)
	assert.Equal(t, buf.Records[1].Identifier, Identifier{ID: 0, Seq: 1})
}
