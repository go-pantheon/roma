package core

import (
	"github.com/go-pantheon/roma/pkg/universe/life"
)

const (
	WorkerEventTypeStorageItemUpdated = life.WorkerEventType(iota + 1) // storage updated. args: item IDs
	WorkerEventTypeStoragePackUpdated                                  // storage updated. args: pack IDs
	WorkerEventTypeHeroUpdated                                         // hero updated. args: none
	WorkerEventTypeEffectParamUpdated                                  // effect param updated. args: none
)

type argKeyStorageItemIDs struct{}
type argKeyStoragePackIDs struct{}
type argKeyHeroIDs struct{}
type argKeyEffectParamIDs struct{}

var (
	ArgKeyStorageItemIDs = argKeyStorageItemIDs{}
	ArgKeyStoragePackIDs = argKeyStoragePackIDs{}
	ArgKeyHeroIDs        = argKeyHeroIDs{}
	ArgKeyEffectParamIDs = argKeyEffectParamIDs{}
)
