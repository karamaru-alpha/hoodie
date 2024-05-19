package repository

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/scylladb/go-set/strset"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/core"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/generator/spanner/input"
	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

const (
	all               = "All"
	and               = "And"
	eq                = "Eq"
	in                = "In"
	indexPrefix       = "Idx"
	indexMethodSuffix = "With"
)

//go:embed repository.gen.go.tpl
var templateFileBytes []byte

type Column struct {
	GoName    string
	CamelName string
	LocalName string
	Type      string
	IsList    bool
	IsEnum    bool
	PK        bool
}

type Method struct {
	Name         string
	Args         string
	SliceArgName string
	ReturnName   string
	SelectType   string
	Wheres       string
	UseCache     bool
}

type Type struct {
	GoName  string
	Key     string
	Columns []*Column
}

type Table struct {
	PkgName            string
	GoName             string
	TableName          string
	CamelName          string
	PKColumns          []*Column
	Methods            []*Method
	Types              []*Type
	NeedCommonResponse bool
}

type creator struct {
	tpl *template.Template
}

func NewCreator(tpl *template.Template) core.EachTemplateCreator[*input.Message] {
	return &creator{tpl: template.Must(tpl.New("pkg/infra/spanner/repository/repository.gen.go.tpl").Parse(string(templateFileBytes)))}
}

func (c *creator) Create(message *input.Message) (*core.TemplateInfo, error) {
	data := &Table{
		PkgName:            message.PkgName,
		GoName:             message.GoName,
		TableName:          message.SnakeName,
		CamelName:          message.CamelName,
		PKColumns:          make([]*Column, 0),
		Methods:            make([]*Method, 0),
		Types:              make([]*Type, 0),
		NeedCommonResponse: message.NeedCommonResponse,
	}

	columns := make([]*Column, 0, len(message.Fields))
	columnMap := make(map[string]*Column)
	for _, field := range message.Fields {
		column := &Column{
			GoName:    field.GoName,
			CamelName: field.CamelName,
			LocalName: core.ToLocalName(field.CamelName),
			Type:      field.GoType,
			IsEnum:    field.IsEnum,
			IsList:    field.IsList,
			PK:        field.PK,
		}
		if field.PK {
			data.PKColumns = append(data.PKColumns, column)
		}
		columns = append(columns, column)
		columnMap[field.GoName] = column
	}
	data.Types = append(data.Types, &Type{
		GoName:  data.GoName,
		Key:     all,
		Columns: columns,
	})

	var methodName string
	var methodArgs string
	var methodWheres string
	last := len(data.PKColumns) - 1
	for i, column := range data.PKColumns {
		if i == last {
			break
		}
		argName := column.CamelName
		methodName += column.GoName
		methodArgs += argName
		methodWheres += column.GoName
		data.Methods = append(data.Methods,
			&Method{
				Name:       methodName,
				Args:       methodArgs + " " + column.Type,
				ReturnName: data.GoName,
				SelectType: all,
				Wheres:     methodWheres + eq + "(" + argName + ")",
				UseCache:   true,
			},
			&Method{
				Name:         methodName + "s",
				Args:         methodArgs + "s" + " []" + column.Type,
				SliceArgName: argName + "s",
				ReturnName:   data.GoName,
				SelectType:   all,
				Wheres:       methodWheres + in + "(" + argName + "s)",
				UseCache:     true,
			},
		)
		methodName += and
		methodArgs += " " + column.Type + ", "
		methodWheres += eq + "(" + argName + ")." + and + "()."
	}

	for _, index := range message.Indexes {
		typ := &Type{}
		keySet := strset.NewWithSize(len(index.Keys) + len(index.PascalStoring))
		goNames := make([]string, 0, len(index.Keys))
		for _, key := range index.Keys {
			if column, ok := columnMap[key.GoName]; ok {
				goNames = append(goNames, column.GoName)
				keySet.Add(column.GoName)
			}
		}
		typ.Key = indexPrefix + strings.Join(goNames, "")
		typ.GoName = data.GoName + typ.Key
		keySet.Add(index.PascalStoring...)

		for _, column := range columns {
			if column.PK || keySet.Has(column.GoName) {
				typ.Columns = append(typ.Columns, column)
			}
		}
		data.Types = append(data.Types, typ)

		methodName = typ.Key + indexMethodSuffix
		methodArgs = ""
		methodWheres = ""
		for i, key := range index.Keys {
			column, ok := columnMap[key.GoName]
			if !ok {
				continue
			}
			methodName += column.GoName
			methodArgs += column.LocalName
			methodWheres += column.GoName

			isNotNull := ""
			if index.NullFiltered {
				for _, k := range index.Keys[i+1:] {
					if col, ok := columnMap[k.GoName]; ok {
						isNotNull += "." + and + "()." + col.GoName + "IsNotNull()"
					}
				}
			}

			data.Methods = append(data.Methods,
				&Method{
					Name:       methodName,
					Args:       methodArgs + " " + column.Type,
					ReturnName: data.GoName + typ.Key,
					SelectType: typ.Key,
					Wheres:     methodWheres + eq + "(" + column.LocalName + ")" + isNotNull,
					UseCache:   false,
				},
				&Method{
					Name:       methodName + "s",
					Args:       methodArgs + "s []" + column.Type,
					ReturnName: data.GoName + typ.Key,
					SelectType: typ.Key,
					Wheres:     methodWheres + in + "(" + column.LocalName + "s)" + isNotNull,
					UseCache:   false,
				},
			)
			methodName += and
			methodArgs += " " + column.Type + ", "
			methodWheres += eq + "(" + column.LocalName + ")." + and + "()."
		}
	}

	buf := &bytes.Buffer{}
	if err := c.tpl.Execute(buf, data); err != nil {
		return nil, perrors.New(err.Error())
	}

	return &core.TemplateInfo{
		Data:     buf.Bytes(),
		FilePath: strings.Join([]string{"pkg/infra/spanner/repository", message.SnakeName + ".gen.go"}, "/"),
	}, nil
}
