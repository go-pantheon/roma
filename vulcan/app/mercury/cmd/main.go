package main

import (
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/app/mercury/internal/template"
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
	modPath := path.Join(baseDir, apiPathPrefix, apiModFile)
	seqDirPath := path.Join(baseDir, apiPathPrefix, apiSeqDir)

	mcs, err := compilers.NewModCompilers(modPath)
	if err != nil {
		panic(err)
	}

	scs := make([]*compilers.SeqCompiler, 0, 32)
	for _, mc := range mcs {
		for _, mod := range mc.Mods {
			if sc, err := compilers.NewSeqCompilers(path.Join(seqDirPath, string(mod)+".proto"), mc.Group); err != nil {
				panic(err)
			} else {
				scs = append(scs, sc)
			}
		}
		slog.Info("=== prepare to generate api:", "modules", mc.Mods)
	}

	if err := gen(baseDir, scs); err != nil {
		panic(err)
	}
}

func gen(_ string, scs []*compilers.SeqCompiler) error {
	// handle temp directory
	if _, err := os.Stat(destTmpPath); err == nil {
		if err = os.RemoveAll(destTmpPath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "please manually delete the directory: %s", destTmpPath)
	}
	if err := os.Mkdir(destTmpPath, 0755); err != nil {
		return errors.Wrapf(err, "failed to create temp directory: %s", destTmpPath)
	}

	// generate to temp directory
	if err := genTask(destTmpPath, scs); err != nil {
		return err
	}

	if err := os.RemoveAll(destPath); err != nil {
		return err
	}
	if err := os.Rename(destTmpPath, destPath); err != nil {
		return err
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
