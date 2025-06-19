package template

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/go-pantheon/fabrica-util/errors"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(data *File, tmpl *template.Template) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, errors.Wrapf(err, "failed to execute template")
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to format generated code:\n%s", buf.String())
	}

	return formatted, nil
}
