package common

import (
	"sync"
)

import (
	"github.com/glory-go/monkey"
)

type DebugMetadata struct {
	ID       string
	GuardMap map[string]*GuardInfo
}

type GuardInfo struct {
	Guard *monkey.PatchGuard
	Lock  sync.Mutex
}
