package autowire

import (
	"github.com/glory-go/glory/autowire/util"
)

// Autowire
type Autowire interface {
	TagKey() string
	Factory(sdId string) interface{}
	Construct(sdID string, impledPtr, param interface{}) error
	ParseParam(sdID string, fi *FieldInfo) interface{}
	GetAllStructDescribers() map[string]*StructDescriber

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
	ParseSDID(field *FieldInfo) string

	// IsSingleton means struct can't be boot entrance, and only have one impl globally, only created once.
	IsSingleton() bool
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
	ConstructFunc func(impl interface{}, param interface{}) error // injected
	DestroyFunc   func(impl interface{})

	impledStructPtr interface{} // impledStructPtr is only used to get name
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
