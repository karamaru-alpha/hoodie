package enum

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/enum/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

//go:embed enum.gen.go.tpl
var templateFileBytes []byte

type Value struct {
	UpperSnakeName string
}

type Enum struct {
	PascalName string
	Values     []*Value
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.EachTemplateCreator[*input.Enum] {
	return &creator{tpl: template.Must(tpl.New("pkg/domain/enum/enum.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(enum *input.Enum) (*core.TemplateInfo, error) {
	data := &Enum{
		PascalName: core.ToPascalCase(enum.SnakeName),
		Values:     make([]*Value, 0, len(enum.Values)),
	}
	for _, value := range enum.Values {
		if value.RawName == "UNKNOWN" {
			continue
		}
		data.Values = append(data.Values, &Value{
			UpperSnakeName: core.ToUpperSnakeCase(value.RawName),
		})
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/domain/enum", enum.SnakeName + ".gen.go"}, "/"),
	}, nil
}
