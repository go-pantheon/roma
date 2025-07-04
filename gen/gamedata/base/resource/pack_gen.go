// Code generated by gen-data-base. DO NOT EDIT.

package resource

import (
	"os"
	"path/filepath"

	"github.com/go-pantheon/fabrica-util/errors"
	jsoniter "github.com/json-iterator/go"
)

// PackDataBaseGens excel/Resource/item.xlsx:Pack
type PackDataBaseGens struct {
	DataBases []*PackDataBaseGen
}

// PackDataBaseGen excel/Resource/item.xlsx:Pack
type PackDataBaseGen struct {
	ID          int64              // ID
	Name        string             // Name
	ItemTypeInt uint64             // Type
	Max         uint64             // Max Amount. Unlimited if Zero
	Items       map[int64]uint64   // Item list
	Radios      map[int64]uint64   // Radio List
	GroupItems  []map[int64]uint64 // Group Item Array
	GroupRadios []map[int64]uint64 // Group Radio Array
	GroupWeight []uint64           // Group Random Weight
}

func (d *PackDataBaseGens) Table() string {
	return "excel/Resource/item.xlsx:Pack"
}

func (d *PackDataBaseGen) Table() string {
	return "excel/Resource/item.xlsx:Pack"
}

func LoadPackDataBaseGens(filename string) *PackDataBaseGens {
	filename = filepath.FromSlash(filename)

	json, err := os.ReadFile(filename)
	if err != nil {
		panic(errors.Wrapf(err, "Load json failed. file=%s", filename))
	}

	baseList := []PackDataBaseGen{}
	err = jsoniter.Unmarshal(json, &baseList)
	if err != nil {
		panic(errors.Wrapf(err, "Unmarshal json failed. file=%s", filename))
	}

	datas := &PackDataBaseGens{}

	for _, base := range baseList {
		datas.DataBases = append(datas.DataBases, &base)
	}

	return datas
}

func (d *PackDataBaseGen) Id() int64 {
	return d.ID
}
