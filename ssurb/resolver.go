package ssurb

import (
	"log"

	"github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast/models"
)

// ModuleType indicates the type of module
type ModuleType int

const (
	// URB refers to UrbModule
	URB ModuleType = 0
	// HBFD refers to HbfdModule
	HBFD ModuleType = 1
	// THETAFD refers to ThetafdModule
	THETAFD ModuleType = 2
)

// IResolver defines what interface functions are available for inter-module communication
// Also makes it testable..
type IResolver interface {
	hb() []int
	trusted() []int

	Dispatch(*models.Message)
}

// Resolver facilitates inter-module communication
type Resolver struct {
	Modules map[ModuleType]interface{}
}

// Dispatch routes an incoming message to the correct module
func (r *Resolver) Dispatch(m *models.Message) {
	urbModule := r.Modules[URB].(UrbModule)

	switch m.Type {
	case models.MSG:
		urbModule.onMSG(m)
	case models.MSGack:
		urbModule.onMSGack(m)
	case models.GOSSIP:
		urbModule.onGOSSIP(m)
	default:
		log.Fatalf("Got unrecognized message %v", m)
	}
}

func (r *Resolver) hb() []int {
	m := r.Modules[HBFD].(HbfdModule)
	return m.HB()
}

func (r *Resolver) trusted() []int {
	m := r.Modules[THETAFD].(ThetafdModule)
	return m.Trusted()
}
