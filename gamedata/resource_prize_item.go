package gamedata

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/pkg/errs"
)

var (
	EmptyItemPrizes = NewEmptyItemPrizes()
)

func TryNewItemPrizesList(itemSlice []map[int64]uint64) ([]*ItemPrizes, error) {
	ret := make([]*ItemPrizes, 0, len(itemSlice))
	for _, items := range itemSlice {
		if len(items) == 0 {
			continue
		}
		itemPrizes, err := TryNewItemPrizes(items)
		if err != nil {
			if errors.Is(err, errs.ErrEmptyPrize) {
				itemPrizes = EmptyItemPrizes
			} else {
				return nil, err
			}
		}
		ret = append(ret, itemPrizes)
	}
	return ret, nil
}

// ItemPrizes is a map of item id to item prize. It is built from game data and is immutable.
type ItemPrizes struct {
	items []*ItemPrize
}

func TryNewItemPrizes(itemAmounts map[int64]uint64) (*ItemPrizes, error) {
	if len(itemAmounts) == 0 {
		return EmptyItemPrizes, nil
	}

	items := make([]*ItemPrize, 0, len(itemAmounts))
	for itemId, amount := range itemAmounts {
		itemPrize, err := TryNewItemPrize(itemId, amount)
		if err != nil {
			return nil, err
		}
		items = append(items, itemPrize)
	}

	return NewItemPrizes(items...), nil
}

func NewEmptyItemPrizes() *ItemPrizes {
	return &ItemPrizes{
		items: make([]*ItemPrize, 0),
	}
}

func NewItemPrizes(items ...*ItemPrize) *ItemPrizes {
	return &ItemPrizes{
		items: items,
	}
}

func (p *ItemPrizes) CloneWith(others ...*ItemPrizes) *ItemPrizes {
	if len(others) == 0 {
		return p
	}

	items := make(map[int64]*ItemPrize, len(p.items)*len(others))

	for _, itemPrize := range p.items {
		if itemPrize == nil || itemPrize.itemData == nil || itemPrize.amount == 0 {
			continue
		}
		items[itemPrize.itemData.ID] = NewItemPrize(itemPrize.itemData, itemPrize.amount)
	}

	for _, other := range others {
		if other == nil {
			continue
		}
		for _, itemPrize := range other.items {
			if itemPrize == nil || itemPrize.itemData == nil || itemPrize.amount == 0 {
				continue
			}
			if _, ok := items[itemPrize.itemData.ID]; ok {
				items[itemPrize.itemData.ID].amount += itemPrize.amount
			} else {
				items[itemPrize.itemData.ID] = NewItemPrize(itemPrize.itemData, itemPrize.amount)
			}
		}
	}

	ret := &ItemPrizes{
		items: make([]*ItemPrize, 0, len(items)),
	}
	for _, itemPrize := range items {
		ret.items = append(ret.items, itemPrize)
	}
	return ret
}

func IsItemPrizesValid(p *ItemPrizes) bool {
	return p != nil && len(p.items) > 0
}

func (p *ItemPrizes) Items() []*ItemPrize {
	return slices.Clone(p.items)
}

func (p *ItemPrizes) IsEmpty() bool {
	return len(p.items) == 0
}

type ItemPrize struct {
	itemData *ResourceItemData
	amount   uint64
}

func TryNewItemPrize(itemId int64, amount uint64) (*ItemPrize, error) {
	if amount == 0 {
		return nil, errors.Wrapf(errs.ErrEmptyPrize, "id=%d", itemId)
	}

	itemData := GetResourceItemData(itemId)
	if itemData == nil {
		return nil, errors.Wrapf(errs.ErrGameDataNotFound, "id=%d", itemId)
	}
	return NewItemPrize(itemData, amount), nil
}

func NewItemPrize(itemData *ResourceItemData, amount uint64) *ItemPrize {
	return &ItemPrize{
		itemData: itemData,
		amount:   amount,
	}
}

func (o *ItemPrize) Data() *ResourceItemData {
	return o.itemData
}

func (o *ItemPrize) Amount() uint64 {
	return o.amount
}
