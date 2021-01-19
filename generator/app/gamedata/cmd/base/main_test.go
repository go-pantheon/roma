package main

import (
	"log/slog"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/filewriter"
)

func TestMain(t *testing.T) {
	baseDir := filewriter.BasePath()

	excelPath := "../../../../fixtures/excel"
	destPath := "../../../../fixtures/outputs/base"

	excelPath = filepath.FromSlash(path.Join(baseDir, excelPath))
	destPath = filepath.FromSlash(path.Join(baseDir, destPath))

	slog.Info("excel directory", "dir", excelPath)
	slog.Info("dest directory", "dir", destPath)

	err := filewriter.RebuildDir(destPath)
	require.NoError(t, err)

	excel, err := parser.Parse(excelPath)
	require.NoError(t, err)

	err = genLoader(destPath, excel)
	require.NoError(t, err)

	excel.Walk(func(sh sheet.Sheet) bool {
		err = genBase(destPath, sh)
		require.NoError(t, err)
		return true
	})
}
