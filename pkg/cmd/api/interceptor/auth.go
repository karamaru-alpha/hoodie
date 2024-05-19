package interceptor

import (
	"context"

	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/dcontext"
)

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// TODO: implement me
			dcontext.ExtractUser(ctx).SetUser(&dcontext.SetUserParam{
				UserID: "dummyUserID",
			})
			return next(ctx, req)
		}
	}
}
