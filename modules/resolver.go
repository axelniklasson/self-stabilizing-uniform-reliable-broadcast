package modules

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
}

// Resolver facilitates inter-module communication
type Resolver struct {
	Modules map[ModuleType]interface{}
}

func (r *Resolver) hb() []int {
	m := r.Modules[HBFD].(HbfdModule)
	return m.HB()
}

func (r *Resolver) trusted() []int {
	m := r.Modules[THETAFD].(ThetafdModule)
	return m.Trusted()
}
