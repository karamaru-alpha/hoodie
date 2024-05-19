package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/karamaru-alpha/days/pkg/cmd/api/usecase/user"
	"github.com/karamaru-alpha/days/pkg/dcontext"
	"github.com/karamaru-alpha/days/pkg/pb/rpc/api"
	"github.com/karamaru-alpha/days/pkg/pb/rpc/api/apiconnect"
)

type handler struct {
	userUsecase user.Usecase
}

func NewHandler(userUsecase user.Usecase) apiconnect.UserHandler {
	return &handler{
		userUsecase: userUsecase,
	}
}

func (h *handler) UpdateName(ctx context.Context, req *connect.Request[api.UserUpdateNameRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := h.userUsecase.UpdateName(ctx, dcontext.ExtractUser(ctx).GetUserID(), req.Msg.GetName()); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}
