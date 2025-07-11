package gamedata

import (
	"math"

	"github.com/go-pantheon/roma/pkg/util/maths/u64"
)

// DefaultMaxAmount is used when the item has no maximum amount.
// It is the maximum value of int64 for compatibility with other programming languages.
const DefaultMaxAmount = math.MaxInt64

type ItemType int64

const (
	ItemTypeUnspecified      ItemType = iota // 0
	ItemTypeRechargeCurrency                 // 1
	ItemTypeGameCurrency                     // 2
	ItemTypeRecovery                         // 3
	ItemTypeHeroPiece                        // 4
)

type ResourceItemData struct {
	*ResourceItemDataGen

	Type ItemType
}

func (d *ResourceItemData) build() {
	d.ResourceItemDataGen.build()

	d.Type = ItemType(d.ItemTypeInt)
}

func (d *ResourceItemData) Max() uint64 {
	if d.ResourceItemDataGen.Max == 0 {
		return DefaultMaxAmount
	}

	return u64.Min(d.ResourceItemDataGen.Max, DefaultMaxAmount)
}
