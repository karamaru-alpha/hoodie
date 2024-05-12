package api

import (
	"os"

	"github.com/karamaru-alpha/hoodie/pkg/domain/config"
	"github.com/karamaru-alpha/hoodie/pkg/herrors"
)

func New() (*config.APIConfig, error) {
	env := os.Getenv("ENV")
	var cfg *config.APIConfig
	switch config.Env(env) {
	case config.EnvLocal:
		cfg = localAPIConfig
	default:
		return nil, herrors.New(herrors.Internal, "environment variable 'ENV' is invalid").SetValues(map[string]any{"ENV": env})
	}

	return cfg, nil
}
