{{ range . -}}
{{- $tableName := .TableName -}}
CREATE TABLE {{ .TableName }} (
  {{- range .Columns }}
  {{ .ColumnName }} {{ .Type }}
    {{- if not .Nullable }} NOT NULL{{ end -}}
    ,
  {{- end }}
) PRIMARY KEY (
  {{- range $i, $pk := .PKColumns -}}
    {{- if $i }}, {{ end }}{{ $pk.ColumnName }}{{ if $pk.Desc }} DESC{{ end }}
  {{- end -}}
)
{{- if .InterleaveTableName -}}
  ,
  INTERLEAVE IN PARENT {{ .InterleaveTableName }} ON DELETE CASCADE
{{- end -}}
{{- if .DeletionPolicy -}}
  ,
  ROW DELETION POLICY (OLDER_THAN({{ .DeletionPolicy.TimestampColumn }}, INTERVAL {{ .DeletionPolicy.Days }} DAY))
{{- end -}}
;
{{ range .Indexes -}}
CREATE {{ if .Unique }}UNIQUE {{ end }}{{ if .NullFiltered }}NULL_FILTERED {{ end }}INDEX Idx{{ $tableName }}By{{ range .Keys }}{{ .ColumnName }}{{ end }} ON {{ $tableName }}(
  {{- range $i, $col := .Keys }}{{ if $i }}, {{ end }}{{ $col.ColumnName }}{{ if $col.Desc }} DESC{{ end }}{{ end -}}
)
{{- if len .Storing }} STORING ({{ range $i, $s := .Storing }}{{ if $i }}, {{ end }}{{ $s }}{{ end }}){{ end -}}
;
{{ end }}
{{ end -}}
