package di

import (
	"connectrpc.com/connect"

	"github.com/karamaru-alpha/days/pkg/cmd/api/interceptor"
	"github.com/karamaru-alpha/days/pkg/domain/config"
)

func newHandlerOption(
	cfg *config.APIConfig,
) ([]connect.HandlerOption, error) {
	interceptors := []connect.Interceptor{
		interceptor.NewContextInterceptor(),
		interceptor.NewErrorInterceptor(cfg.Env),
		interceptor.NewRequestInterceptor(),
		interceptor.NewAuthInterceptor(),
		interceptor.NewValidationInterceptor(),
	}
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}, nil
}
