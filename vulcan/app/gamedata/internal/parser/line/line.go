package line

import (
	"fmt"
	"strings"

	field "github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/pkg/errors"
)

const (
	JsonSubLineName = "SubDatas"
)

type Line struct {
	IdField *field.Field

	Fields   []*field.Field
	FieldMap map[string]*field.Field

	SubLines   []*Line
	SubLineMap map[int64]*Line
}

func NewLines(mds []*field.Metadata, values [][]string) ([]*Line, error) {
	rows, err := newRows(mds, values)
	if err != nil {
		return nil, err
	}

	lines := make([]*Line, 0, len(rows))

	for _, rs := range rows {
		line, err := newLine(mds, rs)
		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func NewLine(mds []*field.Metadata, values []string) (*Line, error) {
	row, err := newRow(mds, values)
	if err != nil {
		return nil, err
	}

	return newLine(mds, newRowGroup(row.IdField, row.SubIdField, []*Row{row}))
}

func newSubLine(subIdField *field.Field) *Line {
	sl := &Line{
		IdField:  subIdField,
		FieldMap: make(map[string]*field.Field),
	}

	return sl
}

func newLine(mds []*field.Metadata, rg *RowGroup) (*Line, error) {
	if rg == nil || len(rg.rows) == 0 {
		return nil, errors.Errorf("rows is empty")
	}

	l := &Line{
		IdField:    rg.idField,
		FieldMap:   make(map[string]*field.Field, len(mds)),
		SubLineMap: make(map[int64]*Line, len(rg.rows)),
	}

	mergeFieldNames := make(map[string]struct{}, len(mds))
	mergeFields := make([]string, 0, len(mds))

	for _, r := range rg.rows {
		subId, isSubLine := r.subId()
		if isSubLine {
			if err := l.addSubLine(newSubLine(r.SubIdField)); err != nil {
				return nil, err
			}
		}

		for _, f := range r.Fields {
			mergeField, err := praseField(l, isSubLine, subId, f)
			if err != nil {
				return nil, err
			}

			if mergeField != "" {
				if _, ok := mergeFieldNames[mergeField]; !ok {
					mergeFields = append(mergeFields, mergeField)
					mergeFieldNames[mergeField] = struct{}{}
				}
			}
		}
	}

	for _, name := range mergeFields {
		if err := l.addMergeField(rg, name); err != nil {
			return nil, err
		}
	}

	if len(l.Fields) == 0 && len(l.SubLines) == 0 {
		return nil, errors.Errorf("line fields cannot be empty. id=%d", l.Id())
	}

	l.removeEmptySubLines()

	return l, nil
}

func praseField(l *Line, isSubLine bool, subId int64, f *field.Field) (mergeField string, err error) {
	if f.IsMerged() {
		return f.FieldName, nil
	}

	if f.Type == field.SharedType {
		l.addSharedField(f.Metadata, f)
		return "", nil
	}

	if isSubLine {
		return "", l.SubLineMap[subId].addField(f.Metadata, f)
	}

	return "", l.addField(f.Metadata, f)
}

func (l *Line) addMergeField(rg *RowGroup, name string) error {
	fds := make([]*field.Field, 0, len(rg.rows))

	for _, row := range rg.rows {
		if f, ok := row.FieldMap[name]; ok && f.IsMerged() {
			fds = append(fds, f)
		} else {
			return errors.Errorf("field not found or not merged. field=%s", name)
		}
	}

	mf, err := field.Merge(fds)
	if err != nil {
		return err
	}

	return l.addField(mf.Metadata, mf)
}

func (l *Line) removeEmptySubLines() {
	newSubLines := make([]*Line, 0, len(l.SubLines))
	newSubLineMap := make(map[int64]*Line, len(l.SubLines))

	for _, sl := range l.SubLines {
		if len(sl.Fields) > 0 {
			newSubLines = append(newSubLines, sl)
			newSubLineMap[sl.Id()] = sl
		}
	}

	l.SubLines = newSubLines
	l.SubLineMap = newSubLineMap
}

func (l *Line) addSharedField(md *field.Metadata, fd *field.Field) {
	if field.IsZeroValue(fd.Value) {
		return
	}

	l.Fields = append(l.Fields, fd)
	l.FieldMap[md.FieldName] = fd
}

func (l *Line) addField(md *field.Metadata, fd *field.Field) error {
	if _, ok := l.FieldMap[md.FieldName]; ok {
		return errors.Errorf("field name duplicated. field=%s name1=%s name2=%s", md.FieldName, md.Name, l.FieldMap[md.FieldName].Name)
	}

	l.Fields = append(l.Fields, fd)
	l.FieldMap[md.FieldName] = fd

	return nil
}

func (l *Line) addSubLine(sl *Line) error {
	if _, ok := l.SubLineMap[sl.Id()]; ok {
		return errors.Errorf("subId duplicated. subId=%d", sl.Id())
	}

	l.SubLines = append(l.SubLines, sl)
	l.SubLineMap[sl.Id()] = sl

	return nil
}

func (l *Line) Id() int64 {
	return l.IdField.Value.(int64)
}

func (l *Line) EncodeToJson() (string, error) {
	if len(l.Fields) == 0 && len(l.SubLines) == 0 {
		return "", errors.Errorf("fields and sublines cannot be empty. id=%d", l.Id())
	}

	parts := make([]string, 0, len(l.Fields)+len(l.SubLines))

	idJson, err := l.IdField.EncodeToJson()
	if err != nil {
		return "", errors.Wrapf(err, "line=%d encode id failed", l.Id())
	}

	parts = append(parts, idJson)

	for _, f := range l.Fields {
		jsonData, err := f.EncodeToJson()
		if err != nil {
			return "", errors.Wrapf(err, "field=%s", f.FieldName)
		}

		parts = append(parts, jsonData)
	}

	subJsonBuilders := &strings.Builder{}

	if len(l.SubLines) > 0 {
		for _, sl := range l.SubLines {
			sj, err := sl.EncodeToJson()
			if err != nil {
				return "", errors.WithMessagef(err, "subline=%d", sl.Id())
			}

			subJsonBuilders.WriteString(sj)
			subJsonBuilders.WriteRune(',')
		}

		subJson := fmt.Sprintf("[%s]", strings.TrimSuffix(subJsonBuilders.String(), ","))

		parts = append(parts, fmt.Sprintf(`"%s": %s`, JsonSubLineName, subJson))
	}

	return fmt.Sprintf("{%s}", strings.Join(parts, ",")), nil
}
