package gamedata

import "github.com/go-pantheon/roma/pkg/zerrors"

// Prizes is a combination of item prizes and pack prizes. It is used to calculate the total prize. It is mutable.
type Prizes struct {
	items *ItemPrizes
	packs *PackPrizes
}

func TryNewPrizes(items *ItemPrizes, packs *PackPrizes) (*Prizes, error) {
	if items == nil && packs == nil {
		return nil, zerrors.ErrEmptyPrize
	}

	return &Prizes{
		items: items,
		packs: packs,
	}, nil
}

func (p *Prizes) Merge(others ...*Prizes) *Prizes {
	if len(others) == 0 {
		return p
	}

	for _, other := range others {
		if other == nil {
			continue
		}
		p.items = p.items.CloneWith(other.items)
		p.packs = p.packs.CloneWith(other.packs)
	}

	return p
}

func (p *Prizes) Walk(itemFunc func(itemPrize *ItemPrize) (isContinue bool), packFunc func(packPrize *PackPrize) (isContinue bool)) {
	for _, itemPrize := range p.items.items {
		if !itemFunc(itemPrize) {
			break
		}
	}
	for _, packPrize := range p.packs.packs {
		if !packFunc(packPrize) {
			break
		}
	}
}

func (p *Prizes) ItemPrizes() []*ItemPrize {
	return p.items.items
}

func (p *Prizes) PackPrizes() []*PackPrize {
	return p.packs.packList
}
