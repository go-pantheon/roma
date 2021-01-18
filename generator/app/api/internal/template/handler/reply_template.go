package handler

import (
	"bytes"
	"text/template"

	"github.com/vulcan-frame/vulcan-game/vulcan/pkg/compilers"
	"github.com/vulcan-frame/vulcan-util/camelcase"
)

var replyTemplate = `
{{- /* delete empty line */ -}}
// Code generated by gen-api. DO NOT EDIT.

package handler

import (
	"google.golang.org/protobuf/proto"
	"{{.Project}}/game/gen/api/server/{{.LowerGroup}}/intra/v1"
	"github.com/pkg/errors"
)

func New{{.UpperGroup}}Response(mod, seq int32, Obj int64, in proto.Message) (ret []byte, err error) {
	data, err := proto.Marshal(in)
	if err != nil {
		err = errors.Wrapf(err, "proto marshal failed. mod=%d seq=%d obj=%d", mod, seq, Obj)
		return
	}
	return New{{.UpperGroup}}ResponseByData(mod, seq, Obj, data)
}

func New{{.UpperGroup}}ResponseByData(mod, seq int32, Obj int64, data []byte) (ret []byte, err error) {
	p := New{{.UpperGroup}}ResponseProtoByData(mod, seq, Obj, data)
	ret, err = proto.Marshal(p)
	if err != nil {
		err = errors.Wrapf(err, "proto marshal failed. mod=%d seq=%d obj=%d", mod, seq, Obj)
	}
	return
}

func New{{.UpperGroup}}ResponseProto(mod, seq int32, Obj int64, in proto.Message) (ret *intrav1.TunnelResponse, err error) {
	data, err := proto.Marshal(in)
	if err != nil {
		err = errors.Wrapf(err, "proto marshal failed. mod=%d seq=%d obj=%d", mod, seq, Obj)
		return
	}
	ret = New{{.UpperGroup}}ResponseProtoByData(mod, seq, Obj, data)
	return 
}

func New{{.UpperGroup}}ResponseProtoByData(mod, seq int32, Obj int64, data []byte) (p *intrav1.TunnelResponse) {
	p = &intrav1.TunnelResponse{
		Mod:  mod,
		Seq:  seq,
		Obj:  Obj,
		Data: data,
	}
	return
}

`

type ReplyService struct {
	Project    string
	UpperGroup string
	LowerGroup string
}

func NewReplyService(project string, c *compilers.ModsCompiler) *ReplyService {
	s := &ReplyService{
		Project:    project,
		UpperGroup: camelcase.ToUpperCamel(string(c.Group)),
		LowerGroup: camelcase.ToLowerCamel(string(c.Group)),
	}
	return s
}

func (s *ReplyService) Execute() ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("handler_reply").Parse(replyTemplate)
	if err != nil {
		return nil, err
	}
	if err = tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
