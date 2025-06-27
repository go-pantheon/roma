package gamedata

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/zerrors"
)

// Costs is a map of item id to item cost. It is built from game data and is immutable.
type Costs struct {
	items []*ItemCost
}

func TryNewCosts(itemAmounts map[int64]uint64) (*Costs, error) {
	if len(itemAmounts) == 0 {
		return &Costs{
			items: make([]*ItemCost, 0),
		}, nil
	}

	items := make([]*ItemCost, 0, len(itemAmounts))
	for itemId, amount := range itemAmounts {
		itemCost, err := TryNewItemCost(itemId, amount)
		if err != nil {
			return nil, err
		}
		items = append(items, itemCost)
	}

	return NewCosts(items...), nil
}

func NewCosts(items ...*ItemCost) *Costs {
	return &Costs{
		items: items,
	}
}

func (c *Costs) CloneWith(others ...*Costs) *Costs {
	if len(others) == 0 {
		return c
	}

	items := make(map[int64]*ItemCost, len(c.items)+len(others))

	for _, itemCost := range c.items {
		if itemCost == nil || itemCost.itemData == nil || itemCost.amount == 0 {
			continue
		}
		items[itemCost.itemData.ID] = NewItemCost(itemCost.itemData, itemCost.amount)
	}

	for _, other := range others {
		if other == nil {
			continue
		}

		for _, itemCost := range other.items {
			if itemCost == nil || itemCost.itemData == nil || itemCost.amount == 0 {
				continue
			}

			if _, ok := items[itemCost.itemData.ID]; ok {
				items[itemCost.itemData.ID].amount += itemCost.amount
			} else {
				items[itemCost.itemData.ID] = NewItemCost(itemCost.itemData, itemCost.amount)
			}
		}
	}

	ret := &Costs{
		items: make([]*ItemCost, 0, len(items)),
	}

	for _, itemCost := range items {
		ret.items = append(ret.items, itemCost)
	}

	return ret
}

func (c *Costs) Walk(f func(itemCost *ItemCost) (isContinue bool)) {
	for _, itemCost := range c.items {
		if !f(itemCost) {
			break
		}
	}
}

func (c *Costs) ItemList() []*ItemCost {
	return c.items
}

func (c *Costs) IsEmpty() bool {
	return len(c.items) == 0
}

// ItemCost costs for items
type ItemCost struct {
	itemData *ResourceItemData
	amount   uint64
}

func TryNewItemCost(itemId int64, amount uint64) (*ItemCost, error) {
	if amount == 0 {
		return nil, errors.Wrapf(zerrors.ErrEmptyCost, "itemId=%d", itemId)
	}

	itemData := GetResourceItemData(itemId)
	if itemData == nil {
		return nil, errors.Wrapf(zerrors.ErrGameDataNotFound, "id=%d", itemId)
	}
	
	return NewItemCost(itemData, amount), nil
}

func NewItemCost(itemData *ResourceItemData, amount uint64) *ItemCost {
	return &ItemCost{
		itemData: itemData,
		amount:   amount,
	}
}

func (c *ItemCost) Data() *ResourceItemData {
	return c.itemData
}

func (c *ItemCost) Amount() uint64 {
	return c.amount
}
