package config

type APIConfig struct {
	Env             Env
	Port            string
	ProjectID       string
	SpannerInstance string
	SpannerDB       string
}

// pkg/cmd/api/di/spanner.go#SpannerConfig

func (c *APIConfig) GetSpannerProjectID() string {
	return c.ProjectID
}

func (c *APIConfig) GetSpannerInstance() string {
	return c.SpannerInstance
}

func (c *APIConfig) GetSpannerDB() string {
	return c.SpannerDB
}
