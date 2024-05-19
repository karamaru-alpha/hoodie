package repository

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

const (
	and               = "And"
	indexPrefix       = "Idx"
	indexMethodSuffix = "With"
)

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Column struct {
	GoName string
	Type   string
}

type Method struct {
	Name       string
	Args       string
	ReturnName string
}

type Table struct {
	PkgName string
	GoName  string
	Methods []*Method
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &creator{tpl: template.Must(tpl.New("pkg/domain/repository/repository.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(message *input.Message) (*core.TemplateInfo, error) {
	data := &Table{
		PkgName: message.PkgName,
		GoName:  message.GoName,
		Methods: make([]*Method, 0),
	}

	pkColumns := make([]*Column, 0)
	columnMap := make(map[string]*Column, len(message.Fields))
	for _, field := range message.Fields {
		column := &Column{
			GoName: field.GoName,
			Type:   field.GoType,
		}
		if field.PK {
			pkColumns = append(pkColumns, column)
		}
		columnMap[field.GoName] = column
	}

	var methodName string
	var methodArgs string
	last := len(pkColumns) - 1
	for i, column := range pkColumns {
		if i == last {
			break
		}
		methodName += column.GoName
		methodArgs += column.GoName
		data.Methods = append(data.Methods,
			&Method{
				Name:       methodName,
				Args:       methodArgs + " " + column.Type,
				ReturnName: data.GoName,
			},
			&Method{
				Name:       methodName + "s",
				Args:       methodArgs + "s" + " []" + column.Type,
				ReturnName: data.GoName,
			},
		)
		methodName += and
		methodArgs += " " + column.Type + ", "
	}

	for _, index := range message.Indexes {
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			if column, ok := columnMap[key.GoName]; ok {
				goNames = append(goNames, column.GoName)
			}
		}
		goName := indexPrefix + strings.Join(goNames, "")
		methodName = goName + indexMethodSuffix
		methodArgs = ""
		for _, key := range index.Keys {
			column, ok := columnMap[key.GoName]
			if !ok {
				continue
			}
			methodName += column.GoName
			methodArgs += column.GoName
			data.Methods = append(data.Methods,
				&Method{
					Name:       methodName,
					Args:       methodArgs + " " + column.Type,
					ReturnName: data.GoName + goName,
				},
				&Method{
					Name:       methodName + "s",
					Args:       methodArgs + "s" + " []" + column.Type,
					ReturnName: data.GoName + goName,
				},
			)
			methodName += and
			methodArgs += " " + column.Type + ", "
		}
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/domain/repository", message.PkgName, message.SnakeName + ".gen.go"}, "/"),
	}, nil
}
