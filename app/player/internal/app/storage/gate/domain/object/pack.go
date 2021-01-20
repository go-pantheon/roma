package object

import (
	"github.com/vulcan-frame/vulcan-game/gamedata"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/u64"
)

type PackInfo struct {
	data   *gamedata.ResourcePackData
	amount uint64
}

func NewPackInfo(packData *gamedata.ResourcePackData, amount uint64) *PackInfo {
	return &PackInfo{
		data:   packData,
		amount: amount,
	}
}

func (o *PackInfo) Data() *gamedata.ResourcePackData {
	return o.data
}

func (o *PackInfo) Amount() uint64 {
	return o.amount
}

func (o *PackInfo) Add(toAdd uint64) {
	o.amount = u64.Min(u64.Add(o.amount, toAdd), o.data.Max())
}

func (o *PackInfo) Sub(toSub uint64) {
	o.amount = u64.Sub(o.amount, toSub)
}
