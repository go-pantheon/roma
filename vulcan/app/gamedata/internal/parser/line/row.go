package line

import (
	"reflect"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/go-pantheon/roma/vulcan/pkg/align"
	"github.com/pkg/errors"
)

type RowGroup struct {
	idField   *field.Field
	subIdName string
	rows      []*Row
}

type Row struct {
	IdField    *field.Field
	SubIdField *field.Field
	FieldMap   map[string]*field.Field
	Fields     []*field.Field
}

func newRowGroup(idField *field.Field, subIdField *field.Field, rows []*Row) *RowGroup {
	rg := &RowGroup{
		idField: idField,
		rows:    rows,
	}

	if subIdField != nil {
		rg.subIdName = subIdField.Name
	}
	return rg
}

// newRows groups rows based on ID and SubID like Rows[ID][SubID]
func newRows(mds []*field.Metadata, values [][]string) (rows []*RowGroup, err error) {
	rows = make([]*RowGroup, 0, len(values))

	var (
		rowGroups   = make(map[int64]*RowGroup)
		rowIdOrders = make([]int64, 0, len(values))
	)

	for _, vs := range values {
		if len(vs) <= 1 {
			continue
		}

		var row *Row
		row, err = newRow(mds, vs)
		if err != nil {
			err = errors.WithMessagef(err, "values<%+v>", vs)
			return
		}

		id := row.id()

		if _, ok := rowGroups[id]; !ok {
			rowIdOrders = append(rowIdOrders, id)
		}

		rg := rowGroups[id]
		if rg == nil {
			rg = newRowGroup(row.IdField, row.SubIdField, []*Row{})
			rowGroups[id] = rg
		}
		rg.rows = append(rg.rows, row)
	}

	for _, id := range rowIdOrders {
		rg := rowGroups[id]
		subIdMap := make(map[int64]struct{}, len(rg.rows))
		for _, row := range rg.rows {
			subId, _ := row.subId()
			if _, ok := subIdMap[subId]; ok {
				err = errors.Errorf("sub id already exists. sub id=%d", subId)
				return
			}
			subIdMap[subId] = struct{}{}
		}
		rows = append(rows, rg)
	}

	return
}

func newRow(mds []*field.Metadata, values []string) (row *Row, err error) {
	values = align.Align(values, len(mds))

	row = &Row{
		FieldMap: make(map[string]*field.Field, len(mds)),
		Fields:   make([]*field.Field, 0, len(mds)),
	}

	for i, md := range mds {
		f, e := field.NewField(md, values[i])
		if e != nil {
			err = errors.WithMessagef(e, "field=%s", md.FieldName)
			return
		}

		switch md.Type {
		case field.IdType, field.SharedIdType:
			if row.IdField != nil {
				err = errors.Errorf("id field already exists. field=%s", md.FieldName)
				return
			}
			if md.FieldType != reflect.TypeOf(int64(0)) {
				err = errors.Errorf("id field type must be int64. field=%s", md.FieldName)
				return
			}
			if f.Value.(int64) <= 0 {
				err = errors.Errorf("id field value cannot be 0. field=%s", md.FieldName)
				return
			}
			row.IdField = f
		case field.SharedSubIdType:
			if row.SubIdField != nil {
				err = errors.Errorf("sub id field already exists. field=%s", md.FieldName)
				return
			}
			if md.FieldType != reflect.TypeOf(int64(0)) {
				err = errors.Errorf("sub id field type must be int64. field=%s", md.FieldName)
				return
			}
			if f.Value.(int64) <= 0 {
				err = errors.Errorf("sub id field value cannot be 0. field=%s", md.FieldName)
				return
			}
			row.SubIdField = f
		default:
			if _, ok := row.FieldMap[md.FieldName]; ok {
				err = errors.Errorf("field already exists. field=%s", md.FieldName)
				return
			}
			row.FieldMap[md.FieldName] = f
			row.Fields = append(row.Fields, f)
		}
	}

	if row.IdField == nil {
		err = errors.New("id field not found")
		return
	}

	return
}

func (r *Row) id() int64 {
	return r.IdField.Value.(int64)
}

func (r *Row) subId() (v int64, ok bool) {
	if r.SubIdField == nil {
		return 0, false
	}
	return r.SubIdField.Value.(int64), true
}
