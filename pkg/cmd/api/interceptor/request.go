package interceptor

import (
	"context"
	"log/slog"
	"net/netip"

	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/dcontext"
)

func NewRequestInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			param := &dcontext.SetRequestParam{}
			// IP
			if addrPort, err := netip.ParseAddrPort(req.Peer().Addr); err != nil {
				slog.ErrorContext(ctx, err.Error())
			} else {
				param.IP = addrPort.Addr().String()
			}
			dcontext.ExtractRequest(ctx).SetRequest(param)
			return next(ctx, req)
		}
	}
}
