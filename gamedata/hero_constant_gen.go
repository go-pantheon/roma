// Code generated by gen-datas. DO NOT EDIT.

package gamedata

import (
	hero_base "github.com/vulcan-frame/vulcan-game/gen/gamedata/base/hero"
	"github.com/pkg/errors"
)

var _ = errors.New("import holding")

// HeroConstantDataGen excel/Hero/hero.xlsx:Constant
type HeroConstantDataGen struct {
	*hero_base.ConstantDataBaseGen

	HeroDataList []*HeroBaseData
}

func newHeroConstantData(base *hero_base.ConstantDataBaseGen) (d *HeroConstantData, err error) {
	d = &HeroConstantData{}

	d.HeroConstantDataGen, err = newHeroConstantDataGen(base)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func newHeroConstantDataGen(base *hero_base.ConstantDataBaseGen) (d *HeroConstantDataGen, err error) {
	d = &HeroConstantDataGen{
		ConstantDataBaseGen: base,
	}

	return d, nil
}

func (d *HeroConstantDataGen) build() {
	for _, id := range d.Hero {
		d.HeroDataList = append(d.HeroDataList, GetHeroBaseData(id))
	}

}

func (d *HeroConstantDataGen) init() error {
	return nil
}
