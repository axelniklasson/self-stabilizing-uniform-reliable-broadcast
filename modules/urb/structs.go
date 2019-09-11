package urb

type UrbModule struct {
	ID             int
	P              []int
	Seq            int
	Buffer         Buffer
	SeqMin         []int
	BufferUnitSize int
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

type Message struct {
	Type       MessageType
	SenderID   int
	Seq        int
	ReceiverID int
	Contents   []byte
}

type Identifier struct {
	ID  int
	Seq int
}

type Buffer struct {
	Records []BufferRecord
}

func (b Buffer) Contains(id Identifier) *BufferRecord {
	for _, r := range b.Records {
		if r.Identifier == id {
			return &r
		}
	}
	return nil
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
	// TODO
	PrevHB []int
}
