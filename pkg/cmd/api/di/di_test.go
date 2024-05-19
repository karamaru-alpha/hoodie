package di

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/domain/config"
)

func TestDI_Initialize(t *testing.T) {
	t.Parallel()

	assert.NoError(t, fx.ValidateApp(fx.Options(
		basicOption(),
		fx.Provide(
			func() *config.APIConfig {
				return &config.APIConfig{
					Env: config.EnvLocal,
				}
			},
			func() *spanner.Client { return nil },
		),
		hooks(),
	)))
}
