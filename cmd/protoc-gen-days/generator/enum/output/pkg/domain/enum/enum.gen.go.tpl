{{ template "autogen_comment" }}
package enum

import (
	"strconv"
)

{{ $Name := .PascalName -}}
type {{ .PascalName }}Slice []{{ .PascalName }}

var {{ .PascalName }}Values = {{ .PascalName }}Slice{
    {{- range .Values }}
    {{ $Name }}_{{ .UpperSnakeName }},
    {{- end }}
}

func (e {{ .PascalName }}) EncodeSpanner() (any, error) {
	return int64(e), nil
}

func (e {{ .PascalName }}Slice) EncodeSpanner() (any, error) {
	slice := make([]int64, 0, len(e))
	for _, v := range e {
		slice = append(slice, int64(v))
	}
	return slice, nil
}
