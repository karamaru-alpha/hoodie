package di

import (
	"net/http"

	"connectrpc.com/connect"
	"go.uber.org/fx"

	"github.com/karamaru-alpha/hoodie/pkg/cmd/api/handler/example"
	"github.com/karamaru-alpha/hoodie/pkg/pb/rpc/api/apiconnect"
)

var handlerSet = fx.Provide(
	example.New,
)

func invokeHandler(
	mux *http.ServeMux,
	opts []connect.HandlerOption,
	exampleHandler apiconnect.ExampleHandler,
) {
	mux.Handle(apiconnect.NewExampleHandler(exampleHandler, opts...))
}
