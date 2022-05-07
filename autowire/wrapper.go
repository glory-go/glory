package autowire

import (
	"os"
	"reflect"
	"runtime"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/glory-go/glory/autowire/util"
)

type WrapperAutowire interface {
	Autowire
	Impl(sdID string, param interface{}) (interface{}, error)
}

func getWrappedAutowire(autowire Autowire, allAutowires map[string]WrapperAutowire) WrapperAutowire {
	return &WrapperAutowireImpl{
		Autowire:     autowire,
		allAutowires: allAutowires,
		initedMap:    map[string]bool{},
	}
}

type WrapperAutowireImpl struct {
	Autowire
	initedMap    map[string]bool
	latestPtr    interface{}
	allAutowires map[string]WrapperAutowire
}

// Impl is used to get impled struct
func (w *WrapperAutowireImpl) Impl(sdID string, param interface{}) (interface{}, error) {
	// 1. check singleton
	if inited, ok := w.initedMap[sdID]; w.Autowire.IsSingleton() && inited && ok {
		return w.latestPtr, nil
	}

	// 2. factory
	impledPtr := w.Autowire.Factory(sdID)

	// 3. in ject
	if err := w.inject(impledPtr, sdID); err != nil {
		return nil, err
	}

	// 4. construct field
	if err := w.Autowire.Construct(sdID, impledPtr, param); err != nil {
		return nil, err
	}

	// 5. record
	w.initedMap[sdID] = true
	w.latestPtr = impledPtr
	return impledPtr, nil
}

// inject do tag autowire and monkey inject
func (w *WrapperAutowireImpl) inject(impledPtr interface{}, sdId string) error {
	sd := w.Autowire.GetAllStructDescribers()[sdId]

	// 1. reflect
	valueOf := reflect.ValueOf(impledPtr)
	valueOfElem := valueOf.Elem()
	typeOf := valueOfElem.Type()
	if typeOf.Kind() != reflect.Struct {
		// not struct, no needs to inject tag and monkey, just return
		return nil
	}

	// deal with struct
	// 3. tag inject
	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		t := typeOf.Field(i)
		var subImpledPtr interface{}
		tagKey := ""
		tagValue := ""
		for _, aw := range w.allAutowires {
			if val, ok := t.Tag.Lookup(aw.TagKey()); ok {
				fieldInfo := &FieldInfo{
					FieldName: t.Name,
					FieldType: t.Type.Name(),
					TagKey:    aw.TagKey(),
					TagValue:  val,
				}
				// create param from field info
				subSDID := aw.ParseSDID(fieldInfo)
				param := aw.ParseParam(subSDID, fieldInfo)
				var err error
				subImpledPtr, err = aw.Impl(subSDID, param)
				if err != nil {
					return err
				}
				tagKey = aw.TagKey()
				tagValue = val
				break // only one tag is support
			}
		}
		if tagKey == "" && tagValue == "" {
			continue
		}
		// set field
		subService := valueOfElem.Field(i)
		if !(subService.IsValid() && subService.CanSet()) {
			err := perrors.Errorf("Failed to autowire interface %s's impl %s service. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
				util.GetStructName(sd.Interface), util.GetStructName(impledPtr), t.Type.Name(), tagKey, tagValue)
			return err
		}
		subService.Set(reflect.ValueOf(subImpledPtr))
	}
	// 2. monkey
	if os.Getenv("GOARCH") == "amd64" || runtime.GOARCH == "amd64" {
		// only service, only amd64 mod can inject monkey function
		mf(impledPtr, sd.ID())
	}
	return nil
}
