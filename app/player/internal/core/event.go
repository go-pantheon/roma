package core

import (
	"github.com/go-pantheon/roma/pkg/universe/life"
)

const (
	WorkerEventTypeStorageItemUpdated = life.WorkerEventType(iota + 1) // storage updated. args: item IDs
	WorkerEventTypeStoragePackUpdated = life.WorkerEventType(iota + 1) // storage updated. args: pack IDs
	WorkerEventTypeHeroUpdated                                         // hero updated. args: none
	WorkerEventTypeEffectParamUpdated                                  // effect param updated. args: none
)
