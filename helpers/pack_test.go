package helpers

import (
	"reflect"
	"testing"

	"gotest.tools/assert"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

func TestPackAndUnpack(t *testing.T) {
	msg := models.Message{Type: models.MSG, Sender: 0, Data: map[string]interface{}{"foo": "bar"}}
	encoded, err := Pack(&msg)
	assert.Assert(t, encoded != nil)
	assert.Assert(t, err == nil)

	decoded, err := Unpack(encoded)
	assert.Assert(t, decoded != nil)
	assert.Assert(t, err == nil)

	// check that message was correctly decoded
	assert.Equal(t, decoded.Type, models.MSG)
	assert.Equal(t, decoded.Sender, 0)
	assert.Assert(t, reflect.DeepEqual(decoded.Data, map[string]interface{}{"foo": "bar"}))
}
