package api

import "github.com/karamaru-alpha/hoodie/pkg/domain/config"

var localAPIConfig = &config.APIConfig{
	Env:  config.EnvLocal,
	Port: "8080",
}
