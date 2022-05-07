package config

import (
	"log"
	"strings"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
	"github.com/glory-go/glory/config"
)

func init() {
	autowire.RegisterAutowire(&Autowire{})
}

const Name = "config"

type Autowire struct {
}

func (a *Autowire) Factory(sdId string) interface{} {
	sd := a.GetAllStructDescribers()[sdId]
	defer func() { // assure the impl procedure of one service run once
		if r := recover(); r != nil {
			log.Printf("recover panic = %s\n", r)
		}
	}()
	return sd.Factory()
}

func (a *Autowire) Construct(sdId string, impledPtr, param interface{}) error {
	sd := a.GetAllStructDescribers()[sdId]
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, param)
	}
	return nil
}

func (a *Autowire) ParseParam(sdID string, fi *autowire.FieldInfo) interface{} {
	splitedTagValue := strings.Split(fi.TagValue, ",")
	configPath := splitedTagValue[1]
	sd := a.GetAllStructDescribers()[sdID]
	if sd.ParamFactory == nil {
		return nil
	}
	param := sd.ParamFactory()
	if err := config.LoadConfigPathValue(configPath, param); err != nil {
		log.Println("load config path "+configPath+" error = ", err)
	}
	return param
}

func (a *Autowire) ParseSDID(fi *autowire.FieldInfo) string {
	splitedTagValue := strings.Split(fi.TagValue, ",")
	return util.GetIdByNamePair(splitedTagValue[0], splitedTagValue[0])
}

func (a *Autowire) IsSingleton() bool {
	return false
}

func (a *Autowire) TagKey() string {
	return Name
}

func (a *Autowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return configStructDescriberMap
}

var configStructDescriberMap = make(map[string]*autowire.StructDescriber)

func RegisterStructDescriber(s *autowire.StructDescriber) {
	configStructDescriberMap[s.ID()] = s
}

func GetImpl(extensionId string, configPrefix string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, configPrefix)
}
