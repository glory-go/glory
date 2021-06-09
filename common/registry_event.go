package common

type RegistryEventOption uint32

const (
	RegistryAddEvent                  = RegistryEventOption(0)
	RegistryUpdateEvent               = RegistryEventOption(1)
	RegistryDeleteEvent               = RegistryEventOption(2)
	RegistryUpdateToSerivcesListEvent = RegistryEventOption(3)
)

type RegistryChangeEvent struct {
	Opt      RegistryEventOption
	Addr     Address
	Serivces []Address
}

func NewRegistryChangeEvent(opt RegistryEventOption, addr Address) *RegistryChangeEvent {
	return &RegistryChangeEvent{
		Opt:  opt,
		Addr: addr,
	}
}

func NewReigstryUpdateToServiceEvent(addrs []Address) *RegistryChangeEvent {
	return &RegistryChangeEvent{
		Opt:      RegistryUpdateToSerivcesListEvent,
		Serivces: addrs,
	}
}
