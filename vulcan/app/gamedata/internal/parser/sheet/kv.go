package sheet

import (
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/line"
	"github.com/pkg/errors"
)

var _ Sheet = (*Kv)(nil)

type Kv struct {
	*Metadata

	fieldMds []*field.Metadata
	line     *line.Line
}

func newKv(tmd *Metadata, rows [][]string) (*Kv, error) {
	tmd.Type = SheetTypeKV

	kv := &Kv{
		Metadata: tmd,
	}

	mds, err := field.NewMetadataList(rows[:field.TableMetadataLineSize])
	if err != nil {
		return nil, err
	}

	l, err := line.NewLine(mds, rows[field.TableMetadataLineSize])
	if err != nil {
		return nil, err
	}

	if len(l.SubLines) > 0 {
		return nil, errors.New("kv sheet cannot have sub lines")
	}

	kv.fieldMds = mds
	kv.line = l

	for _, md := range mds {
		if md.FormulaValue != "" {
			kv.HasFormulaField = true
			break
		}
	}

	return kv, nil
}

func (kv *Kv) FullName() string {
	return kv.Metadata.FullName
}

func (kv *Kv) GetMetadata() *Metadata {
	return kv.Metadata
}

func (kv *Kv) GetIdFieldMetadata() *field.Metadata {
	return nil
}

func (kv *Kv) GetSubIdFieldMetadata() *field.Metadata {
	return nil
}

func (kv *Kv) EncodeToJson() (string, error) {
	json, err := kv.line.EncodeToJson()
	if err != nil {
		return "", errors.WithMessagef(err, "kv=%s", kv.FullName())
	}

	return json, nil
}

func (kv *Kv) WalkLine(f func(l *line.Line) (err error)) error {
	return f(kv.line)
}

func (kv *Kv) WalkLineField(fieldName string, f func(l *field.Field) error) error {
	field := kv.line.FieldMap[fieldName]
	if field == nil {
		return nil
	}

	return f(field)
}

func (kv *Kv) WalkFieldMetadata(f func(md *field.Metadata) error) error {
	for _, md := range kv.fieldMds {
		if err := f(md); err != nil {
			return err
		}
	}

	return nil
}

func (kv *Kv) WalkSubFieldMetadata(_ func(md *field.Metadata) error) error {
	return nil
}
