{{ template "autogen_comment" }}
{{- $pkColumns := .PKColumns }}
{{- $goName := .GoName }}
{{- $pkgName := .PkgName }}
package {{ .PkgName }}

import (
	"github.com/karamaru-alpha/days/pkg/domain/entity"
)

// {{ .GoName }} {{ .Comment }}
type {{ .GoName }} struct {
{{- range .Columns }}
	// {{ .Comment }}
	{{ .GoName }} {{ .GoType }} `json:"{{ .CamelName }},omitempty"`
{{- end }}
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

type {{ .GoName }}PKs []*{{ .GoName }}PK

func (e *{{ .GoName }}PK) ToEntity() *{{ .GoName }} {
	return &{{ .GoName }}{
	{{- range .PKColumns }}
		{{ .GoName }}: e.{{ .GoName }},
	{{- end }}
	}
}
