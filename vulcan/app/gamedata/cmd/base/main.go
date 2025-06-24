package main

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/template/base"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/template/baseload"
	"github.com/go-pantheon/roma/vulcan/pkg/cmd"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
)

const (
	project = "github.com/go-pantheon/roma"
)

const (
	excelPath   = "exceldata/excel/"
	destPath    = "gen/gamedata/base/"
	destTmpPath = "gen/gamedata/base_tmp/"
)

func main() {
	filewriter.Init(destPath, destTmpPath)

	baseDir := filewriter.BasePath()

	slog.Info("=== from excel directory:", "path", excelPath)
	slog.Info("=== to base directory:", "path", destPath)

	destDir := path.Join(baseDir, destPath)
	tmpDir := path.Join(baseDir, destTmpPath)
	if err := filewriter.RebuildDir(tmpDir); err != nil {
		panic(err)
	}

	excelDir := filepath.Join(baseDir, excelPath)
	sheets, err := parser.Parse(excelDir)
	if err != nil {
		panic(err)
	}

	if err := genLoader(tmpDir, sheets); err != nil {
		panic(err)
	}

	sheets.Walk(func(sh sheet.Sheet) bool {
		if err := genBase(tmpDir, sh); err != nil {
			panic(err)
		}
		return true
	})

	// move files from tmpDir to destDir
	if err := os.RemoveAll(destDir); err != nil {
		panic(err)
	}

	if err := os.Rename(tmpDir, destDir); err != nil {
		panic(err)
	}

	_, err = cmd.CmdExecute(destDir, "gofmt", "-w", destDir)
	if err != nil {
		panic(err)
	}

	slog.Info("=== Base directory generated.", "dir", destPath)
}

func genBase(dir string, sh sheet.Sheet) (err error) {
	md := sh.GetMetadata()

	destDir := path.Join(dir, camelcase.ToUnderScore(md.Package))
	if _, err = filewriter.CreateDir(destDir); err != nil {
		return
	}

	var s filewriter.GenService
	switch md.Type {
	case sheet.SheetTypeTable:
		if s, err = base.NewTableService(project, sh); err != nil {
			return
		}
	case sheet.SheetTypeKV:
		if s, err = base.NewKvService(project, sh); err != nil {
			return
		}
	}

	to := path.Join(destDir, camelcase.ToUnderScore(md.Name)+"_gen.go")
	if err = filewriter.GenFile(to, s); err != nil {
		return
	}

	slog.Info("generated base", "file", filewriter.SprintGenPath(to))
	return
}

func genLoader(dir string, excelDir *parser.Sheets) (err error) {
	to := filepath.Join(dir, "loader_gen.go")
	s := baseload.NewService(project, excelDir)
	if err = filewriter.GenFile(to, s); err != nil {
		return
	}

	slog.Info("generated base loader", "file", filewriter.SprintGenPath(to))
	return
}
