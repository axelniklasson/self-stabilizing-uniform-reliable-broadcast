package helpers

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetNTPTime(t *testing.T) {
	time := GetNTPTime()
	time2 := GetNTPTime()
	assert.Assert(t, time2.After(time))
}
