package gamedata

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/pkg/errs"
)

// PackPrizes is a map of pack id to pack prize. It is built from game data and is immutable.
type PackPrizes struct {
	packs    map[int64]*PackPrize
	packList []*PackPrize
}

func TryNewPackPrizes(packs map[int64]uint64) (*PackPrizes, error) {
	if len(packs) == 0 {
		return &PackPrizes{
			packs:    make(map[int64]*PackPrize),
			packList: make([]*PackPrize, 0),
		}, nil
	}

	packDatas := make(map[int64]*PackPrize, len(packs))
	packList := make([]*PackPrize, 0, len(packs))
	for packId, amount := range packs {
		packPrize, err := TryNewPackPrize(packId, amount)
		if err != nil {
			return nil, err
		}
		packDatas[packId] = packPrize
		packList = append(packList, packPrize)
	}

	return &PackPrizes{
		packs:    packDatas,
		packList: packList,
	}, nil
}

func (p *PackPrizes) CloneWith(others ...*PackPrizes) *PackPrizes {
	if len(others) == 0 {
		return p
	}

	ret := &PackPrizes{
		packs: make(map[int64]*PackPrize),
	}

	// clone current packs
	for packId, packPrize := range p.packs {
		if packPrize == nil || packPrize.packData == nil || packPrize.amount == 0 {
			continue
		}
		ret.packs[packId] = NewPackPrize(packPrize.packData, packPrize.amount)
	}

	for _, other := range others {
		if other == nil {
			continue
		}
		for packId, packPrize := range other.packs {
			if packPrize == nil || packPrize.packData == nil || packPrize.amount == 0 {
				continue
			}
			if _, ok := ret.packs[packId]; ok {
				ret.packs[packId].amount += packPrize.amount
			} else {
				ret.packs[packId] = NewPackPrize(packPrize.packData, packPrize.amount)
			}
		}
	}
	return ret
}

func IsPackPrizesValid(p *PackPrizes) bool {
	return p != nil && len(p.packList) > 0
}

func (p *PackPrizes) Packs() []*PackPrize {
	return slices.Clone(p.packList)
}

func (p *PackPrizes) IsEmpty() bool {
	return len(p.packs) == 0
}

type PackPrize struct {
	packData *ResourcePackData
	amount   uint64
}

func TryNewPackPrize(packId int64, amount uint64) (*PackPrize, error) {
	if amount == 0 {
		return nil, errors.Wrapf(errs.ErrEmptyPrize, "id=%d", packId)
	}

	packData := GetResourcePackData(packId)
	if packData == nil {
		return nil, errors.Wrapf(errs.ErrGameDataNotFound, "id=%d", packId)
	}
	return NewPackPrize(packData, amount), nil
}

func NewPackPrize(packData *ResourcePackData, amount uint64) *PackPrize {
	return &PackPrize{
		packData: packData,
		amount:   amount,
	}
}

func (p *PackPrize) Data() *ResourcePackData {
	return p.packData
}

func (p *PackPrize) Amount() uint64 {
	return p.amount
}
