package gamedata

import (
	"slices"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/zerrors"
)

var (
	EmptyRadioPrizes = NewEmptyRadioPrizes()
)

func TryNewRadioPrizesList(radioSlice []map[int64]uint64) ([]*RadioPrizes, error) {
	ret := make([]*RadioPrizes, 0, len(radioSlice))

	for _, radios := range radioSlice {
		if len(radios) == 0 {
			continue
		}

		radioPrizes, err := TryNewRadioPrizes(radios)
		if err != nil {
			if errors.Is(err, zerrors.ErrEmptyPrize) {
				radioPrizes = EmptyRadioPrizes
			} else {
				return nil, err
			}
		}

		ret = append(ret, radioPrizes)
	}

	return ret, nil
}

// RadioPrizes is a map of radio id to radio prize. It is built from game data and is immutable.
type RadioPrizes struct {
	radios    map[int64]*RadioPrize
	radioList []*RadioPrize
}

func NewEmptyRadioPrizes() *RadioPrizes {
	return &RadioPrizes{
		radios: make(map[int64]*RadioPrize),
	}
}

func TryNewRadioPrizes(radios map[int64]uint64) (*RadioPrizes, error) {
	if len(radios) == 0 {
		return &RadioPrizes{
			radios:    make(map[int64]*RadioPrize),
			radioList: make([]*RadioPrize, 0),
		}, nil
	}

	radioDatas := make(map[int64]*RadioPrize, len(radios))
	radioList := make([]*RadioPrize, 0, len(radios))

	for radioId, amount := range radios {
		radioPrize, err := TryNewRadioPrize(radioId, amount)
		if err != nil {
			return nil, err
		}

		radioDatas[radioId] = radioPrize
		radioList = append(radioList, radioPrize)
	}

	return &RadioPrizes{
		radios:    radioDatas,
		radioList: radioList,
	}, nil
}

func (p *RadioPrizes) CloneWith(others ...*RadioPrizes) *RadioPrizes {
	if len(others) == 0 {
		return p
	}

	ret := &RadioPrizes{
		radios: make(map[int64]*RadioPrize),
	}

	// clone current packs
	for radioId, radioPrize := range p.radios {
		if radioPrize == nil || radioPrize.radioData == nil || radioPrize.amount == 0 {
			continue
		}

		ret.radios[radioId] = NewRadioPrize(radioPrize.radioData, radioPrize.amount)
	}

	for _, other := range others {
		if other == nil {
			continue
		}

		for radioId, radioPrize := range other.radios {
			if radioPrize == nil || radioPrize.radioData == nil || radioPrize.amount == 0 {
				continue
			}

			if _, ok := ret.radios[radioId]; ok {
				ret.radios[radioId].amount += radioPrize.amount
			} else {
				ret.radios[radioId] = NewRadioPrize(radioPrize.radioData, radioPrize.amount)
			}
		}
	}

	return ret
}

func (rs *RadioPrizes) Rand() *ItemPrizes {
	ret := NewEmptyItemPrizes()

	for _, r := range rs.radioList {
		index := r.radioData.GroupWeightsGroupWeights.Rand()
		if index >= 0 && int(index) < len(r.radioData.GroupItemsItemPrizesList) {
			prize := r.radioData.GroupItemsItemPrizesList[index]
			if IsItemPrizesValid(prize) {
				ret = ret.CloneWith(prize)
			}
		}
	}

	return ret
}

func IsRadioPrizesValid(p *RadioPrizes) bool {
	return p != nil && len(p.radios) > 0
}

func (p *RadioPrizes) Radios() []*RadioPrize {
	return slices.Clone(p.radioList)
}

func (p *RadioPrizes) IsEmpty() bool {
	return len(p.radios) == 0
}

type RadioPrize struct {
	radioData *ResourceRadioData
	amount    uint64
}

func TryNewRadioPrize(radioId int64, amount uint64) (*RadioPrize, error) {
	if amount == 0 {
		return nil, errors.Wrapf(zerrors.ErrEmptyPrize, "radioID=%d, amount=%d", radioId, amount)
	}

	radioData := GetResourceRadioData(radioId)
	if radioData == nil {
		return nil, errors.Wrapf(zerrors.ErrGameDataNotFound, "radioID=%d", radioId)
	}

	return NewRadioPrize(radioData, amount), nil
}

func NewRadioPrize(radioData *ResourceRadioData, amount uint64) *RadioPrize {
	return &RadioPrize{
		radioData: radioData,
		amount:    amount,
	}
}

func (p *RadioPrize) Data() *ResourceRadioData {
	return p.radioData
}

func (p *RadioPrize) Amount() uint64 {
	return p.amount
}
