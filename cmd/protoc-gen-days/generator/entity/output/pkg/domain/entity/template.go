package entity

import (
	"bytes"
	_ "embed"
	"log/slog"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/entity/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

type Data struct {
	PkgName    string
	TableNames []string
}

type Column struct {
	GoName    string
	CamelName string
	Comment   string
	Type      string
	GoType    string
	IsList    bool
	IsEnum    bool

	PK bool
}

type Index struct {
	GoName  string
	Comment string
	Columns []*Column
}

type Table struct {
	PkgName   string
	GoName    string
	Comment   string
	Columns   []*Column
	PKColumns []*Column
}

//go:embed types.gen.go.tpl
var typesTemplateFileBytes []byte

func NewTypesCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &typesCreator{
		tpl: template.Must(tpl.New("pkg/domain/entity/types.gen.go.tpl").Parse(string(typesTemplateFileBytes))),
	}
}

type typesCreator struct {
	tpl *template.Template
}

func (c *typesCreator) Create(message *input.Message) (*core.TemplateInfo, error) {
	data := &Table{
		PkgName:   message.PkgName,
		GoName:    message.GoName,
		Comment:   message.Comment,
		Columns:   make([]*Column, 0, len(message.Fields)),
		PKColumns: make([]*Column, 0),
	}

	for _, field := range message.Fields {
		column := &Column{
			GoName:    field.GoName,
			CamelName: field.CamelName,
			Comment:   field.Comment,
			Type:      field.Type,
			GoType:    field.Type,
			PK:        field.PK,
		}
		if field.IsList {
			if field.IsEnum {
				column.GoType = field.Type + "Slice"
			} else {
				column.GoType = "[]" + field.Type
			}
		}
		data.Columns = append(data.Columns, column)
		if field.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	slog.Info(strings.Join([]string{message.FileDirName, message.SnakeName + ".gen.go"}, "/"))

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{message.FileDirName, message.SnakeName + ".gen.go"}, "/"),
	}, nil
}
