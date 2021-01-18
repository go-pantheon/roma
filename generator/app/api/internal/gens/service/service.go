package service

import (
	"log/slog"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/api/internal/template/service"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/compilers"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/filewriter"
)

func Gen(project, base string, modCompiler *compilers.ModsCompiler) error {
	dir := path.Join(base, "service/")
	if err := os.Mkdir(dir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create tmp dir: %s", dir)
	}

	if err := genModService(project, dir, modCompiler); err != nil {
		return err
	}
	return nil
}

func genModService(project, dir string, c *compilers.ModsCompiler) error {
	s := service.NewSvcService(project, c)
	to := path.Join(dir, "service_gen.go")
	if err := filewriter.GenFile(to, s); err != nil {
		return err
	}
	slog.Info("generate mod services files completed", "to", to)
	return nil
}
