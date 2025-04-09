package object

import (
	"github.com/go-pantheon/roma/gamedata"
	"github.com/go-pantheon/roma/pkg/util/maths/u64"
)

type ItemInfo struct {
	data   *gamedata.ResourceItemData
	amount uint64
}

func NewItemInfo(itemData *gamedata.ResourceItemData, amount uint64) *ItemInfo {
	return &ItemInfo{
		data:   itemData,
		amount: amount,
	}
}

func (o *ItemInfo) Data() *gamedata.ResourceItemData {
	return o.data
}

func (o *ItemInfo) Amount() uint64 {
	return o.amount
}

func (o *ItemInfo) Add(toAdd uint64) {
	o.amount = u64.Min(u64.Add(o.amount, toAdd), o.data.Max())
}

func (o *ItemInfo) Sub(toSub uint64) {
	o.amount = u64.Sub(o.amount, toSub)
}
