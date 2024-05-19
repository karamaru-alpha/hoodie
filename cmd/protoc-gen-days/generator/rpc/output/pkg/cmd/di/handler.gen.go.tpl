{{ template "autogen_comment" }}
{{- $pkgName := .PkgName }}

package di

import (
	"net/http"

	"connectrpc.com/connect"
	"go.uber.org/fx"

	{{ range .Handlers -}}
	"github.com/karamaru-alpha/days/pkg/cmd/{{ $pkgName }}/handler/{{ .PkgName }}"
	{{ end -}}
	"{{ .GoPkgName }}/{{ .PkgName }}connect"
)

var handlerSet = fx.Provide(
	{{ range .Handlers -}}
	{{ .PkgName }}.NewHandler,
	{{ end -}}
)

func invokeHandler(
	mux *http.ServeMux,
	opts []connect.HandlerOption,
	{{ range .Handlers -}}
	{{ .CamelName }}Handler {{ $pkgName }}connect.{{ .GoName }}Handler,
	{{ end -}}
) {
	{{ range .Handlers -}}
	mux.Handle({{ $pkgName }}connect.New{{ .GoName }}Handler({{ .CamelName }}Handler, opts...))
	{{ end -}}
}
