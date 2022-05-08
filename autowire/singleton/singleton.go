package singleton

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/base"
	"github.com/glory-go/glory/autowire/param_loader"
	"github.com/glory-go/glory/autowire/sdid_parser"
)

func init() {
	autowire.RegisterAutowire(NewSingletonAutowire(nil, nil, nil))
}

const Name = "singleton"

var singletonStructDescribersMap = make(map[string]*autowire.StructDescriber)

// autowire APIs

// NewSingletonAutowire create a singleton autowire based autowire, e.g. grpc, base.facade can be re-write to outer autowire
func NewSingletonAutowire(sp autowire.SDIDParser, pl autowire.ParamLoader, facade autowire.Autowire) autowire.Autowire {
	if sp == nil {
		sp = sdid_parser.GetDefaultSDIDParser()
	}
	if pl == nil {
		pl = param_loader.GetDefaultParamLoader()
	}
	singletonAutowire := &SingletonAutowire{
		paramLoader: pl,
		sdIDParser:  sp,
	}
	if facade == nil {
		facade = singletonAutowire
	}
	singletonAutowire.AutowireBase = base.New(facade, sp, pl)
	return singletonAutowire

}

type SingletonAutowire struct {
	base.AutowireBase
	paramLoader autowire.ParamLoader
	sdIDParser  autowire.SDIDParser
}

// GetAllStructDescribers should be re-write by facade
func (s *SingletonAutowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return singletonStructDescribersMap
}

// TagKey should be re-writed by facade autowire
func (s *SingletonAutowire) TagKey() string {
	return Name
}

func (s *SingletonAutowire) IsSingleton() bool {
	return true
}

// developer APIs

func RegisterStructDescriber(sd *autowire.StructDescriber) {
	sd.SetAutowireType(Name)
	singletonStructDescribersMap[sd.ID()] = sd
}

func GetImpl(sdId string) (interface{}, error) {
	return autowire.Impl(Name, sdId, nil)
}
