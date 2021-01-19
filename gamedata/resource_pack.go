package gamedata

import "github.com/vulcan-frame/vulcan-game/pkg/util/maths/u64"

type ResourcePackData struct {
	*ResourcePackDataGen
}

func (d *ResourcePackData) Max() uint64 {
	if d.ResourcePackDataGen.Max == 0 {
		return DefaultMaxAmount
	}
	return u64.Min(d.ResourcePackDataGen.Max, DefaultMaxAmount)
}

func (d *ResourcePackData) Rand() *ItemPrizes {
	ret := d.ItemsItemPrizes

	prize := d.RadiosRadioPrizes.Rand()
	if IsItemPrizesValid(prize) {
		ret = ret.CloneWith(prize)
	}

	index := d.GroupWeightGroupWeights.Rand()
	if index < 0 {
		return ret
	}

	if int(index) < len(d.GroupItemsItemPrizesList) {
		prize := d.GroupItemsItemPrizesList[index]
		if IsItemPrizesValid(prize) {
			ret = ret.CloneWith(prize)
		}
	}

	if int(index) < len(d.GroupRadiosRadioPrizesList) {
		radio := d.GroupRadiosRadioPrizesList[index]
		if IsRadioPrizesValid(radio) {
			prize := radio.Rand()
			if IsItemPrizesValid(prize) {
				ret = ret.CloneWith(prize)
			}
		}
	}
	return ret
}
