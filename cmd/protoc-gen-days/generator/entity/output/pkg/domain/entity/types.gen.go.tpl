{{ template "autogen_comment" }}
{{- $pkColumns := .PKColumns }}
{{- $goName := .GoName }}
{{- $pkgName := .PkgName }}
package {{ .PkgName }}

import (
	"github.com/karamaru-alpha/days/pkg/domain/entity"
	"github.com/karamaru-alpha/days/pkg/domain/dto"
)

const (
	{{ .GoName }}TableName   = "{{ .GoName }}"
	{{ .GoName }}Comment     = "{{ .Comment }}"

	{{ range .Columns -}}
	{{ $goName }}ColumnName_{{ .GoName }} = "{{ .CamelName }}"
	{{ end -}}
)

// {{ .GoName }} {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .GoType }} `json:"{{ .CamelName }},omitempty"`
{{- end }}
}


func (e *{{ .GoName }}) GetPK() *{{ .GoName }}PK {
	return &{{ .GoName }}PK{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}) FullDeepCopy() *{{ .GoName }} {
	{{- range .Columns }}
	{{- if and (not .PK) (.IsList) }}
	{{ .CamelName }} := make([]{{ .Type }}, len(e.{{ .GoName }}))
	copy({{ .CamelName }}, e.{{ .GoName }})
	{{- end }}
	{{- end }}
	return &{{ .GoName }}{
		{{- range .Columns }}
		{{ if .IsList }}{{ .GoName }}: {{ .CamelName }},
		{{- else -}}{{ .GoName }}: e.{{ .GoName }},
		{{- end }}
		{{- end }}
	}
}

func (e *{{ .GoName }}) ToKeyValue() map[string]any {
	return map[string]any{
	{{- range .Columns }}
		{{ $goName }}ColumnName_{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

type {{ .GoName }}Slice []*{{ .GoName }}

type {{ .GoName }}MapByPK {{ range .PKColumns }}map[{{ .Type }}]{{ end }}*{{ .GoName }}

func (s {{ .GoName }}Slice) CreateMapByPK() {{ .GoName }}MapByPK {
	m := make({{ .GoName }}MapByPK, len(s))
	for _, row := range s {
		{{- range $i, $_ := slice .PKColumns 0 (sub (len .PKColumns) 1) }}
		{{- $cols := slice $pkColumns 0 (add1 $i) }}
		if _, ok := m{{ range $cols }}[row.{{ .GoName }}]{{ end }}; !ok {
			m{{ range slice $cols }}[row.{{ .GoName }}]{{ end }} = make({{ range slice $pkColumns (add1 $i) (len $pkColumns)}}map[{{ .Type }}]{{ end }}*{{ $goName }})
		}
		{{- end }}
		m{{ range .PKColumns }}[row.{{ .GoName }}]{{ end }} = row
	}
	return m
}

type {{ .GoName }}PK struct {
	{{ range .PKColumns -}}
		{{ .GoName }} {{ .GoType }}
	{{ end -}}
}

func (e *{{ .GoName }}PK) Key() string {
	return strings.Join([]string{
	{{- range .PKColumns }}
		fmt.Sprint(e.{{ .GoName }}),
	{{- end }}
	}, ".")
}

func (e *{{ .GoName }}PK) Generate() []any {
	return []any{
	{{- range .PKColumns }}
		e.{{ .GoName }},
	{{- end }}
	}
}

type {{ .GoName }}PKs []*{{ .GoName }}PK

func (e *{{ .GoName }}PK) ToEntity() *{{ .GoName }} {
	return &{{ .GoName }}{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

{{- range .Indexes }}

// {{ .GoName }} {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .GoType }}
{{- end }}
}

func (e *{{ .GoName }}) GetPK() *{{ $goName }}PK {
	return &{{ $goName }}PK{
	{{- range $pkColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

func (e *{{ .GoName }}) ToEntity() *{{ $goName }} {
	return &{{ $goName }}{
	{{- range .Columns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}

type {{ .GoName }}Slice []*{{ .GoName }}

func (s {{ .GoName }}Slice) GetPKs() {{ $goName }}PKs {
	pks := make({{ $goName }}PKs, 0, len(s))
	for _, row := range s {
		pks = append(pks, row.GetPK())
	}
	return pks
}
{{- end }}

var {{ .GoName }}ColumnMap = map[string]*dto.Column{
	{{- range .Columns }}
	{{ $goName }}ColumnName_{{ .GoName }}: {
		Name:     {{ $goName }}ColumnName_{{ .GoName }},
		Type:     "{{ .Type }}",
		IsList:   {{ .IsList }},
		PK:       {{ .PK }},
		Comment:  "{{ .Comment }}",
	},
	{{- end }}
}

var {{ .GoName }}Columns = dto.Columns{
	{{- range .Columns }}
	{{ $goName }}ColumnMap[{{ $goName }}ColumnName_{{ .GoName }}],
	{{- end }}
}
