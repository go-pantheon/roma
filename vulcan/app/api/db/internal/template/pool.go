package template

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/fabrica-util/errors"
)

const poolTemplate = `
package {{ .Package }}

import (
	"sync"
)

var (
	{{- range .Messages }}
		{{ .Name }}Pool = new{{ .Name }}Pool()
	{{- end }}
)

{{- range .Messages }}

type {{ .Name | ToLowerCamel }}Pool struct {
	sync.Pool
}

func new{{ .Name }}Pool() *{{ .Name | ToLowerCamel }}Pool {
	return &{{ .Name | ToLowerCamel }}Pool{
		Pool: sync.Pool{
			New: func() any {
				return &{{ .Name }}{}
			},
		},
	}
}

func (pool *{{ .Name | ToLowerCamel }}Pool) Get() *{{ .Name }} {
	return pool.Pool.Get().(*{{ .Name }})
}

func (pool *{{ .Name | ToLowerCamel }}Pool) Put(p *{{ .Name }}) {
	{{- range .Fields }}
		{{- if .IsMap }}
			{{- if .ValueIsMessage }}
				for _, v := range p.{{ .Name | ToUpperCamel }} {
					{{ .ValueType }}Pool.Put(v)
				}
			{{- end }}
		{{- else if .IsRepeated }}
			{{- if .ValueIsMessage }}
				for _, v := range p.{{ .Name | ToUpperCamel }} {
					{{ .ValueType }}Pool.Put(v)
				}
			{{- end }}
		{{- else if .IsMessage }}
			{{- if .IsMessage }}
				{{ .Type }}Pool.Put(p.{{ .Name | ToUpperCamel }})
			{{- end }}
		{{- end }}
	{{- end }}

	p.Reset()
	pool.Pool.Put(p)
}
{{- end }}
`

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(data *File) ([]byte, error) {
	tmpl, err := template.New("protoPool").Funcs(template.FuncMap{
		"ToLowerCamel": camelcase.ToLowerCamel,
		"ToUpperCamel": camelcase.ToUpperCamel,
	}).Parse(poolTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template")
	}

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
