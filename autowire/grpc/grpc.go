package grpc

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/singleton"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		grpcAutowire := &Autowire{}
		grpcAutowire.Autowire = singleton.NewSingletonAutowire(&sdIDParser{}, &paramLoader{}, grpcAutowire)
		return grpcAutowire
	}())
}

const Name = "grpc"

type Autowire struct {
	autowire.Autowire
}

// TagKey re-write SingletonAutowire
func (a *Autowire) TagKey() string {
	return Name
}

// GetAllStructDescribers re-write SingletonAutowire
func (a *Autowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return grpcStructDescriberMap
}

var grpcStructDescriberMap = make(map[string]*autowire.StructDescriber)

func RegisterStructDescriber(s *autowire.StructDescriber) {
	s.SetAutowireType(Name)
	grpcStructDescriberMap[s.ID()] = s
}

func GetImpl(extensionId string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, nil)
}
