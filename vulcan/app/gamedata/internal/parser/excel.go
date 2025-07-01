package parser

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

var (
	excelFileNameReg = regexp.MustCompile(`^[\p{L}\p{N}][\p{L}\p{N}_-]+\.xlsx$`)
)

type Sheets struct {
	dir    string
	sheets map[string]sheet.Sheet
}

func (gd *Sheets) Walk(f func(gf sheet.Sheet) (continued bool)) {
	for _, s := range gd.sheets {
		if !f(s) {
			break
		}
	}
}

func Parse(path string) (*Sheets, error) {
	dir := &Sheets{
		dir:    path,
		sheets: make(map[string]sheet.Sheet, 256),
	}

	absDir, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "parse directory failed. dir=%s", path)
	}

	err = filepath.Walk(absDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			err = errors.Wrapf(err, "load excel file failed. file=%s", filePath)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !excelFileNameReg.MatchString(filepath.Base(filePath)) {
			slog.Info("ignore file", "path", filePath)
			return nil
		}

		sheets, err := parseFile(filePath)
		if err != nil {
			return err
		}

		for _, s := range sheets {
			if _, ok := dir.sheets[s.FullName()]; ok {
				return errors.Errorf("Sheet duplicated. sheet=%s", s.FullName())
			}

			dir.sheets[s.FullName()] = s
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dir, nil
}

func parseFile(path string) (map[string]sheet.Sheet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		err = errors.Wrapf(err, "Load excel file failed. file=%s", path)
		return nil, err
	}

	sheets, err := sheet.Parse(f)
	if err != nil {
		return nil, errors.WithMessagef(err, "Parse excel file failed. file=%s", path)
	}

	return sheets, nil
}
