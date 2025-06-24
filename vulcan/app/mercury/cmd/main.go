package main

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/app/mercury/internal/template"
	"github.com/go-pantheon/roma/vulcan/pkg/cmd"
	"github.com/go-pantheon/roma/vulcan/pkg/compilers"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/pkg/errors"
)

const (
	project = "github.com/go-pantheon/roma"
)

const (
	apiPathPrefix = "api/client"
	apiModFile    = "module/modules.proto"
	apiSeqDir     = "sequence"
	destPath      = "mercury/gen/task"
	destTmpPath   = "mercury/gen/task_tmp"
)

func main() {
	filewriter.Init(destPath, destTmpPath)

	baseDir := filewriter.BasePath()
	modFile := path.Join(baseDir, apiPathPrefix, apiModFile)
	seqDir := path.Join(baseDir, apiPathPrefix, apiSeqDir)

	destDir := filepath.Join(baseDir, destPath)
	destTmpDir := filepath.Join(baseDir, destTmpPath)

	mcs, err := compilers.NewModCompilers(modFile)
	if err != nil {
		panic(err)
	}

	scs := make([]*compilers.SeqCompiler, 0, 32)
	for _, mc := range mcs {
		for _, mod := range mc.Mods {
			if sc, err := compilers.NewSeqCompilers(path.Join(seqDir, string(mod)+".proto"), mc.Group); err != nil {
				panic(err)
			} else {
				scs = append(scs, sc)
			}
		}
		slog.Info("=== prepare to generate api:", "modules", mc.Mods)
	}

	if err := gen(destTmpDir, destDir, scs); err != nil {
		panic(err)
	}
}

func gen(destTmpDir, destDir string, scs []*compilers.SeqCompiler) error {
	// handle temp directory
	if _, err := os.Stat(destTmpDir); err == nil {
		if err = os.RemoveAll(destTmpDir); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "please manually delete the directory: %s", destTmpPath)
	}
	if err := os.Mkdir(destTmpDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create temp directory: %s", destTmpPath)
	}

	// generate to temp directory
	if err := genTask(destTmpDir, scs); err != nil {
		return err
	}

	if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	if err := os.Rename(destTmpDir, destDir); err != nil {
		return err
	}

	_, err := cmd.CmdExecute(destDir, "gofmt", "-w", destDir)
	if err != nil {
		panic(err)
	}

	slog.Info("=== api files generated.", "dir", destPath)

	return nil
}

func genTask(base string, cs []*compilers.SeqCompiler) error {
	for _, c := range cs {
		dir := base + "/" + string(c.Mod())
		if err := os.Mkdir(dir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create temp directory: %s", dir)
		}

		for _, api := range c.Apis {
			if strings.TrimSpace(api.CS) == "" {
				continue
			}

			s := template.NewTaskService(project, c, api)
			to := dir + "/" + camelcase.ToUnderScore(api.UpperCamelName) + "_gen.go"
			if err := filewriter.GenFile(to, s); err != nil {
				return err
			}
			slog.Info("generated task", "file", filewriter.SprintGenPath(to))
		}
	}
	return nil
}
