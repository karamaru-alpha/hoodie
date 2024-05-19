package di

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
)

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Data struct {
	GoName     string
	TableNames []string
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.BulkTemplateCreator[*input.Message] {
	return &creator{tpl: template.Must(tpl.New("pkg/di/repository.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(pkgName string, messages []*input.Message) (*core.TemplateInfo, error) {
	data := &Data{
		GoName:     core.ToGolangPascalCase(pkgName),
		TableNames: make([]string, 0, len(messages)),
	}

	for _, message := range messages {
		data.TableNames = append(data.TableNames, message.GoName)
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/di", pkgName + "_repository.gen.go"}, "/"),
	}, nil
}
