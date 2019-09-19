package helpers

import (
	"encoding/json"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// Pack encodes the message as bytes
func Pack(m *models.Message) ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Unpack decodes a slice of bytes to a Message
func Unpack(bytes []byte) (*models.Message, error) {
	var m models.Message
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
