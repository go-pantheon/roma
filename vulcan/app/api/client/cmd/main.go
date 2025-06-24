package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/go-pantheon/roma/vulcan/app/api/client/internal/gens/codec"
	"github.com/go-pantheon/roma/vulcan/app/api/client/internal/gens/handler"
	"github.com/go-pantheon/roma/vulcan/app/api/client/internal/gens/service"
	"github.com/go-pantheon/roma/vulcan/pkg/cmd"
	"github.com/go-pantheon/roma/vulcan/pkg/compilers"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/pkg/errors"
)

const (
	project = "github.com/go-pantheon/roma"
)

const (
	apiPathPrefix = "api/client/"
	apiModFile    = "module/modules.proto"
	apiSeqDir     = "sequence/"
	destPath      = "gen/app"
	destTmpPath   = "gen/app_tmp"
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
			seqFile := path.Join(seqDirPath, fmt.Sprintf("%s.proto", mod))
			if sc, err := compilers.NewSeqCompilers(seqFile, mc.Group); err != nil {
				slog.Error("new seq compiler failed", "mod", mod, "error", err)
			} else {
				scs = append(scs, sc)
			}
		}
		slog.Info("=== prepare to generate api:", "mods", mc.Mods)
	}

	if err = gen(baseDir, mcs, scs); err != nil {
		slog.Error("generate failed", "error", err)
		panic(err)
	}
}

func gen(base string, mcs []*compilers.ModsCompiler, scs []*compilers.SeqCompiler) error {
	tmpDir := path.Join(base, destTmpPath)
	dir := path.Join(base, destPath)

	// rebuild tmp dir
	if _, err := os.Stat(tmpDir); err == nil {
		if err = os.RemoveAll(tmpDir); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to remove tmp dir: %s", tmpDir)
	}
	if err := os.Mkdir(tmpDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create tmp dir: %s", tmpDir)
	}

	// generate to tmp dir
	if err := codec.Gen(project, tmpDir, mcs, scs); err != nil {
		return err
	}

	for _, mc := range mcs {
		tmpGroupDir := path.Join(tmpDir, string(mc.Group))
		if err := os.Mkdir(tmpGroupDir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create tmp group dir: %s", tmpGroupDir)
		}
		// generate to tmp dir
		if err := service.Gen(project, tmpGroupDir, mc); err != nil {
			return err
		}
		if err := handler.Gen(project, tmpGroupDir, mc, scs); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if err := os.Rename(tmpDir, dir); err != nil {
		return err
	}

	_, err := cmd.CmdExecute(dir, "gofmt", "-w", dir)
	if err != nil {
		panic(err)
	}

	slog.Info("=== api files generated.", "dir", destPath)

	return nil
}
