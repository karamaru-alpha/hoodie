package entity

import (
	"bytes"
	_ "embed"
	"log/slog"
	"strings"
	"text/template"

	"github.com/scylladb/go-set/strset"

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
	Indexes   []*Index
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
		Indexes:   make([]*Index, 0),
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
	for _, index := range message.Indexes {
		i := &Index{
			GoName:  data.GoName + "Idx",
			Comment: data.GoName + " Index(",
			Columns: make([]*Column, 0, len(data.PKColumns)+len(index.Keys)+len(index.PascalStoring)),
		}
		keySet := strset.NewWithSize(len(index.Keys) + len(index.PascalStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			goNames = append(goNames, key.GoName)
			keySet.Add(key.GoName)
		}
		i.GoName += strings.Join(goNames, "")
		i.Comment += strings.Join(goNames, ", ") + ")"
		keySet.Add(index.PascalStoring...)

		for _, column := range data.Columns {
			if column.PK || keySet.Has(column.GoName) {
				i.Columns = append(i.Columns, column)
			}
		}
		data.Indexes = append(data.Indexes, i)
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
