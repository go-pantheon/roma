package handler

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/go-pantheon/roma/vulcan/app/api/client/internal/template/handler"
	"github.com/go-pantheon/roma/vulcan/pkg/compilers"
	"github.com/go-pantheon/roma/vulcan/pkg/filewriter"
	"github.com/pkg/errors"
)

func Gen(project, base string, mc *compilers.ModsCompiler, scs []*compilers.SeqCompiler) error {
	dir := path.Join(base, "handler/")
	if err := os.Mkdir(dir, 0750); err != nil {
		return errors.Wrapf(err, "failed to create tmp dir: %s", dir)
	}

	if err := genCodecHandler(project, dir, mc); err != nil {
		return err
	}

	if err := genModsHandler(project, dir, mc, scs); err != nil {
		return err
	}

	return genHandler(project, dir, mc)
}

func genHandler(project, dir string, mc *compilers.ModsCompiler) error {
	ms := handler.NewHandlersService(project, mc)
	to := path.Join(dir, "handler_gen.go")

	if err := filewriter.GenFile(to, ms); err != nil {
		return err
	}

	slog.Info("generate handlers", "file", filewriter.SprintGenPath(to))

	return nil
}

func genCodecHandler(project, dir string, mc *compilers.ModsCompiler) error {
	ms := handler.NewRespService(project, mc)
	to := path.Join(dir, "resp_gen.go")

	if err := filewriter.GenFile(to, ms); err != nil {
		return err
	}

	slog.Info("generate resp handler", "file", filewriter.SprintGenPath(to))

	return nil
}

func genModsHandler(project, dir string, mc *compilers.ModsCompiler, cs []*compilers.SeqCompiler) error {
	for _, c := range cs {
		if mc.Group != c.Group {
			continue
		}

		s := handler.NewModService(project, c.Mod(), c)
		to := path.Join(dir, fmt.Sprintf("%s_gen.go", c.Mod()))

		if err := filewriter.GenFile(to, s); err != nil {
			return err
		}

		slog.Info("generate mod handler", "file", filewriter.SprintGenPath(to))
	}

	return nil
}
