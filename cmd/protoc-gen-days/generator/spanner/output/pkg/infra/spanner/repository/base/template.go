package base

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/scylladb/go-set/strset"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

//go:embed buffer.gen.go.tpl
var bufferTemplateFileBytes []byte

type bufferCreator struct {
	tpl *template.Template
}

func NewBufferCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &bufferCreator{tpl: template.Must(tpl.New("pkg/infra/spanner/repository/base/buffer.gen.go.tpl").Parse(string(bufferTemplateFileBytes)))}
}

func (c *bufferCreator) Create(message *input.Message) (*core.TemplateInfo, error) {
	type Column struct {
		GoName string
		Type   string
		IsList bool
	}
	type Table struct {
		PkgName   string
		GoName    string
		CamelName string
		Columns   []*Column
		PKColumns []*Column
	}
	data := &Table{
		PkgName:   message.PkgName,
		GoName:    message.GoName,
		CamelName: message.CamelName,
		Columns:   make([]*Column, 0, len(message.Fields)),
		PKColumns: make([]*Column, 0),
	}

	for _, field := range message.Fields {
		column := &Column{
			GoName: field.GoName,
			Type:   field.GoType,
			IsList: field.IsList,
		}
		if field.IsList {
			if field.IsEnum {
				column.Type = field.GoType + "Slice"
			} else {
				column.Type = "[]" + field.GoType
			}
		}
		if field.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		data.Columns = append(data.Columns, column)
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner/repository/base", fmt.Sprintf("%s_buffer.gen.go", message.SnakeName)}, "/"),
	}, nil
}

//go:embed cache.gen.go.tpl
var cacheTemplateFileBytes []byte

type cacheCreator struct {
	tpl *template.Template
}

func NewCacheCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &cacheCreator{tpl: template.Must(tpl.New("pkg/infra/spanner/repository/base/cache.gen.go.tpl").Parse(string(cacheTemplateFileBytes)))}
}

func (c *cacheCreator) Create(message *input.Message) (*core.TemplateInfo, error) {
	type Column struct {
		GoName  string
		Type    string
		SetType string
		IsEnum  bool
	}
	type Table struct {
		PkgName   string
		GoName    string
		CamelName string
		Columns   []*Column
		PKColumns []*Column
	}
	data := &Table{
		PkgName:   message.PkgName,
		GoName:    message.GoName,
		CamelName: message.CamelName,
		Columns:   make([]*Column, 0, len(message.Fields)),
		PKColumns: make([]*Column, 0),
	}

	for _, field := range message.Fields {
		column := &Column{
			GoName:  field.GoName,
			Type:    field.GoType,
			SetType: field.SetType,
			IsEnum:  field.IsEnum,
		}

		if field.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		data.Columns = append(data.Columns, column)
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner/repository/base", fmt.Sprintf("%s_cache.gen.go", message.SnakeName)}, "/"),
	}, nil
}

//go:embed definition.gen.go.tpl
var definitionTemplateFileBytes []byte

type definitionCreator struct {
	tpl *template.Template
}

func NewDefinitionCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &definitionCreator{tpl: template.Must(tpl.New("pkg/infra/spanner/repository/base/definition.gen.go.tpl").Parse(string(definitionTemplateFileBytes)))}
}

func (c *definitionCreator) Create(message *input.Message) (*core.TemplateInfo, error) {
	type Column struct {
		GoName     string
		ColumnName string
	}
	type Table struct {
		TableName string
		GoName    string
		Columns   []*Column
	}
	data := &Table{
		TableName: message.GoName,
		GoName:    message.GoName,
		Columns:   make([]*Column, 0, len(message.Fields)),
	}

	for _, field := range message.Fields {
		column := &Column{
			GoName:     field.GoName,
			ColumnName: field.GoName,
		}
		data.Columns = append(data.Columns, column)
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner/repository/base", message.SnakeName + "_definition.gen.go"}, "/"),
	}, nil
}

//go:embed query_builder.gen.go.tpl
var queryBuilderTemplateFileBytes []byte

type queryBuilderCreator struct {
	tpl *template.Template
}

func NewQueryBuilderCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &queryBuilderCreator{tpl: template.Must(tpl.New("pkg/infra/spanner/repository/base/query_builder.gen.go.tpl").Parse(string(queryBuilderTemplateFileBytes)))}
}

func (c *queryBuilderCreator) Create(message *input.Message) (*core.TemplateInfo, error) {
	type Column struct {
		GoName  string
		Type    string
		SetType string
		IsList  bool
		IsEnum  bool
		PK      bool
	}
	type Type struct {
		GoName  string
		Key     string
		Columns []*Column
		IsIndex bool
	}
	type Table struct {
		GoName    string
		CamelName string
		Columns   []*Column
		Types     []*Type
	}
	data := &Table{
		GoName:    message.GoName,
		CamelName: message.CamelName,
		Columns:   make([]*Column, 0, len(message.Fields)),
		Types:     make([]*Type, 0),
	}

	for _, field := range message.Fields {
		column := &Column{
			GoName:  field.GoName,
			Type:    field.GoType,
			SetType: field.SetType,
			IsList:  field.IsList,
			IsEnum:  field.IsEnum,
			PK:      field.PK,
		}
		data.Columns = append(data.Columns, column)
	}
	data.Types = append(data.Types, &Type{
		GoName:  "All",
		Key:     "All",
		Columns: data.Columns,
		IsIndex: false,
	})

	for _, index := range message.Indexes {
		keySet := strset.NewWithSize(len(index.Keys) + len(index.PascalStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			goNames = append(goNames, key.GoName)
			keySet.Add(key.GoName)
		}
		keySet.Add(index.PascalStoring...)

		columns := make([]*Column, 0, len(data.Columns))
		for _, column := range data.Columns {
			if column.PK || keySet.Has(column.GoName) {
				columns = append(columns, column)
			}
		}
		data.Types = append(data.Types, &Type{
			GoName:  "Idx" + strings.Join(goNames, ""),
			Key:     "Idx" + data.GoName + "By" + strings.Join(goNames, ""),
			Columns: columns,
			IsIndex: true,
		})
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner/repository/base", fmt.Sprintf("%s_query_builder.gen.go", message.SnakeName)}, "/"),
	}, nil
}
