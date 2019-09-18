package models

import (
	"fmt"
)

// Processor represents a server/node in the network
type Processor struct {
	ID       int
	Hostname string
	IPString string
	IP       []byte
}

func (p Processor) String() string {
	return fmt.Sprintf("Processor %d - %s - %s", p.ID, p.Hostname, p.IPString)
}
