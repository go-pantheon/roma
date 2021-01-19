package main

import (
	"flag"
	"log/slog"
	"path"
	"path/filepath"

	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/filewriter"
)

var (
	xlsxBaseDir string
	jsonBaseDir string
)

func init() {
	flag.StringVar(&xlsxBaseDir, "xlsx_dir", "exceldata/excel", "xlsx dir path, eg: -xlsx_dir fixtures/excel")
	flag.StringVar(&jsonBaseDir, "json_dir", "gen/gamedata/json", "json dir path, eg: -json_dir outputs/json")
}

func main() {
	flag.Parse()

	baseDir := filewriter.BasePath()

	xlsxBaseDir = filepath.FromSlash(path.Join(baseDir, xlsxBaseDir))
	jsonBaseDir = filepath.FromSlash(path.Join(baseDir, jsonBaseDir))

	slog.Info("excel directory", "path", xlsxBaseDir)
	slog.Info("json directory", "path", jsonBaseDir)

	if err := filewriter.RebuildDir(jsonBaseDir); err != nil {
		panic(err)
	}

	sheets, err := parser.Parse(xlsxBaseDir)
	if err != nil {
		panic(err)
	}

	sheets.Walk(func(sh sheet.Sheet) bool {
		content, err := sh.EncodeToJson()
		if err != nil {
			panic(err)
		}

		jsonFilePath := filepath.FromSlash(path.Join(jsonBaseDir, sh.FullName()+".json"))
		if err = filewriter.WriteFile(jsonFilePath, []byte(content)); err != nil {
			panic(err)
		}
		slog.Info("json file generated", "path", jsonFilePath)
		return true
	})
}
