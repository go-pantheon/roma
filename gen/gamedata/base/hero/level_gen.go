// Code generated by gen-gamedata. DO NOT EDIT.

package hero

import (
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// LevelDataBaseGens excel/Hero/hero.xlsx:Level
type LevelDataBaseGens struct {
	DataBases []*LevelDataBaseGen
}

// LevelDataBaseGen excel/Hero/hero.xlsx:Level
type LevelDataBaseGen struct {
	ID       int64 // ID
	SubDatas []*LevelSubDataBaseGen
}

// LevelSubDataBaseGen excel/Hero/hero.xlsx:Level
type LevelSubDataBaseGen struct {
	Level int64            // Level
	Name  string           // Name
	Atk   uint64           // Hero Level Atk
	Hp    uint64           // Hero Level Hp
	Cost  map[int64]uint64 // required to upgrade to this level
	Prize map[int64]uint64 // prizes when upgraded
}

func (d *LevelDataBaseGens) Table() string {
	return "excel/Hero/hero.xlsx:Level"
}

func (d *LevelDataBaseGen) Table() string {
	return "excel/Hero/hero.xlsx:Level"
}

func (d *LevelSubDataBaseGen) Table() string {
	return "excel/Hero/hero.xlsx:Level"
}

func LoadLevelDataBaseGens(filename string) *LevelDataBaseGens {
	filename = filepath.FromSlash(filename)

	json, err := os.ReadFile(filename)
	if err != nil {
		panic(errors.Wrapf(err, "Load json failed. file=%s", filename))
	}

	baseList := []LevelDataBaseGen{}
	err = jsoniter.Unmarshal(json, &baseList)
	if err != nil {
		panic(errors.Wrapf(err, "Unmarshal json failed. file=%s", filename))
	}

	datas := &LevelDataBaseGens{}

	for _, base := range baseList {
		datas.DataBases = append(datas.DataBases, &base)
	}

	return datas
}

func (d *LevelDataBaseGen) Id() int64 {
	return d.ID
}

func (d *LevelSubDataBaseGen) Id() int64 {
	return d.Level
}
