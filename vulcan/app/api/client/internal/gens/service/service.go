package service

import (
	"log/slog"
	"os"
	"path"

	"github.com/go-pantheon/roma/vulcan/app/api/client/internal/template/service"
	"github.com/go-pantheon/roma/vulcan/pkg/compilers"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/pkg/errors"
)

func Gen(project, base string, modCompiler *compilers.ModsCompiler) error {
	dir := path.Join(base, "service/")
	if err := os.Mkdir(dir, 0750); err != nil {
		return errors.Wrapf(err, "failed to create tmp dir: %s", dir)
	}

	return genModService(project, dir, modCompiler)
}

func genModService(project, dir string, c *compilers.ModsCompiler) error {
	s := service.NewSvcService(project, c)
	to := path.Join(dir, "service_gen.go")

	if err := filewriter.GenFile(to, s); err != nil {
		return err
	}

	slog.Info("generate service", "file", filewriter.SprintGenPath(to))

	return nil
}
