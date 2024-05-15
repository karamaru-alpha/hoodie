package example

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/karamaru-alpha/days/pkg/pb/rpc/api"
	"github.com/karamaru-alpha/days/pkg/pb/rpc/api/apiconnect"
)

type handler struct{}

func New() apiconnect.ExampleHandler {
	return &handler{}
}

func (h *handler) Ping(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[api.ExamplePingResponse], error) {
	return connect.NewResponse(&api.ExamplePingResponse{
		Message: "hello, world!",
	}), nil
}
