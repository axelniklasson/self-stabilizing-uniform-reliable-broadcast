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
	Hb() []int
	Trusted() []int
	UrbBroadcast(*UrbMessage)
	Dispatch(*models.Message)
}

// Resolver facilitates inter-module communication
type Resolver struct {
	Modules map[ModuleType]interface{}
}

// Hb calles the HB funciton in the hbfd module
func (r *Resolver) Hb() []int {
	m := r.Modules[HBFD].(*HbfdModule)
	return m.HB()
}

// Trusted calls the Trusted function in the theta fd module
func (r *Resolver) Trusted() []int {
	m := r.Modules[THETAFD].(*ThetafdModule)
	return m.Trusted()
}

// UrbBroadcast is called by the API whenever a message came from the application layer to be broadcasted
func (r *Resolver) UrbBroadcast(msg *UrbMessage) {
	m := r.Modules[URB].(*UrbModule)
	m.UrbBroadcast(msg)
}

// Dispatch routes an incoming message to the correct module
func (r *Resolver) Dispatch(m *models.Message) {
	urbModule := r.Modules[URB].(*UrbModule)
	hbfdModule := r.Modules[HBFD].(*HbfdModule)
	thetafdModule := r.Modules[THETAFD].(*ThetafdModule)

	switch m.Type {
	case models.MSG:
		urbModule.onMSG(m)
	case models.MSGack:
		urbModule.onMSGack(m)
	case models.GOSSIP:
		urbModule.onGOSSIP(m)
	case models.HBFDheartbeat:
		hbfdModule.onHeartbeat(m.Sender)
	case models.THETAheartbeat:
		thetafdModule.onHeartbeat(m.Sender)
	default:
		log.Fatalf("Got unrecognized message %v", m)
	}
}

// GetUrbModule is used to get the current isntance of the urb module
func (r *Resolver) GetUrbModule() *UrbModule {
	urbModule := r.Modules[URB].(*UrbModule)
	return urbModule
}

// GetUrbData is called by the API to return data to clients
func (r *Resolver) GetUrbData() map[string]interface{} {
	urbModule := r.Modules[URB].(*UrbModule)
	return urbModule.GetData()
}
