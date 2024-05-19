{{ template "autogen_comment" }}
{{- $goName := .GoName }}

package di

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/infra/spanner/repository"
)

var {{ .GoName }}RepositorySet = fx.Provide(
{{- range .TableNames }}
	repository.New{{ . }}Repository,
{{- end }}
)
