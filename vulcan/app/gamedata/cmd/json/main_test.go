package main

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"

	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser"
	"github.com/go-pantheon/roma/vulcan/app/gamedata/internal/parser/sheet"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/stretchr/testify/require"
)

func TestJsonExample(t *testing.T) {
	baseDir := filewriter.BasePath()

	xlsxBaseDir := "../../../../fixtures/excel"
	jsonBaseDir := "../../../../fixtures/outputs/json"

	xlsxBaseDir = filepath.FromSlash(path.Join(baseDir, xlsxBaseDir))
	jsonBaseDir = filepath.FromSlash(path.Join(baseDir, jsonBaseDir))

	t.Log("xlsx dir", "path", xlsxBaseDir)
	t.Log("json dir", "path", jsonBaseDir)

	err := filewriter.RebuildDir(jsonBaseDir)
	require.NoError(t, err)

	sheets, err := parser.Parse(xlsxBaseDir)
	require.NoError(t, err)

	sheets.Walk(func(sh sheet.Sheet) bool {
		content, err := sh.EncodeToJson()
		require.NoError(t, err)
		fmt.Println(content)

		jsonFilePath := filepath.FromSlash(path.Join(jsonBaseDir, sh.FullName()+".json"))
		err = filewriter.WriteFile(jsonFilePath, []byte(content))
		require.NoError(t, err)
		return true
	})
}
