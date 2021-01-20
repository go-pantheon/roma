package core

import (
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
)

const (
	WorkerEventTypeStorageItemUpdated = life.WorkerEventType(iota + 1) // storage updated. args: item IDs
	WorkerEventTypeStoragePackUpdated = life.WorkerEventType(iota + 1) // storage updated. args: pack IDs
	WorkerEventTypeEffectParamUpdated                                  // effect param updated. args: none
)
