package sheet

import (
	"fmt"
	"strings"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/line"
	"github.com/pkg/errors"
)

var _ Sheet = (*Table)(nil)

type Table struct {
	*Metadata

	idFieldMetadata    *field.Metadata
	fieldMetadatas     []*field.Metadata
	subIdFieldMetadata *field.Metadata
	subFieldMetadatas  []*field.Metadata

	lines []*line.Line
}

func newTable(md *Metadata, rows [][]string) (*Table, error) {
	md.Type = SheetTypeTable

	o := &Table{
		Metadata: md,
	}

	mds, err := field.NewMetadataList(rows[:field.TableMetadataLineSize])
	if err != nil {
		return nil, err
	}

	lines, err := line.NewLines(mds, rows[field.TableMetadataLineSize:])
	if err != nil {
		return nil, err
	}
	o.lines = lines

	if len(lines) == 0 {
		return o, nil
	}

	l := lines[0]
	o.idFieldMetadata = l.IdField.Metadata
	for _, f := range l.Fields {
		o.fieldMetadatas = append(o.fieldMetadatas, f.Metadata)
	}
	if len(l.SubLines) > 0 {
		sl := l.SubLines[0]
		o.subIdFieldMetadata = sl.IdField.Metadata
		for _, f := range sl.Fields {
			o.subFieldMetadatas = append(o.subFieldMetadatas, f.Metadata)
		}
	}

	for _, md := range mds {
		if md.FormulaValue != "" {
			o.Metadata.HasFormulaField = true
			break
		}
	}

	return o, nil
}

func (t *Table) FullName() string {
	return t.Metadata.FullName
}

func (t *Table) GetMetadata() *Metadata {
	return t.Metadata
}

func (t *Table) GetIdFieldMetadata() *field.Metadata {
	return t.idFieldMetadata
}

func (t *Table) GetSubIdFieldMetadata() *field.Metadata {
	return t.subIdFieldMetadata
}

func (t *Table) EncodeToJson() (string, error) {
	var parts []string

	for _, l := range t.lines {
		jsonData, err := l.EncodeToJson()
		if err != nil {
			return "", errors.WithMessagef(err, "table=%s", t.FullName())
		}
		parts = append(parts, jsonData)
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ",")), nil
}

func (t *Table) WalkLine(f func(l *line.Line) (err error)) error {
	for _, l := range t.lines {
		if err := f(l); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) WalkLineField(fieldName string, f func(l *field.Field) error) error {
	for _, l := range t.lines {
		field := l.FieldMap[fieldName]
		if field == nil {
			continue
		}
		if err := f(field); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) WalkFieldMetadata(f func(md *field.Metadata) error) error {
	for _, md := range t.fieldMetadatas {
		if err := f(md); err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) WalkSubFieldMetadata(f func(md *field.Metadata) error) error {
	for _, md := range t.subFieldMetadatas {
		if err := f(md); err != nil {
			return err
		}
	}
	return nil
}
