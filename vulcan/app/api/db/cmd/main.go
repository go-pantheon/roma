package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-pantheon/roma/vulcan/app/api/db/internal/gen/pool"
	tmpl "github.com/go-pantheon/roma/vulcan/app/api/db/internal/template"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
)

const (
	project = "github.com/go-pantheon/roma"
)

const (
	apiPathPrefix = "api/db/"
	destPath      = "gen"
)

func main() {
	baseDir := filewriter.BasePath()
	destDir := filepath.Join(baseDir, destPath)

	data, err := pool.ParseData(apiPathPrefix, project)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		panic(err)
	}

	log.Printf("dir: %s", destDir)

	service := tmpl.NewService()
	for _, file := range data.Files {
		formatted, err := service.Execute(file)
		if err != nil {
			panic(err)
		}

		destFilePath := filepath.Join(destDir, file.Dir, file.FileName+"_pool.go")

		if err := os.WriteFile(destFilePath, formatted, 0644); err != nil {
			panic(err)
		}

	}
}
