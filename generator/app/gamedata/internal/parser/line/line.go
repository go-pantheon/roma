package line

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	field "github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/field"
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
		subId, ok := r.subId()
		if ok {
			sl := newSubLine(r.SubIdField)
			if err := l.addSubLine(sl); err != nil {
				return nil, err
			}
		}
		for _, f := range r.Fields {
			switch f.Type {
			case field.SharedType:
				l.addSharedField(f.Metadata, f)
			case field.MergedListType, field.MergedListNonNilType, field.MergedMapType:
				if _, ok := mergeFieldNames[f.Metadata.FieldName]; !ok {
					mergeFields = append(mergeFields, f.Metadata.FieldName)
					mergeFieldNames[f.Metadata.FieldName] = struct{}{}
				}
			default:
				if ok {
					if err := l.SubLineMap[subId].addField(f.Metadata, f); err != nil {
						return nil, err
					}
					continue
				}
				if err := l.addField(f.Metadata, f); err != nil {
					return nil, err
				}
			}
		}
	}

	for _, name := range mergeFields {
		fds := make([]*field.Field, 0, len(rg.rows))
		for _, row := range rg.rows {
			fds = append(fds, row.FieldMap[name])
		}
		mf, err := field.Merge(fds)
		if err != nil {
			return nil, err
		}
		if err := l.addField(mf.Metadata, mf); err != nil {
			return nil, err
		}
	}

	if len(l.Fields) == 0 && len(l.SubLines) == 0 {
		return nil, errors.Errorf("line fields cannot be empty. id=%d", l.Id())
	}

	// remove empty sublines
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

	return l, nil
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

	var parts []string

	idJson, err := l.IdField.EncodeToJson()
	if err != nil {
		return "", errors.Wrapf(err, "line=%d encode id failed", l.Id())
	}
	parts = append(parts, idJson)

	for _, f := range l.Fields {
		jsonData, err := f.EncodeToJson()
		if err != nil {
			return "", errors.Wrapf(err, "field=%s", f.Metadata.FieldName)
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
