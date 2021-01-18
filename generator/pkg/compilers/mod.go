package compilers

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"regexp"
	"sort"
	"strings"
)

type ModsCompiler struct {
	filename string

	Group GroupType
	Mods  []ModType
}

func NewModCompilers(filename string) ([]*ModsCompiler, error) {
	c := &ModsCompiler{filename: filename}
	if err := c.Compile(); err != nil {
		return nil, err
	}

	// 特殊处理
	mcs := make([]*ModsCompiler, 0, len(groupModMap)+1)
	playerCompiler := &ModsCompiler{filename: filename, Group: PlayerGroup}
	mcs = append(mcs, playerCompiler)

	for g := range groupModMap {
		mc := &ModsCompiler{filename: filename, Group: g}
		mcs = append(mcs, mc)
	}

	for _, mc := range mcs {
		for _, mod := range c.Mods {
			if GroupByMod(mod) == mc.Group {
				mc.Mods = append(mc.Mods, mod)
			}
		}
	}

	return mcs, nil
}

func (c *ModsCompiler) Compile() error {
	f, err := os.Open(c.filename)
	if err != nil {
		return errors.Wrapf(err, "加载mod文件错误。文件<%s>", c.filename)
	}

	modReg, err := regexp.Compile(`\s*=\s*[1-9]([0-9])*;`)
	if err != nil {
		return errors.Wrapf(err, "正则错误。文件<%s>", c.filename)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)

		if !modReg.MatchString(text) {
			continue
		}
		subs := strings.Split(text, equalTag)
		c.Mods = append(c.Mods, ModType(strings.ToLower(strings.TrimSpace(subs[0]))))
	}
	sort.Sort(ModSlice(c.Mods))

	if err = scanner.Err(); err != nil {
		return errors.Wrapf(err, "读取api文件错误。文件<%s>", c.filename)
	}
	return nil
}

func (c *ModsCompiler) Filename() string {
	return c.filename
}
