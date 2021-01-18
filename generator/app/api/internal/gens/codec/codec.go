package codec

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/vulcan/app/api/internal/template/codec"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/compilers"
	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/filewriter"
)

func Gen(project, base string, mcs []*compilers.ModsCompiler, scs []*compilers.SeqCompiler) error {
	dir := path.Join(base, "/codec")
	if err := os.Mkdir(dir, 0755); err != nil {
		return errors.Wrapf(err, "create codec dir failed. path: %s", dir)
	}

	if err := genCodec(project, dir, mcs); err != nil {
		return errors.Wrapf(err, "gen global codec failed")
	}

	for _, c := range scs {
		if err := genModsCodec(project, dir, c); err != nil {
			return errors.Wrapf(err, "gen mod codec failed. mod: %s", c.Mod())
		}
	}
	return nil
}

func genCodec(project, dir string, mcs []*compilers.ModsCompiler) error {
	s := codec.NewService(project, mcs)
	to := path.Join(dir, "codec_gen.go")
	if err := filewriter.GenFile(to, s); err != nil {
		return err
	}
	fmt.Println(to)
	return nil
}

func genModsCodec(project, dir string, c *compilers.SeqCompiler) error {
	s := codec.NewModService(project, c.Mod(), c)
	to := path.Join(dir, fmt.Sprintf("%s_gen.go", c.Mod()))
	if err := filewriter.GenFile(to, s); err != nil {
		return err
	}
	fmt.Println(to)
	return nil
}
