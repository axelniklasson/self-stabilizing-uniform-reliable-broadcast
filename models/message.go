package models

// MessageType indicates the type of message
type MessageType int

const (
	// MSG represents a broadcasted message
	MSG MessageType = 0
	// MSGack represents an acknowledgement of a broadcasted message
	MSGack MessageType = 1
	// GOSSIP represents messages used by processors to update each other
	GOSSIP MessageType = 2
	// HBFDheartbeat represents a hbfd message
	HBFDheartbeat MessageType = 3
	// THETAheartbeat represents a hbfd message
	THETAheartbeat MessageType = 4
)

// Message represents a message sent between two processors over UDP
type Message struct {
	Type   MessageType
	Sender int
	Data   map[string]interface{}
}
