package autowire

import (
	"fmt"
)

var allWrapperAutowires = make(map[string]WrapperAutowire)

func Load() error {
	// get all autowires
	allWrapperAutowires = GetAllWrapperAutowires()

	// autowire all struct that can be entrance
	for _, aw := range allWrapperAutowires {
		for sdID := range aw.GetAllStructDescribers() {
			if aw.CanBeEntrance() {
				_, err := aw.ImplWithoutParam(sdID)
				if err != nil {
					panic(fmt.Errorf("[Boot ImplWithoutParam] sd %s failed, reason is %s", sdID, err))
				}
			}
		}
	}
	return nil
}

func Impl(autowireType, structDescriberID string, param interface{}) (interface{}, error) {
	for _, wrapperAutowire := range allWrapperAutowires {
		if wrapperAutowire.TagKey() == autowireType {
			return wrapperAutowire.ImplWithParam(structDescriberID, param)
		}
	}
	return nil, nil
}
