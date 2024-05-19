package ddl

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

//go:embed init.gen.sql.tpl
var templateFileBytes []byte

type Column struct {
	ColumnName string
	Type       string
	Nullable   bool
	Desc       bool
}

type IndexKey struct {
	ColumnName string
	Desc       bool
}

type Index struct {
	Keys         []*IndexKey
	Unique       bool
	NullFiltered bool
	Storing      []string
}

type DeletionPolicy struct {
	TimestampColumn string
	Days            int32
}

type Table struct {
	TableName           string
	Columns             []*Column
	PKColumns           []*Column
	Indexes             []*Index
	InterleaveTableName string
	DeletionPolicy      *DeletionPolicy
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.BulkTemplateCreator[*input.Message] {
	return &creator{tpl: template.Must(tpl.New("db/ddl/init.gen.sql.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(pkgName string, messages []*input.Message) (*core.TemplateInfo, error) {
	data := make([]*Table, 0, len(messages))
	for _, message := range messages {
		table := &Table{
			TableName: message.GoName,
			Columns:   make([]*Column, 0, len(message.Fields)),
			PKColumns: make([]*Column, 0),
			Indexes:   make([]*Index, 0, len(message.Indexes)),
		}
		if message.Interleave != nil {
			table.InterleaveTableName = message.Interleave.GoName
		}
		if message.TTL != nil {
			table.DeletionPolicy = &DeletionPolicy{
				TimestampColumn: message.TTL.TimestampColumnGoName,
				Days:            message.TTL.Days,
			}
		}

		indexColumnMap := make(map[string]*Column)
		for _, field := range message.Fields {
			typ := field.DBType
			if field.IsList {
				typ = "ARRAY<" + typ + ">"
			}
			column := &Column{
				ColumnName: field.GoName,
				Type:       typ,
				Nullable:   !field.PK,
				Desc:       field.Desc,
			}

			table.Columns = append(table.Columns, column)
			if field.PK {
				table.PKColumns = append(table.PKColumns, column)
			}

			indexColumnMap[field.GoName] = column
		}

		for _, index := range message.Indexes {
			i := &Index{
				Keys:         make([]*IndexKey, 0, len(index.Keys)),
				Unique:       index.Unique,
				NullFiltered: index.NullFiltered,
				Storing:      index.PascalStoring,
			}
			for _, key := range index.Keys {
				for columnName, column := range indexColumnMap {
					if columnName == key.GoName {
						i.Keys = append(i.Keys, &IndexKey{
							ColumnName: column.ColumnName,
							Desc:       key.Desc,
						})
					}
				}
			}
			table.Indexes = append(table.Indexes, i)
		}

		if message.Interleave != nil {
			table.InterleaveTableName = message.Interleave.GoName
		}
		if message.TTL != nil {
			table.DeletionPolicy = &DeletionPolicy{
				TimestampColumn: message.TTL.TimestampColumnGoName,
				Days:            message.TTL.Days,
			}
		}

		data = append(data, table)
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"db/ddl", fmt.Sprintf("%s.gen.sql", pkgName)}, "/"),
	}, nil
}
