package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func testGetNTPTime(t *testing.T) {
	time := GetNTPTime()
	time2 := GetNTPTime()
	assert.Assert(t, time2.After(time))
}
