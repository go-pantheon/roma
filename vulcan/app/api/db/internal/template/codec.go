package template

import (
	"text/template"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/fabrica-util/errors"
)

const codecTemplate = `
package {{ .Package }}

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

{{- range .Messages }}

{{- $msgName := .Name | ToUpperCamel }}

{{ if .HasOneof }}
	{{- range .Fields }}
	  {{- $filedName := .Name | ToUpperCamel }}
		{{- if .IsOneof }}
			func Encode{{ $msgName }}(module life.Module) (*{{ $msgName }}, error) {		
				p := module.EncodeServer()
				mp := UserModuleProtoPool.Get()

				switch p.(type) {
				{{- range .OneofElements }}
				case *{{ .Type }}:
					mp.Module = &{{ $msgName }}_{{ .Name | ToUpperCamel }}{
						{{ .Name | ToUpperCamel }}: p.(*{{ .Type }}),
					}
					return mp, nil
				{{- end }}
				default:
					return nil, errors.Errorf("{{ $msgName }} encode invalid type: %T", module)
			  }				
		  }

		  func Decode{{ $msgName }}(p *{{ $msgName }}, module life.Module) error {
				if p.{{ .Name | ToUpperCamel }} == nil {
					return errors.New("{{ $msgName }}.{{ .Name | ToUpperCamel }} is nil")
				}
				
				switch p.{{ .Name | ToUpperCamel }}.(type) {
		    {{- range .OneofElements }}
				case *{{ $msgName }}_{{ .Name | ToUpperCamel }}:
				  return module.DecodeServer(p.Get{{ .Name | ToUpperCamel }}())
				{{- end }}
				default:
					return errors.Errorf("{{ $msgName }} decode invalid type: %T", p.{{ .Name | ToUpperCamel }})
			  }
		  }
		{{- end }}
	{{ end }}	 
{{- end }}

{{- end }}
`

type CodecService struct {
}

func NewCodecTemplate() (*template.Template, error) {
	tmpl, err := template.New("protoPool").Funcs(template.FuncMap{
		"ToLowerCamel": camelcase.ToLowerCamel,
		"ToUpperCamel": camelcase.ToUpperCamel,
	}).Parse(codecTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse template")
	}

	return tmpl, nil
}
