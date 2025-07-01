package object

import (
	"errors"
	"time"

	"github.com/go-pantheon/roma/gamedata"
)

type UpdateType bool

const (
	UpdateTypeAdd UpdateType = true
	UpdateTypeSub UpdateType = false
)

type UpdateInfo struct {
	tp        UpdateType
	items     map[int64]*updateItemInfo
	packs     map[int64]*updatePackInfo
	updatedAt time.Time
}

func NewUpdateInfo(ctime time.Time, tp UpdateType) *UpdateInfo {
	return &UpdateInfo{
		items:     make(map[int64]*updateItemInfo),
		packs:     make(map[int64]*updatePackInfo),
		updatedAt: ctime,
		tp:        tp,
	}
}

func (u *UpdateInfo) WalkItem(f func(d *gamedata.ResourceItemData, amount uint64)) {
	for _, item := range u.items {
		f(item.itemData, item.changeAmount)
	}
}

func (u *UpdateInfo) WalkPack(f func(d *gamedata.ResourcePackData, amount uint64)) {
	for _, pack := range u.packs {
		f(pack.packData, pack.changeAmount)
	}
}

func (u *UpdateInfo) Merge(other *UpdateInfo) error {
	if u.tp != other.tp {
		return errors.New("type mismatch")
	}

	for _, item := range other.items {
		u.AddItem(item.itemData, item.changeAmount)
	}

	for _, pack := range other.packs {
		u.AddPack(pack.packData, pack.changeAmount)
	}

	return nil
}

func (u *UpdateInfo) AddItem(itemData *gamedata.ResourceItemData, changeAmount uint64) {
	if u.items[itemData.Id()] == nil {
		u.items[itemData.Id()] = newUpdateItemInfo(itemData, changeAmount)
	} else {
		u.items[itemData.Id()].changeAmount += changeAmount
	}
}

func (u *UpdateInfo) AddPack(packData *gamedata.ResourcePackData, changeAmount uint64) {
	if u.packs[packData.Id()] == nil {
		u.packs[packData.Id()] = newUpdatePackInfo(packData, changeAmount)
	} else {
		u.packs[packData.Id()].changeAmount += changeAmount
	}
}

func (u *UpdateInfo) Type() UpdateType {
	return u.tp
}

type updateItemInfo struct {
	itemData     *gamedata.ResourceItemData
	changeAmount uint64
}

func newUpdateItemInfo(itemData *gamedata.ResourceItemData, changeAmount uint64) *updateItemInfo {
	return &updateItemInfo{
		itemData:     itemData,
		changeAmount: changeAmount,
	}
}

type updatePackInfo struct {
	packData     *gamedata.ResourcePackData
	changeAmount uint64
}

func newUpdatePackInfo(packData *gamedata.ResourcePackData, changeAmount uint64) *updatePackInfo {
	return &updatePackInfo{
		packData:     packData,
		changeAmount: changeAmount,
	}
}
