package base

import (
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/field"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/vulcan-frame/vulcan-util/camelcase"
)

type service struct {
	TablePath string

	Package   string
	UpperName string
	LowerName string

	Sub       bool
	Struct    string
	SubStruct string

	IdField    *DataBaseField
	SubIdField *DataBaseField
	Fields     []*DataBaseField
	SubFields  []*DataBaseField
}

type DataBaseField struct {
	Name    string
	Type    string
	Comment string
}

func newService(sh sheet.Sheet) (*service, error) {
	md := sh.GetMetadata()
	s := &service{}
	s.TablePath = md.Path + ":" + md.Sheet
	s.Package = camelcase.ToUnderScore(md.Package)
	s.LowerName = camelcase.ToLowerCamel(md.Name)
	s.UpperName = camelcase.ToUpperCamel(md.Name)
	s.Struct = s.UpperName + "DataBaseGen"
	s.SubStruct = s.UpperName + "SubDataBaseGen"

	if sh.GetIdFieldMetadata() != nil {
		s.IdField = newDataBaseField(sh.GetIdFieldMetadata())
	}

	if sh.GetSubIdFieldMetadata() != nil {
		s.SubIdField = newDataBaseField(sh.GetSubIdFieldMetadata())
	}

	_ = sh.WalkFieldMetadata(func(md *field.Metadata) error {
		s.Fields = append(s.Fields, newDataBaseField(md))
		return nil
	})

	_ = sh.WalkSubFieldMetadata(func(md *field.Metadata) error {
		s.SubFields = append(s.SubFields, newDataBaseField(md))
		return nil
	})

	if len(s.SubFields) > 0 {
		s.Sub = true
	}

	return s, nil
}

func newDataBaseField(md *field.Metadata) *DataBaseField {
	ret := &DataBaseField{
		Name:    camelcase.ToUpperCamel(md.FieldName),
		Type:    md.FieldType.String(),
		Comment: md.Comment,
	}
	return ret
}
