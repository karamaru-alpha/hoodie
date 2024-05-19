package connect

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/rpc/input"
)

//go:embed specs.gen.go.tpl
var templateFileBytes []byte

type Method struct {
	GoName           string
	Description      string
	IdempotencyLevel string
}

type Data struct {
	PkgName string
	Methods []*Method
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.BulkTemplateCreator[*input.Service] {
	return &creator{tpl: template.Must(tpl.New("pkg/pb/rpc/connect/specs.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(pkgName string, messages []*input.Service) (*core.TemplateInfo, error) {
	data := &Data{
		PkgName: pkgName,
		Methods: make([]*Method, 0, len(messages)),
	}

	var goPkgName string
	for _, message := range messages {
		goPkgName = message.GoPkgName
		for _, method := range message.Methods {
			data.Methods = append(data.Methods, &Method{
				GoName:           method.GoName,
				Description:      method.Comment,
				IdempotencyLevel: method.IdempotencyLevel,
			})
		}
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, err
	}

	dir := core.GoPackageNameToFileDirName(goPkgName)

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{dir, fmt.Sprintf("%sconnect", pkgName), "specs.gen.go"}, "/"),
	}, nil
}
