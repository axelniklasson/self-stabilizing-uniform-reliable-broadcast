package modules

// Message models a message that can be broadcasted using the URB algorithm
type Message struct {
	Type       MessageType
	SenderID   int
	Seq        int
	ReceiverID int
	Contents   []byte
}
