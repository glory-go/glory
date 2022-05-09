package common

import (
	"sync"
)

import (
	"github.com/glory-go/monkey"
)

type DebugMetadata struct {
	GuardMap map[string]*GuardInfo
}

type GuardInfo struct {
	Guard *monkey.PatchGuard
	Lock  sync.Mutex
}
