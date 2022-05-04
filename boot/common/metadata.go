package common

import (
	"github.com/glory-go/monkey"
)

type RegisterServiceMetadata struct {
	InterfaceStruct    interface{}
	SvcStructPtr       interface{}
	ConstructFunctions []func(interface{})
	GuardMap           map[string]*monkey.PatchGuard
	IsController       bool
}
