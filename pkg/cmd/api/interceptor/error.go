package interceptor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"

	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/config"
)

const stackTraceSize = 4 << 10 // 4KB

func NewErrorInterceptor(env config.Env) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					e, ok := r.(error)
					if !ok {
						e = fmt.Errorf("%+v", r)
					}
					stack := make([]byte, stackTraceSize)
					length := runtime.Stack(stack, true)
					qerr := derrors.Wrap(e, derrors.Internal, "panic occurred").SetValues(map[string]any{
						"panicMessage": fmt.Sprintf("%+v", r),
						"stackTrace":   string(stack[:length]),
					})
					err = errorHandler(ctx, qerr, env)
				}
			}()

			res, err = next(ctx, req)
			if err == nil {
				return res, nil
			}

			return nil, errorHandler(ctx, err, env)
		}
	}
}

func errorHandler(ctx context.Context, err error, env config.Env) error {
	qerr, _ := derrors.As(err)
	pattern := qerr.ErrorPattern()
	slog.ErrorContext(ctx, "", "err", err.Error())

	// hide error message if env isn't local
	if !env.IsLocal() {
		err = errors.New(pattern.ErrorCode.String())
	}
	e := connect.NewError(pattern.ConnectStatusCode, err)
	return e
}
