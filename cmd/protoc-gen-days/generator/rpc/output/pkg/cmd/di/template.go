package di

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/rpc/input"
)

//go:embed handler.gen.go.tpl
var templateFileBytes []byte

type Handler struct {
	GoName    string
	PkgName   string
	CamelName string
}

type Data struct {
	PkgName   string
	GoPkgName string
	Handlers  []*Handler
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.BulkTemplateCreator[*input.Service] {
	return &creator{tpl: template.Must(tpl.New("pkg/cmd/di/handler.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(pkgName string, messages []*input.Service) (*core.TemplateInfo, error) {
	data := &Data{
		PkgName:  pkgName,
		Handlers: make([]*Handler, 0, len(messages)),
	}

	var goPkgName string
	for _, message := range messages {
		goPkgName = message.GoPkgName
		data.Handlers = append(data.Handlers, &Handler{
			GoName:    message.GoName,
			PkgName:   core.ToSnakeCase(message.GoName),
			CamelName: message.CamelName,
		})
	}

	data.GoPkgName = goPkgName

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/cmd", pkgName, "di", "handler.gen.go"}, "/"),
	}, nil
}
