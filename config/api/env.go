package api

import (
	"os"

	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/config"
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

	return cfg, nil
}
