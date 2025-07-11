package base

import (
	"path/filepath"
	"strings"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
)

type service struct {
	Org       string
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

func newService(project string, sh sheet.Sheet) *service {
	md := sh.GetMetadata()
	s := &service{}
	s.Org = filepath.Clean(strings.Replace(project, filepath.Base(project), "", 1))
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

	return s
}

func newDataBaseField(md *field.Metadata) *DataBaseField {
	ret := &DataBaseField{
		Name:    camelcase.ToUpperCamel(md.FieldName),
		Type:    md.FieldType.String(),
		Comment: md.Comment,
	}

	return ret
}
