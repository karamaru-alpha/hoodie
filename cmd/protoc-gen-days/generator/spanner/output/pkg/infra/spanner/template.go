package spanner

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

//go:embed buffer_write.gen.go.tpl
var templateFileBytes []byte

type Table struct {
	GoName string
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.BulkTemplateCreator[*input.Message] {
	return &creator{tpl: template.Must(tpl.New("pkg/infra/spanner/buffer_write.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(_ string, messages []*input.Message) (*core.TemplateInfo, error) {
	data := make([]*Table, 0, len(messages))
	for _, message := range messages {
		data = append(data, &Table{
			GoName: message.GoName,
		})
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner", "buffer_write.gen.go"}, "/"),
	}, nil
}
