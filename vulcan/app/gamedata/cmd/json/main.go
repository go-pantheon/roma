package main

import (
	"flag"
	"log/slog"
	"path"
	"path/filepath"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
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

	filewriter.Init(jsonBaseDir, jsonBaseDir)

	baseDir := filewriter.BasePath()

	xlsxBasePath := filepath.FromSlash(path.Join(baseDir, xlsxBaseDir))
	jsonBasePath := filepath.FromSlash(path.Join(baseDir, jsonBaseDir))

	slog.Info("=== from excel directory:", "path", xlsxBaseDir)
	slog.Info("=== to json directory:", "path", jsonBaseDir)

	if err := filewriter.RebuildDir(jsonBasePath); err != nil {
		panic(err)
	}

	sheets, err := parser.Parse(xlsxBasePath)
	if err != nil {
		panic(err)
	}

	sheets.Walk(func(sh sheet.Sheet) bool {
		content, err := sh.EncodeToJson()
		if err != nil {
			panic(err)
		}

		jsonFilePath := filepath.FromSlash(path.Join(jsonBasePath, sh.FullName()+".json"))
		if err = filewriter.WriteFile(jsonFilePath, []byte(content)); err != nil {
			panic(err)
		}

		slog.Info("generated json", "path", filewriter.SprintGenPath(jsonFilePath))

		return true
	})

	slog.Info("=== json files generated.", "dir", jsonBaseDir)
}
