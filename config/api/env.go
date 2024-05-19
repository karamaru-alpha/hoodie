package api

import (
	"log/slog"
	"os"

	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/config"
	"github.com/karamaru-alpha/days/pkg/util/envconfig"
)

func New() (*config.APIConfig, error) {
	env := os.Getenv("ENV")
	var cfg *config.APIConfig
	switch config.Env(env) {
	case config.EnvLocal:
		cfg = localAPIConfig
	default:
		return nil, derrors.New(derrors.Internal, "environment variable 'ENV' is invalid").SetValues(map[string]any{"ENV": env})
	}

	if err := envconfig.Load(cfg); err != nil {
		return nil, derrors.Wrap(err, derrors.Internal, "fail to load env")
	}

	slog.Info(cfg.Port)

	return cfg, nil
}
