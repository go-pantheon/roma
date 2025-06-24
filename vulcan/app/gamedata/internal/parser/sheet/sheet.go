package sheet

import (
	"strings"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/field"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/line"
	"github.com/go-pantheon/roma/vulcan/pkg/align"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type Sheet interface {
	FullName() string
	GetMetadata() *Metadata
	GetIdFieldMetadata() *field.Metadata
	GetSubIdFieldMetadata() *field.Metadata
	WalkLine(func(l *line.Line) error) error
	WalkLineField(fieldName string, f func(l *field.Field) error) error
	WalkFieldMetadata(func(md *field.Metadata) error) error
	WalkSubFieldMetadata(func(md *field.Metadata) error) error
	EncodeToJson() (string, error)
}

func Parse(excel *excelize.File) (sheets map[string]Sheet, err error) {
	sheets = make(map[string]Sheet, len(excel.GetSheetList()))

	kvs, err := parseMetadataSheet(excel)
	if err != nil {
		return
	}

	for _, sheetName := range excel.GetSheetList() {
		t, dataName, err := parseSheetType(sheetName, excel)
		if err != nil {
			return nil, errors.WithMessagef(err, "Parse sheet type failed. sheet=%s", sheetName)
		}
		if t == SheetTypeIgnore {
			continue
		}
		md, err := newMetadata(excel.Path, sheetName, dataName, kvs)
		if err != nil {
			return nil, errors.WithMessagef(err, "Parse sheet metadata failed. sheet=%s", sheetName)
		}

		var s Sheet
		switch t {
		case SheetTypeTable:
			rows, err := excel.GetRows(sheetName, excelize.Options{RawCellValue: true})
			if err != nil {
				return nil, errors.Wrapf(err, "Read excel rows failed. sheet=%s", sheetName)
			}
			for i := range rows {
				rows[i] = rows[i][1:]
			}
			s, err = newTable(md, rows)
			if err != nil {
				return nil, errors.Wrapf(err, "Generate table failed. sheet=%s", sheetName)
			}
		case SheetTypeKV:
			rows, err := excel.GetRows(sheetName, excelize.Options{RawCellValue: true})
			if err != nil {
				return nil, errors.Wrapf(err, "read excel rows failed. sheet=%s", sheetName)
			}
			cols := getCols(rows)
			s, err = newKv(md, cols)
			if err != nil {
				return nil, errors.Wrapf(err, "GSenerate kv failed. sheet=%s", sheetName)
			}
		}
		sheets[dataName] = s
	}
	return
}

func getCols(rows [][]string) [][]string {
	results := make([][]string, field.KvMetadataColSize)
	newRows := make([][]string, 0, len(rows))
	for _, row := range rows {
		newRows = append(newRows, align.Align(row, field.KvMetadataColSize))
	}

	for i, row := range newRows {
		if len(row) != field.KvMetadataColSize {
			break
		}
		if i < field.KvMetadataLineSize {
			continue
		}

		for j, v := range row {
			results[j] = append(results[j], v)
		}
	}

	return results[1:]
}

func parseMetadataSheet(excel *excelize.File) (kvs map[string]string, err error) {
	rows, err := excel.GetRows(MetadataSheetName)
	if err != nil {
		err = errors.Wrapf(err, "read rows failed. sheet=%s", MetadataSheetName)
		return
	}

	kvs = make(map[string]string, len(rows))

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		kvs[row[0]] = row[1]
	}
	return
}

func parseSheetType(sheet string, f *excelize.File) (t SheetType, name string, err error) {
	str, err := f.GetCellValue(sheet, "A1")
	if err != nil {
		err = errors.Wrapf(err, "sheet A1 read faild. sheet=%s", sheet)
		return
	}

	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, string(SheetTypeTable)) {
		name = strings.Replace(str, string(SheetTypeTable), "", 1)
		t = SheetTypeTable
	} else if strings.HasPrefix(str, string(SheetTypeKV)) {
		name = strings.Replace(str, string(SheetTypeKV), "", 1)
		t = SheetTypeKV
	} else {
		t = SheetTypeIgnore
		return
	}
	if name == "" {
		err = errors.Errorf("sheet type format error. sheet=%s, type=%s", sheet, str)
		return
	}
	return
}
