package urb

import (
	"testing"

	"gotest.tools/assert"
)

func TestIsSubSet(t *testing.T) {
	s := map[int]bool{}
	s2 := map[int]bool{0: true, 1: true, 2: true}
	s3 := map[int]bool{0: true, 1: true, 2: true, 3: true}

	assert.Assert(t, isSubset(s, s))
	assert.Assert(t, isSubset(s, s2))
	assert.Assert(t, isSubset(s2, s2))
	assert.Assert(t, isSubset(s2, s3))
	assert.Assert(t, !isSubset(s3, s))
}

func TestMax(t *testing.T) {
	assert.Equal(t, max(5, 3), 5)
	assert.Equal(t, max(-5, 2), 2)
	assert.Equal(t, max(0, 0), 0)
	assert.Equal(t, max(0, 5), 5)
	assert.Equal(t, max(-2, 0), 0)
}

func TestContains(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Assert(t, contains(s, 1))
	assert.Assert(t, !contains(s, -1))
	assert.Assert(t, !contains(s, 0))
	assert.Assert(t, !contains(s, 100))
}
