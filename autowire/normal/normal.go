package normal

import (
	"encoding/json"
	"log"
	"strings"
)

import (
	"github.com/glory-go/glory/autowire"
	"github.com/glory-go/glory/autowire/util"
)

func init() {
	autowire.RegisterAutowire(&NormalAutowire{})
}

const Name = "normal"

type NormalAutowire struct {
}

func (n *NormalAutowire) Factory(sdId string) interface{} {
	return n.GetAllStructDescribers()[sdId].Factory()
}

func (n *NormalAutowire) Construct(sdID string, impledPtr, param interface{}) error {
	sd := n.GetAllStructDescribers()[sdID]
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, param)
	}
	return nil
}

func (n *NormalAutowire) ParseParam(sdID string, fi *autowire.FieldInfo) interface{} {
	sd := normalEntryDescriberMap[sdID]
	if sd.ParamFactory == nil {
		return nil
	}
	splitedTagValue := strings.Split(fi.TagValue, ",")
	kvs := strings.Split(splitedTagValue[1], "&")
	kvMaps := make(map[string]interface{})
	for _, kv := range kvs {
		splitedKV := strings.Split(kv, "=")
		kvMaps[splitedKV[0]] = splitedKV[1]
	}
	data, err := json.Marshal(kvMaps)
	if err != nil {
		log.Printf("error json marshal %s\n", err)
		return nil
	}
	param := sd.ParamFactory()
	if err := json.Unmarshal(data, param); err != nil {
		log.Printf("error jsonun marshal %s\n", err)
		return nil
	}
	return param
}

func (n *NormalAutowire) ParseSDID(field *autowire.FieldInfo) string {
	splitedTagValue := strings.Split(field.TagValue, ",")
	return util.GetIdByNamePair(field.FieldType, splitedTagValue[0])
}

func (n *NormalAutowire) IsSingleton() bool {
	return false
}

func (n *NormalAutowire) TagKey() string {
	return Name
}

func (n *NormalAutowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return normalEntryDescriberMap
}

var normalEntryDescriberMap = make(map[string]*autowire.StructDescriber)

// developer APIs

func RegisterStructDescriber(s *autowire.StructDescriber) {
	normalEntryDescriberMap[s.ID()] = s
}

func GetImpl(extensionId string, param interface{}) (interface{}, error) {
	return autowire.Impl(Name, extensionId, param)
}
