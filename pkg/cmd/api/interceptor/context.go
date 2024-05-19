package interceptor

import (
	"context"

	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/dcontext"
)

func NewContextInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx = dcontext.SetTimeContext(ctx)
			ctx = dcontext.SetRequestInContext(ctx, dcontext.NewRequest())
			ctx = dcontext.SetUserInContext(ctx, dcontext.NewUser())
			return next(ctx, req)
		}
	}
}
