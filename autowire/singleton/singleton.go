package singleton

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
)

func init() {
	autowire.RegisterAutowire(&SingletonAutowire{})
}

const Name = "singleton"

var singletonStructDescribersMap = make(map[string]*autowire.StructDescriber)

// autowire APIs

type SingletonAutowire struct {
}

func (s *SingletonAutowire) Factory(sdId string) interface{} {
	// fixme: here returns one, but wrapper layer still autowire again? and call constructor again?
	sd := s.GetAllStructDescribers()[sdId]
	impledPtr := sd.Factory()
	return impledPtr
}

func (s *SingletonAutowire) ParseParam(_ string, _ *autowire.FieldInfo) interface{} {
	return nil
}

func (s *SingletonAutowire) Construct(sdID string, impledPtr, _ interface{}) error {
	sd := s.GetAllStructDescribers()[sdID]
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, nil)
	}
	return nil
}

func (s *SingletonAutowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return singletonStructDescribersMap
}

func (s *SingletonAutowire) TagKey() string {
	return Name
}

func (s *SingletonAutowire) ParseSDID(field *autowire.FieldInfo) string {
	return util.GetIdByNamePair(field.FieldType, field.TagValue)
}

func (s *SingletonAutowire) IsSingleton() bool {
	return true
}

// developer APIs

func RegisterStructDescriber(sd *autowire.StructDescriber) {
	singletonStructDescribersMap[sd.ID()] = sd
}

func GetImpl(sdId string) (interface{}, error) {
	return autowire.Impl(Name, sdId, nil)
}
