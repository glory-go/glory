package autowire

import (
	"github.com/glory-go/glory/autowire/util"
)

// Autowire
type Autowire interface {
	TagKey() string
	// IsSingleton means struct can be boot entrance, and only have one impl globally, only created once.
	IsSingleton() bool
	/*
		CanBeEntrance means the autowire sturct's param needs not to parse field tag value as param.
		By default, singleton can be boot entrance, and normal can't, because normal needs to try to parse tag value to find config key.
		But for grpc autowire, as singloton, can't be entrance because it needs to parse grpc type from cong tag.
	*/
	CanBeEntrance() bool
	Factory(sdId string) interface{}

	/*
		ParseSDID parse FieldInfo to struct describerId

		if field type is struct ptr like
		MyStruct *MyStruct `autowire-type:"MyStruct"`
		FieldInfo would be
		FieldInfo.FieldName == "MyStruct"
		FieldInfo.FieldType == "" // ATTENTION!!!
		FieldInfo.TagKey == "autowire-type"
		FieldInfo.TagValue == "MyStruct"
		You should make sure tag value contains ptr type if you use it.

		if field type is interface like
		MyStruct MyInterface ` `autowire-type:"MyStruct"`
		FieldInfo would be
		FieldInfo.FieldName == "MyStruct"
		FieldInfo.FieldType == "MyInterface"
		FieldInfo.TagKey == "autowire-type"
		FieldInfo.TagValue == "MyStruct"
	*/
	ParseSDID(field *FieldInfo) (string, error)
	ParseParam(sdID string, fi *FieldInfo) (interface{}, error)
	Construct(sdID string, impledPtr, param interface{}) (interface{}, error)
	GetAllStructDescribers() map[string]*StructDescriber
	InjectPosition() InjectPosition
}

var wrapperAutowireMap = make(map[string]WrapperAutowire)

func RegisterAutowire(autowire Autowire) {
	wrapperAutowireMap[autowire.TagKey()] = getWrappedAutowire(autowire, wrapperAutowireMap)
}

func GetAllWrapperAutowires() map[string]WrapperAutowire {
	return wrapperAutowireMap
}

// FieldInfo

type FieldInfo struct {
	FieldName string
	FieldType string
	TagKey    string
	TagValue  string
}

// StructDescriber

type StructDescriber struct {
	Interface     interface{}
	Factory       func() interface{} // raw struct
	ParamFactory  func() interface{}
	ParamLoader   ParamLoader
	ConstructFunc func(impl interface{}, param interface{}) (interface{}, error) // injected
	DestroyFunc   func(impl interface{})

	impledStructPtr interface{} // impledStructPtr is only used to get name
	autowireType    string
}

func (ed *StructDescriber) SetAutowireType(autowireType string) {
	ed.autowireType = autowireType
}

func (ed *StructDescriber) AutowireType() string {
	return ed.autowireType
}

func (ed *StructDescriber) ID() string {
	if ed.impledStructPtr == nil {
		ed.parse()
	}
	return util.GetIdByInterfaceAndImplPtr(ed.Interface, ed.impledStructPtr)
}

func (ed *StructDescriber) parse() {
	ed.impledStructPtr = ed.Factory()
}

// ParamLoader is interface to load param
type ParamLoader interface {
	Load(sd *StructDescriber, fi *FieldInfo) (interface{}, error)
}

// SDIDParser is interface to parse struct descriptor id
type SDIDParser interface {
	Parse(fi *FieldInfo) (string, error)
}

type InjectPosition int

const (
	AfterFactoryCalled     InjectPosition = 0
	AfterConstructorCalled InjectPosition = 1
)
