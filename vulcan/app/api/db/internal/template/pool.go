package template

import (
	"text/template"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/fabrica-util/errors"
)

const poolTemplate = `
package {{ .Package }}

import (
	"sync"

	{{ if .HasOneof }}
		"github.com/go-kratos/kratos/v2/log"
	{{ end }}
)

var (
	{{- range .Messages }}
		{{ .Name }}Pool = new{{ .Name }}Pool()
	{{- end }}
)

{{- range .Messages }}

{{- $msgName := .Name | ToUpperCamel }}

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
		{{- if .IsOneof }}
			if p.{{ .Name | ToUpperCamel }} != nil {
			  switch p.{{ .Name | ToUpperCamel }}.(type) {
		    {{- range .OneofElements }}
				case *{{ $msgName }}_{{ .Name | ToUpperCamel }}:
				  {{ .Type }}Pool.Put(p.Get{{ .Name | ToUpperCamel }}())
				{{- end }}
				default:
					log.Errorf("{{ $msgName }} put invalid type: %T", p.{{ .Name | ToUpperCamel }})
			  }
			}
		{{- else if .IsMap }}
			{{- if .MapValueIsMessage }}
				for _, v := range p.{{ .Name | ToUpperCamel }} {
					{{ .MapValueType }}Pool.Put(v)
				}
			{{ end }}
		{{- else if .IsRepeated }}
			{{- if .RepeatedValueIsMessage }}
				for _, v := range p.{{ .Name | ToUpperCamel }} {
					{{ .RepeatedValueType }}Pool.Put(v)
				}
			{{ end }}
		{{- else if .IsMessage }}
			{{- if .IsMessage }}
				{{ .Type }}Pool.Put(p.{{ .Name | ToUpperCamel }})
			{{- end }}
		{{ end }}
	{{ end }}
	p.Reset()
	pool.Pool.Put(p)
}

{{- end }}
`

func NewPoolTemplate() (*template.Template, error) {
	tmpl, err := template.New("protoPool").Funcs(template.FuncMap{
		"ToLowerCamel": camelcase.ToLowerCamel,
		"ToUpperCamel": camelcase.ToUpperCamel,
	}).Parse(poolTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template")
	}

	return tmpl, nil
}
