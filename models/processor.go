package models

import "fmt"

type Processor struct {
	ID        int
	Hostname  string
	IPAddress string
}

func (p Processor) String() string {
	return fmt.Sprintf("Processor %d - %s - %s", p.ID, p.Hostname, p.IPAddress)
}
