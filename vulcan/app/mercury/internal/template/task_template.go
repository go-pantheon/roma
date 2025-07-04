package template

import (
	"bytes"
	"text/template"

	"github.com/go-pantheon/fabrica-util/camelcase"
	"github.com/go-pantheon/roma/vulcan/pkg/compilers"
)

var taskTemplate = `
{{- /* delete empty line */ -}}
// Code generated by gen-mercury. DO NOT EDIT.

package {{.UnderScoreMod}}

import (
	"reflect"

	climsg "{{.Project}}/gen/api/client/message"
	climod "{{.Project}}/gen/api/client/module"
	cliseq "{{.Project}}/gen/api/client/sequence"
	"{{.Project}}/mercury/internal/task"
)

var _ task.Taskable = (*{{.Api.UpperCamelName}}Task)(nil)

// {{.Api.UpperCamelName}}Task {{.Api.Comment}}
type {{.Api.UpperCamelName}}Task struct {
	*task.CommonTask
}

func New{{.Api.UpperCamelName}}Task(cs *climsg.{{.Api.CS}}, assert task.AssertFunc) *{{.Api.UpperCamelName}}Task {
	common := task.NewCommonTask(
		task.TypeCommon,
		climod.ModuleID_{{.UpperCamelMod}},
		int32(cliseq.{{.UpperCamelMod}}Seq_{{.Api.UpperCamelName}}),
		reflect.TypeOf(climsg.{{.Api.SC}}{}),
		cs,
		assert,
	)
	o := &{{.Api.UpperCamelName}}Task{
		CommonTask: common,
	}
	return o
}
`

type TaskService struct {
	Project string

	UpperCamelMod string
	UnderScoreMod string

	Api *compilers.Api
}

func NewTaskService(project string, c *compilers.SeqCompiler, api *compilers.Api) *TaskService {
	s := &TaskService{
		Project:       project,
		UpperCamelMod: camelcase.ToUpperCamel(string(c.Mod())),
		UnderScoreMod: camelcase.ToUnderScore(string(c.Mod())),
		Api:           api,
	}

	return s
}

func (s *TaskService) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("task").Parse(taskTemplate)
	if err != nil {
		return nil, err
	}

	if err = tmpl.Execute(buf, s); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
