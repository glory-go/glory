package grpc

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
)

func init() {
	autowire.RegisterAutowire(&Autowire{})
}

const Name = "grpc"

type Autowire struct {
}

func (a *Autowire) TagKey() string {
	return Name
}

func (a *Autowire) Factory(sdId string) interface{} {
	return a.GetAllStructDescribers()[sdId].Factory()
}

func (a *Autowire) Construct(sdID string, impledPtr, _ interface{}) error {
	sd := a.GetAllStructDescribers()[sdID]
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, nil)
	}
	return nil
}

func (a *Autowire) ParseParam(_ string, _ *autowire.FieldInfo) interface{} {
	return nil
}

func (a *Autowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return grpcStructDescriberMap
}

func (a *Autowire) ParseSDID(field *autowire.FieldInfo) string {
	return util.GetIdByNamePair(field.TagValue, field.TagValue)
}

func (a *Autowire) IsSingleton() bool {
	return true
}

var grpcStructDescriberMap = make(map[string]*autowire.StructDescriber)

func RegisterStructDescriber(s *autowire.StructDescriber) {
	grpcStructDescriberMap[s.ID()] = s
}

func GetImpl(extensionId string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, nil)
}
