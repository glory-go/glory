package normal

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/base"
	"github.com/glory-go/glory/autowire/param_loader"
	"github.com/glory-go/glory/autowire/sdid_parser"
)

func init() {
	autowire.RegisterAutowire(NewNormalAutowire(nil, nil, nil))
}

const Name = "normal"

// NewNormalAutowire create a normal autowire based autowire, e.g. config, base.facade can be re-write to outer autowire
func NewNormalAutowire(sp autowire.SDIDParser, pl autowire.ParamLoader, facade autowire.Autowire) autowire.Autowire {
	if sp == nil {
		sp = sdid_parser.GetDefaultSDIDParser()
	}
	if pl == nil {
		pl = param_loader.GetDefaultParamLoader()
	}
	normalAutowire := &NormalAutowire{}
	if facade == nil {
		facade = normalAutowire
	}
	normalAutowire.AutowireBase = base.New(facade, sp, pl)
	return normalAutowire
}

type NormalAutowire struct {
	base.AutowireBase
}

func (n *NormalAutowire) IsSingleton() bool {
	return false
}

// GetAllStructDescribers should be re-write by facade
func (n *NormalAutowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return normalEntryDescriberMap
}

// TagKey should be re-writed by facade autowire
func (n *NormalAutowire) TagKey() string {
	return Name
}

func (s *NormalAutowire) RelyOnTag() bool {
	return false
}

func (s *NormalAutowire) CanBeEntrance() bool {
	return false
}

var normalEntryDescriberMap = make(map[string]*autowire.StructDescriber)

// developer APIs

func RegisterStructDescriber(s *autowire.StructDescriber) {
	s.SetAutowireType(Name)
	normalEntryDescriberMap[s.ID()] = s
}

func GetImpl(sdID string, param interface{}) (interface{}, error) {
	return autowire.Impl(Name, sdID, param)
}
