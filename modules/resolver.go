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

// Resolver facilitates inter-module communication
type Resolver struct {
	Modules map[ModuleType]interface{}
}

func (r *Resolver) hb() []int {
	m := r.Modules[HBFD].(HbfdModule)
	return m.hb()
}

func (r *Resolver) trusted() []int {
	m := r.Modules[THETAFD].(ThetafdModule)
	return m.trusted()
}
