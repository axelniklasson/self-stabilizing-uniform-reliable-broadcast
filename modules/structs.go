package modules

import (
	"fmt"
)

// UrbModule models the URB algorithm in the paper
type UrbModule struct {
	ID       int
	P        []int
	Resolver *Resolver
	Seq      int
	Buffer   Buffer
	RxObsS   []int
	TxObsS   []int
}

// MessageType indicates the type of message
type MessageType int

const (
	// MSG represents a broadcasted message
	MSG MessageType = 0
	// MSGack represents an acknowledgement of a broadcasted message
	MSGack MessageType = 1
	// GOSSIP represents messages used by processors to update each other
	GOSSIP MessageType = 2
)

// Message models a message that can be broadcasted using the URB algorithm
type Message struct {
	Type       MessageType
	SenderID   int
	Seq        int
	ReceiverID int
	Contents   []byte
}

// Identifier is a pair (ID, Seq) associating a message with the sender and its local sequence number
type Identifier struct {
	ID  int
	Seq int
}

// Buffer holds a number of Records
type Buffer struct {
	Records []BufferRecord
}

// Get is a helper function that can be used to check membership for a BufferRecord in a Buffer. Returns nil if no record exists with id
func (b Buffer) Get(id Identifier) *BufferRecord {
	for _, r := range b.Records {
		if r.Identifier == id {
			return &r
		}
	}
	return nil
}

// Add is a wrapper to make it cleaner to add records to the buffer
func (b *Buffer) Add(br BufferRecord) {
	b.Records = append(b.Records, br)
}

// BufferRecord models a record residing in the local buffer of a processor
type BufferRecord struct {
	// the actual message
	Msg *Message
	// identifier of message, made up of ID (sender id) and Seq (local sequence number at sender)
	Identifier Identifier
	// holds false only when the message still needs to be delivered
	Delivered bool
	// set that includes the identifiers of processors that have acknowledge the message msg
	RecBy map[int]bool
	// value of the HB failure detector
	PrevHB []int
}

func (br BufferRecord) String() string {
	return fmt.Sprintf(
		"BufferRecord - Msg: %v, ID: %d, Seq: %d, Delivered: %t, RecBy: %v, PrevHB: %v",
		br.Msg, br.Identifier.ID, br.Identifier.Seq, br.Delivered, br.RecBy, br.PrevHB)
}
