{{ template "autogen_comment" }}
{{ $goName := .GoName -}}
package base

const (
	{{ .GoName }}TableName = "{{ .TableName }}"
	{{ range .Columns -}}
	{{ $goName }}ColumnName{{ .GoName }} = "{{ .ColumnName }}"
	{{ end -}}
)

var (
	{{ .GoName }}ColumnNames = []string{
	{{ range .Columns -}}
		{{ $goName }}ColumnName{{ .GoName }},
	{{ end -}}
	}
)
