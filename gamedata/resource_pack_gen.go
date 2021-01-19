// Code generated by gen-datas. DO NOT EDIT.

package gamedata

import (
	resource_base "github.com/vulcan-frame/vulcan-game/gen/gamedata/base/resource"
	"github.com/pkg/errors"
)

var _ = errors.New("import holding")

// ResourcePackDatas excel/Resource/item.xlsx:Pack
type ResourcePackDatas struct {
	List []*ResourcePackData
	Map  map[int64]*ResourcePackData
}

// ResourcePackDataGen excel/Resource/item.xlsx:Pack
type ResourcePackDataGen struct {
	*resource_base.PackDataBaseGen

	ItemsItemPrizes            *ItemPrizes
	RadiosRadioPrizes          *RadioPrizes
	GroupWeightGroupWeights    *GroupWeights
	GroupItemsItemPrizesList   []*ItemPrizes
	GroupRadiosRadioPrizesList []*RadioPrizes
}

func newResourcePackDatas(bases *resource_base.PackDataBaseGens) (*ResourcePackDatas, error) {
	ds := &ResourcePackDatas{
		List: make([]*ResourcePackData, 0, len(bases.DataBases)),
		Map:  make(map[int64]*ResourcePackData, len(bases.DataBases)),
	}
	for _, base := range bases.DataBases {
		d, err := newResourcePackData(base)
		if err != nil {
			return nil, err
		}
		ds.List = append(ds.List, d)
		ds.Map[d.ID] = d
	}
	return ds, nil
}

func newResourcePackData(base *resource_base.PackDataBaseGen) (d *ResourcePackData, err error) {
	d = &ResourcePackData{}

	d.ResourcePackDataGen, err = newResourcePackDataGen(base)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func newResourcePackDataGen(base *resource_base.PackDataBaseGen) (d *ResourcePackDataGen, err error) {
	d = &ResourcePackDataGen{
		PackDataBaseGen: base,
	}

	return d, nil
}

func (ds *ResourcePackDatas) init() error {
	for _, d := range ds.List {
		if err := d.init(); err != nil {
			return err
		}
	}
	return nil
}

func (ds *ResourcePackDatas) build() {
	for _, d := range ds.List {
		d.build()
	}
}

func (d *ResourcePackDataGen) build() {
	if v, err := TryNewItemPrizes(d.Items); err != nil {
		panic(errors.WithMessagef(err, "File:%s Id=%d", d.Table(), d.Id()))
	} else {
		d.ItemsItemPrizes = v
	}

	if v, err := TryNewRadioPrizes(d.Radios); err != nil {
		panic(errors.WithMessagef(err, "File:%s Id=%d", d.Table(), d.Id()))
	} else {
		d.RadiosRadioPrizes = v
	}

	if v, err := TryNewGroupWeights(d.GroupWeight); err != nil {
		panic(errors.WithMessagef(err, "File:%s Id=%d", d.Table(), d.Id()))
	} else {
		d.GroupWeightGroupWeights = v
	}

	if v, err := TryNewItemPrizesList(d.GroupItems); err != nil {
		panic(errors.WithMessagef(err, "File:%s Id=%d", d.Table(), d.Id()))
	} else {
		d.GroupItemsItemPrizesList = v
	}

	if v, err := TryNewRadioPrizesList(d.GroupRadios); err != nil {
		panic(errors.WithMessagef(err, "File:%s Id=%d", d.Table(), d.Id()))
	} else {
		d.GroupRadiosRadioPrizesList = v
	}

}

func (d *ResourcePackDataGen) init() error {
	return nil
}
