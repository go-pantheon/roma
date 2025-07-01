package compilers

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/pkg/errors"
)

type SeqCompiler struct {
	filename string

	mod   ModType
	Group GroupType
	Apis  []*Api
}

type Api struct {
	UpperCamelName string

	Comment string
	CS      string
	SC      string
}

func NewSeqCompilers(filename string, group GroupType) (*SeqCompiler, error) {
	c := &SeqCompiler{
		filename: filename,
		Group:    group,
	}

	if err := c.Compile(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *SeqCompiler) Mod() ModType {
	return c.mod
}

func (c *SeqCompiler) WalkApis(f func(api *Api) bool) {
	for _, api := range c.Apis {
		if !f(api) {
			return
		}
	}
}

const (
	blankTag   = " "
	commentTag = "//"
	equalTag   = "="
	seqRegRule = `\s*=\s*[1-9]([0-9])*;`
)

func (c *SeqCompiler) Compile() error {
	f, err := os.Open(c.filename)
	if err != nil {
		return errors.Wrapf(err, "load seq file failed. %s", c.filename)
	}

	seqReg := regexp.MustCompile(seqRegRule)

	var seqComment string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if strings.Index(text, "enum") == 0 {
			if len(c.mod) > 0 {
				return errors.Wrapf(err, "multiple modules are defined in a file. %s %s %s", c.filename, c.mod, text)
			}

			c.mod = compileMod(text)

			continue
		}

		if len(c.mod) == 0 {
			continue
		}

		if strings.Index(text, commentTag) != 0 && seqReg.MatchString(text) {
			ignore, api, err := newApi(c.mod, seqComment, text)
			if err != nil {
				return err
			}

			if !ignore {
				c.Apis = append(c.Apis, api)
			}

			seqComment = ""

			continue
		}

		if strings.Index(text, commentTag) == 0 {
			seqComment += " " + strings.TrimSpace(strings.Replace(text, commentTag, "", 1))
		} else {
			seqComment = ""
		}
	}

	if err = scanner.Err(); err != nil {
		return errors.Wrapf(err, "read api file failed. %s", c.filename)
	}

	return nil
}

func newApi(mod ModType, comment string, seqText string) (ignore bool, api *Api, err error) {
	subs := strings.Split(seqText, equalTag)
	if len(subs) != 2 {
		return false, nil, errors.Errorf("no enum value defined. %s", seqText)
	}

	name := strings.TrimSpace(subs[0])

	if mod == ModType(camelcase.ToUnderScore("System")) && name == camelcase.ToUpperCamel("Handshake") {
		return true, nil, nil
	}

	api = &Api{
		UpperCamelName: camelcase.ToUpperCamel(name),
		Comment:        strings.TrimSpace(comment),
	}

	api.SC = "SC" + api.UpperCamelName

	if !strings.HasPrefix(comment, "@push") {
		api.CS = "CS" + api.UpperCamelName
	}

	return false, api, nil
}

func compileMod(text string) (mod ModType) {
	subs := strings.Split(text, blankTag)
	if len(subs) <= 1 {
		return ""
	}

	return ModType(strings.ToLower(camelcase.ToUnderScore(strings.ReplaceAll(strings.TrimSpace(subs[1]), "Seq", ""))))
}
