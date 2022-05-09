package autowire

import (
	"fmt"
)

import (
	"github.com/fatih/color"
)

var allWrapperAutowires = make(map[string]WrapperAutowire)

func printAutowireRegisteredStructDescriber() {
	for autowireType, aw := range allWrapperAutowires {
		color.Blue("[Autowire Type] Found registered autowire type %s", autowireType)
		for sdID := range aw.GetAllStructDescribers() {
			color.Blue("[Autowire Struct Descriptor] Found type %s registered SD %s", autowireType, sdID)
		}
	}
}

func Load() error {
	// get all autowires
	allWrapperAutowires = GetAllWrapperAutowires()

	printAutowireRegisteredStructDescriber()

	// autowire all struct that can be entrance
	for _, aw := range allWrapperAutowires {
		for sdID := range aw.GetAllStructDescribers() {
			if aw.CanBeEntrance() {
				_, err := aw.ImplWithoutParam(sdID)
				if err != nil {
					panic(fmt.Errorf("[Autowire] Impl sd %s failed, reason is %s", sdID, err))
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
