// Code generated by gen-gamedata. DO NOT EDIT.

package hero

import (
	"os"
	"path/filepath"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// SkillDataBaseGens excel/Hero/hero.xlsx:Skill
type SkillDataBaseGens struct {
	DataBases []*SkillDataBaseGen
}

// SkillDataBaseGen excel/Hero/hero.xlsx:Skill
type SkillDataBaseGen struct {
	ID       int64  // ID
	Name     string // Name
	SubDatas []*SkillSubDataBaseGen
}

// SkillSubDataBaseGen excel/Hero/hero.xlsx:Skill
type SkillSubDataBaseGen struct {
	Level int64  // Type
	Atk   uint64 // Atk
}

func (d *SkillDataBaseGens) Table() string {
	return "excel/Hero/hero.xlsx:Skill"
}

func (d *SkillDataBaseGen) Table() string {
	return "excel/Hero/hero.xlsx:Skill"
}

func (d *SkillSubDataBaseGen) Table() string {
	return "excel/Hero/hero.xlsx:Skill"
}

func LoadSkillDataBaseGens(filename string) *SkillDataBaseGens {
	filename = filepath.FromSlash(filename)

	json, err := os.ReadFile(filename)
	if err != nil {
		panic(errors.Wrapf(err, "Load json failed. file=%s", filename))
	}

	baseList := []SkillDataBaseGen{}
	err = jsoniter.Unmarshal(json, &baseList)
	if err != nil {
		panic(errors.Wrapf(err, "Unmarshal json failed. file=%s", filename))
	}

	datas := &SkillDataBaseGens{}

	for _, base := range baseList {
		datas.DataBases = append(datas.DataBases, &base)
	}

	return datas
}

func (d *SkillDataBaseGen) Id() int64 {
	return d.ID
}

func (d *SkillSubDataBaseGen) Id() int64 {
	return d.Level
}
