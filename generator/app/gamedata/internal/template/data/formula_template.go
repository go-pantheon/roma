package data

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/field"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/line"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/vulcan-frame/vulcan-util/camelcase"
)

var formulaTemplate = `
{{- /* delete empty line */ -}}
// Code generated by gen-datas. DO NOT EDIT.

package gamedata

{{- $dataStruct := .DataStruct}}

{{if .HasFields}}
import "github.com/vulcan-frame/vulcan-game/pkg/util/formula"
{{end}}
{{range .Fields}}
func (d *{{$dataStruct}}Gen) Calc{{.Attribute}}(
	{{- range .Params}}
	{{.}} float64,
	{{- end}}
) (float64, error) {
	return formula.Calc(d.Formula, map[string]float64{
		{{- range .Params}}
		"{{.}}": {{.}},
		{{- end}}
	})
}
{{end}}
`

type FormulaService struct {
	Project    string
	Package    string
	DataStruct string
	Fields     []*FormulaField
	HasFields  bool
}

type FormulaField struct {
	Attribute string
	Params    []string
}

func NewFormulaService(project string, sh sheet.Sheet) *FormulaService {
	s := &FormulaService{
		Project: project,
	}
	packageName := sh.GetMetadata().Package
	s.Package = camelcase.ToUnderScore(packageName)
	s.DataStruct = camelcase.ToUpperCamel(sh.GetMetadata().FullName) + "Data"

	sh.WalkFieldMetadata(func(md *field.Metadata) error {
		if md.FormulaValue == "" {
			return nil
		}

		sh.WalkLine(func(l *line.Line) error {
			f := &FormulaField{}
			if ff := l.FieldMap[md.FieldName]; ff != nil {
				f.Params = paramNames(ff.RawValue)
			} else {
				return nil
			}
			if vf := l.FieldMap[md.FormulaValue]; vf != nil {
				f.Attribute = vf.RawValue
			} else {
				return nil
			}

			s.Fields = append(s.Fields, f)
			return nil
		})
		return nil
	})
	s.HasFields = len(s.Fields) > 0

	return s
}

func (s *FormulaService) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("formula").Parse(formulaTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "template new formula error.")
	}
	if err = tmpl.Execute(buf, s); err != nil {
		return nil, errors.Wrapf(err, "template execute formula error.")
	}
	return buf.Bytes(), nil
}

var opers = []string{
	"(",
	")",
	"+",
	"-",
	"*",
	"/",
	"^",
	"%",
}

func paramNames(f string) []string {
	tmp := f
	for _, oper := range opers {
		tmp = strings.ReplaceAll(tmp, oper, " ")
	}
	names := strings.Split(tmp, " ")
	result := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			continue
		}
		if _, err := strconv.ParseFloat(name, 64); err == nil {
			continue
		}
		result = append(result, name)
	}
	return result
}
