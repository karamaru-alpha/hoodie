{{ template "autogen_comment" }}
{{ $goName := .GoName -}}
{{ $camelName := .CamelName -}}
package base

import (
	"strconv"
	"strings"
	"time"

	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/i64set"
	"github.com/scylladb/go-set/strset"

	"github.com/karamaru-alpha/days/pkg/domain/enum"
)

const (
	{{- range .Types }}
	{{ $camelName }}Select{{ .GoName }}Query =
		"SELECT " +
		"`"{{ range $i, $col := .Columns }}{{ if $i }}"`, `"{{ end }} + {{ $goName }}ColumnName{{ $col.GoName }} + {{ end }} "` " +
		"FROM " + "`" + {{ $goName }}TableName + "`"
	{{- end }}
	{{ $camelName }}SelectCountQuery = "SELECT COUNT(*) FROM `" + {{ $goName }}TableName + "`"
)

type {{ .GoName }}QueryBuilder interface {
{{- range .Types }}
	Select{{ .GoName }}() {{ $goName }}QueryBuilderFirstClause
	{{- if .IsIndex }}
	SelectAllWith{{ .GoName }}() {{ $goName }}QueryBuilderFirstClause
	{{- end }}
{{- end }}
	SelectCount() {{ .GoName }}QueryBuilderFirstClause
}

type {{ .GoName }}QueryBuilderFinisher interface {
	OrderBy(orderPairs OrderPairs) {{ .GoName }}QueryBuilderFinisher
	Limit(limit int) {{ .GoName }}QueryBuilderFinisher
	Offset(offset int) {{ .GoName }}QueryBuilderFinisher
	GetQuery() (string, map[string]any)
	GetQueryConditions() []*{{ .GoName }}QueryCondition
}

type {{ .GoName }}QueryBuilderFirstClause interface {
	{{ .GoName }}QueryBuilderFinisher
	Where() {{ .GoName }}QueryBuilderPredicate
}

type {{ .GoName }}QueryBuilderSecondClause interface {
	{{ .GoName }}QueryBuilderFinisher
	And() {{ .GoName }}QueryBuilderPredicate
}

type {{ .GoName }}QueryBuilderPredicate interface {
	{{ range .Columns -}}
	{{ .GoName }}Eq(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Ne(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Gt(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Gte(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Lt(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Lte(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}IsNull() {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}IsNotNull() {{ $goName }}QueryBuilderSecondClause
	{{ if or (eq .Type "string") (eq .Type "byte") -}}
	{{ .GoName }}StartsWith(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ end -}}
	{{ if and (ne .SetType "" ) (not .IsList) -}}
	{{ if .IsEnum -}}
	{{ .GoName }}In(params {{ .Type }}Slice) {{ $goName }}QueryBuilderSecondClause
	{{ else -}}
	{{ .GoName }}In(params []{{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ .GoName }}Nin(params []{{ .Type }}) {{ $goName }}QueryBuilderSecondClause
	{{ end -}}
	{{ end -}}
	{{ end -}}
}

type {{ .GoName }}QueryCondition struct {
	column   string
	operator ConditionOperator
	value    any
}

type {{ .CamelName }}QueryBuilder struct {
	builder         *strings.Builder
	params          map[string]any
	paramIndex      int
	queryConditions []*{{ .GoName }}QueryCondition
}

func New{{ .GoName }}QueryBuilder() {{ .GoName }}QueryBuilder {
	return &{{ .CamelName }}QueryBuilder{
		builder: 		 &strings.Builder{},
		params: 		 make(map[string]any),
		paramIndex: 	 0,
		queryConditions: make( []*{{ .GoName }}QueryCondition, 0),
	}
}

func (qb *{{ .CamelName }}QueryBuilder) addParam(condition string, param any) {
	qb.paramIndex++
	paramKey := ParamBaseKey + strconv.Itoa(qb.paramIndex)
	qb.params[paramKey] = param
	qb.builder.WriteString(condition + "@" + paramKey)
}

{{ range .Types -}}
func (qb *{{ $camelName }}QueryBuilder) Select{{ .GoName }}() {{ $goName }}QueryBuilderFirstClause {
	qb.builder.WriteString({{ $camelName }}Select{{ .GoName }}Query{{- if .IsIndex }} + "@{FORCE_INDEX={{ .Key }}}"{{- end }})
	return qb
}
{{ if .IsIndex }}
func (qb *{{ $camelName }}QueryBuilder) SelectAllWith{{ .GoName }}() {{ $goName }}QueryBuilderFirstClause {
	qb.builder.WriteString({{ $camelName }}SelectAllQuery + "@{FORCE_INDEX={{ .Key }}}")
	return qb
}
{{- end }}
{{ end }}

func (qb *{{ .CamelName }}QueryBuilder) SelectCount() {{ .GoName }}QueryBuilderFirstClause {
	qb.builder.WriteString({{ $camelName }}SelectCountQuery)
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Where() {{ .GoName }}QueryBuilderPredicate {
	qb.builder.WriteString(" WHERE ")
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) And() {{ .GoName }}QueryBuilderPredicate {
	qb.builder.WriteString(" AND ")
	return qb
}

{{ range .Columns -}}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Eq(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.queryConditions = append(qb.queryConditions, &{{ $goName }}QueryCondition{column: {{ $goName }}ColumnName{{ .GoName }}, operator: ConditionOperatorEq, value: param})
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` = ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Ne(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` != ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Gt(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` >", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Gte(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` >= ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Lt(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` < ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Lte(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.addParam("`" + {{ $goName }}ColumnName{{ .GoName }} + "` <= ", param)
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}IsNull() {{ $goName }}QueryBuilderSecondClause {
	qb.builder.WriteString("`" + {{ $goName }}ColumnName{{ .GoName }} + "` IS NULL")
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}IsNotNull() {{ $goName }}QueryBuilderSecondClause {
	qb.builder.WriteString("`" + {{ $goName }}ColumnName{{ .GoName }} + "` IS NOT NULL")
	return qb
}

{{- if or (eq .Type "string") (eq .Type "byte") }}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}StartsWith(param {{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.builder.WriteString("STARTS_WITH(`" + {{ $goName }}ColumnName{{ .GoName }} + "`, ")
	qb.addParam("", param)
	qb.builder.WriteString(")")
	return qb
}
{{- end }}

{{ if and (ne .SetType "" ) (not .IsList) -}}
{{- if .IsEnum }}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}In(params {{ .Type }}Slice) {{ $goName }}QueryBuilderSecondClause {
	qb.queryConditions = append(qb.queryConditions, &{{ $goName }}QueryCondition{column: {{ $goName }}ColumnName{{ .GoName }}, operator: ConditionOperatorIn, value: params.Set()})
{{- else }}
func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}In(params []{{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	{{- if eq .Type "time.Time" }}
	v := i64set.New()
	for _, t := range params {
		v.Add(t.UnixNano())
	}
	qb.queryConditions = append(qb.queryConditions, &{{ $goName }}QueryCondition{column: {{ $goName }}ColumnName{{ .GoName }}, operator: ConditionOperatorIn, value: v})
	{{- else }}
	qb.queryConditions = append(qb.queryConditions, &{{ $goName }}QueryCondition{column: {{ $goName }}ColumnName{{ .GoName }}, operator: ConditionOperatorIn, value: {{ .SetType }}.New(params...)})
	{{- end }}
{{- end }}
	qb.builder.WriteString("`" + {{ $goName }}ColumnName{{ .GoName }} + "` IN (")
	for i, param := range params {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.addParam("", param)
	}
	qb.builder.WriteString(")")
	return qb
}

func (qb *{{ $camelName }}QueryBuilder) {{ .GoName }}Nin(params []{{ .Type }}) {{ $goName }}QueryBuilderSecondClause {
	qb.builder.WriteString("`" + {{ $goName }}ColumnName{{ .GoName }} + "` NOT IN (")
	for i, param := range params {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.addParam("", param)
	}
	qb.builder.WriteString(")")
	return qb
}
{{ end }}
{{ end }}

func (qb *{{ .CamelName }}QueryBuilder) OrderBy(orderPairs OrderPairs) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" ORDER BY ")
	for i, pair := range orderPairs {
		if i != 0 {
			qb.builder.WriteString(", ")
		}
		qb.builder.WriteString("`" + pair.Column + "` " + string(pair.OrderType))
	}
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Limit(limit int) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" LIMIT " + strconv.Itoa(limit))
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) Offset(offset int) {{ .GoName }}QueryBuilderFinisher {
	qb.builder.WriteString(" OFFSET " + strconv.Itoa(offset))
	return qb
}

func (qb *{{ .CamelName }}QueryBuilder) GetQuery() (string, map[string]any) {
	return qb.builder.String(), qb.params
}

func (qb *{{ .CamelName }}QueryBuilder) GetQueryConditions() []*{{ .GoName }}QueryCondition {
	return qb.queryConditions
}
