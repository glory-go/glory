package base

import (
	"github.com/glory-go/glory/autowire"
)

type FacadeAutowire interface {
	GetAllStructDescribers() map[string]*autowire.StructDescriber
}

// New return new AutowireBase
func New(facadeAutowire FacadeAutowire, sp autowire.SDIDParser, pl autowire.ParamLoader) AutowireBase {
	return AutowireBase{
		facadeAutowire: facadeAutowire,
		sdIDParser:     sp,
		paramLoader:    pl,
	}
}

type AutowireBase struct {
	facadeAutowire FacadeAutowire
	sdIDParser     autowire.SDIDParser
	paramLoader    autowire.ParamLoader
}

func (a *AutowireBase) Factory(sdId string) interface{} {
	sd := a.facadeAutowire.GetAllStructDescribers()[sdId]
	return sd.Factory()
}

func (a *AutowireBase) Construct(sdID string, impledPtr, param interface{}) (interface{}, error) {
	sd := a.facadeAutowire.GetAllStructDescribers()[sdID]
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, param)
	}
	return impledPtr, nil
}

func (a *AutowireBase) ParseSDID(field *autowire.FieldInfo) (string, error) {
	return a.sdIDParser.Parse(field)
}

func (a *AutowireBase) ParseParam(sdId string, fi *autowire.FieldInfo) (interface{}, error) {
	sd := a.facadeAutowire.GetAllStructDescribers()[sdId]
	if sd.ParamFactory == nil {
		// doesn't register param factory, do not load param, return with success
		return nil, nil
	}
	if sd.ParamLoader != nil {
		// try to use sd ParamLoader
		param, err := sd.ParamLoader.Load(sd, fi)
		if err == nil {
			return param, nil
		}
	}
	// use autowire defined paramLoader as fall back
	return a.paramLoader.Load(sd, fi)
}
