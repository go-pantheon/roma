package main

import (
	"fmt"
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
		fmt.Printf("准备生成的 API 模块: %+v\n", mc.Mods)
	}

	if err := gen(baseDir, scs); err != nil {
		panic(err)
	}
}

func gen(base string, scs []*compilers.SeqCompiler) error {
	// 处理临时文件夹
	if _, err := os.Stat(destTmpPath); err == nil {
		if err = os.RemoveAll(destTmpPath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrapf(err, "请手动删除文件夹：%s", destTmpPath)
	}
	if err := os.Mkdir(destTmpPath, 0755); err != nil {
		return errors.Wrapf(err, "创建临时文件夹失败。路径：%s", destTmpPath)
	}

	// 生成到临时文件夹
	if err := genTask(destTmpPath, scs); err != nil {
		return err
	}

	if err := os.RemoveAll(destPath); err != nil {
		return err
	}
	if err := os.Rename(destTmpPath, destPath); err != nil {
		return err
	}
	fmt.Println("生成api文件完成。")
	return nil
}

func genTask(base string, cs []*compilers.SeqCompiler) error {
	for _, c := range cs {
		dir := base + "/" + string(c.Mod())
		if err := os.Mkdir(dir, 0755); err != nil {
			return errors.Wrapf(err, "创建临时文件夹失败。路径：%s", dir)
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
			fmt.Println(to)
		}
	}
	return nil
}
