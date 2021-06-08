package common

type RegistryEventOption uint32

const (
	RegistryAddEvent    = RegistryEventOption(0)
	RegistryUpdateEvent = RegistryEventOption(1)
	RegistryDeleteEvent = RegistryEventOption(2)
)

type RegistryChangeEvent struct {
	Opt  RegistryEventOption
	Addr Address
}

func NewRegistryChangeEvent(opt RegistryEventOption, addr Address) *RegistryChangeEvent {
	return &RegistryChangeEvent{
		Opt:  opt,
		Addr: addr,
	}
}
