package di

import (
	"context"
	"net/http"

	"go.uber.org/fx"

	api_config "github.com/karamaru-alpha/days/config/api"
	"github.com/karamaru-alpha/days/pkg/cmd/api"
	"github.com/karamaru-alpha/days/pkg/di"
	"github.com/karamaru-alpha/days/pkg/domain/config"
)

func Initialize() fx.Option {
	return fx.Options(
		// DI
		basicOption(),
		// External DI (mocked in testing)
		externalOption(),
		// Hooks
		hooks(),
	)
}

func basicOption() fx.Option {
	return fx.Options(
		fx.Provide(
			context.Background,
			// alphabetically
			api.NewServer,
			di.NewDBTransactionTxManager,
			newHandlerOption,
			newServeMux,
			func(cfg *config.APIConfig) di.SpannerConfig {
				return cfg
			},
		),
		// alphabetically
		di.TransactionRepositorySet,
		handlerSet,
		serviceSet,
		usecaseSet,
	)
}

func externalOption() fx.Option {
	return fx.Options(
		fx.Provide(
			api_config.New,
			di.NewTransactionClient,
		),
	)
}

func hooks() fx.Option {
	return fx.Options(
		fx.Invoke(invokeHandler),
		fx.Invoke(func(lc fx.Lifecycle, s *http.Server) {
			lc.Append(fx.StartStopHook(api.Serve(s)))
		}),
	)
}
