{{ template "autogen_comment" }}
package {{ .PkgName }}connect

import (
	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/domain/dto"
)

var Specs = dto.ConnectSpecs{
{{- range .Methods }}
	{Procedure: {{ .GoName }}Procedure, Description: "{{ .Description }}", IdempotencyLevel: connect.Idempotency{{ .IdempotencyLevel }} },
{{- end }}
}
