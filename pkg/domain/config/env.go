package config

type Env string

const (
	EnvLocal Env = "local"
)

func (e Env) IsLocal() bool {
	return e == EnvLocal
}
