package main

import (
	"log/slog"
	"os"
	"path"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/template/data"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/template/dataload"
	"github.com/go-pantheon/roma/vulcan/pkg/cmd"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/pkg/errors"
)

const (
	project = "github.com/go-pantheon/roma"
)

const (
	excelPath = "exceldata/excel/"
	destPath  = "gamedata/"
)

func main() {
	filewriter.Init(destPath, destPath)

	baseDir := filewriter.BasePath()

	slog.Info("=== from excel directory:", "path", excelPath)
	slog.Info("=== to gamedata directory:", "path", destPath)

	destDir := path.Join(baseDir, destPath)

	excelDir := path.Join(baseDir, excelPath)
	sheets, err := parser.Parse(excelDir)
	if err != nil {
		panic(err)
	}

	shs := []sheet.Sheet{}
	sheets.Walk(func(sh sheet.Sheet) bool {
		if err := genData(destDir, sh); err != nil {
			panic(err)
		}
		if err := genSpecialFormula(destDir, sh); err != nil {
			panic(err)
		}
		shs = append(shs, sh)
		return true
	})

	if err := genLoader(destDir, shs); err != nil {
		panic(err)
	}

	_, err = cmd.CmdExecute(destDir, "gofmt", "-w", destDir)
	if err != nil {
		panic(err)
	}
	slog.Info("=== gamedata directory generated.", "dir", destPath)
}

func genData(dir string, sh sheet.Sheet) (err error) {
	to := path.Join(dir, camelcase.ToUnderScore(sh.GetMetadata().FullName)+"_gen.go")
	if err = removeFile(to); err != nil {
		return
	}

	s := data.NewDataService(project, sh)
	if err = filewriter.GenFile(to, s); err != nil {
		return
	}
	slog.Info("generated data", "file", filewriter.SprintGenPath(to))
	return
}

func genSpecialFormula(dir string, sh sheet.Sheet) (err error) {
	if !sh.GetMetadata().HasFormulaField {
		return
	}

	to := path.Join(dir, camelcase.ToUnderScore(sh.GetMetadata().FullName)+"_formula_gen.go")
	if err = removeFile(to); err != nil {
		return
	}

	s := data.NewFormulaService(project, sh)
	if err = filewriter.GenFile(to, s); err != nil {
		return
	}
	slog.Info("generated data special formula", "file", filewriter.SprintGenPath(to))
	return
}

func removeFile(filePath string) (err error) {
	if _, err = os.Stat(filePath); err == nil {
		if err = os.Remove(filePath); err != nil {
			return errors.Wrapf(err, "failed to remove file: %s", filePath)
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "please manually delete the file: %s", filePath)
	}
	return nil
}

func genLoader(dir string, shs []sheet.Sheet) (err error) {
	to := path.Join(dir, "loader_gen.go")
	if err = removeFile(to); err != nil {
		return
	}
	s := dataload.NewService(project, shs)
	if err = filewriter.GenFile(to, s); err != nil {
		return
	}
	slog.Info("generated data loader", "file", filewriter.SprintGenPath(to))
	return
}
