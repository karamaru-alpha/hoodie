package interceptor

import (
	"context"

	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/derrors"
)

type validator interface {
	Validate() error
}

func NewValidationInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if v, ok := req.Any().(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, derrors.Wrap(err, derrors.InvalidArgument, "validation error occurred")
				}
			}
			return next(ctx, req)
		}
	}
}
