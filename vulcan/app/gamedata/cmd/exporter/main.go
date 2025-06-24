package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/gamedata"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
)

var (
	xlsxBaseDir string
	jsonBaseDir string
)

func init() {
	flag.StringVar(&xlsxBaseDir, "xlsx_dir", "excel-source", "xlsx dir path, eg: -xlsx_dir gamedocs/excel-source")
	flag.StringVar(&jsonBaseDir, "json_dir", "json", "json dir path, eg: -json_dir gen/parser/json")
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%+v\n", xsync.CatchErr(p))
			os.Exit(1)
		}
	}()

	flag.Parse()

	xlsxBaseDir = filewriter.BasePath() + xlsxBaseDir
	jsonBaseDir = filewriter.BasePath() + jsonBaseDir
	xlsxBaseDir = filepath.FromSlash(xlsxBaseDir)
	jsonBaseDir = filepath.FromSlash(jsonBaseDir)

	slog.Info("excel directory", "path", xlsxBaseDir)
	slog.Info("json directory", "path", jsonBaseDir)

	json()
	check()
}

func json() {
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
		return true
	})
}

func check() {
	gamedata.Load(jsonBaseDir)
}
