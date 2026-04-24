package fiberswagger

type Config struct {
	// Base path for the SwaggerUI.
	//
	// Optional. Default: "/docs"
	BasePath string

	// OpenAPI specification file path to be rendered.
	//
	// Optional. Default: "./openapi.yaml"
	FilePath string
}

var defaultConfig = Config{
	BasePath: "/docs",
	FilePath: "./openapi.yaml",
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return defaultConfig
	}

	cfg := config[0]

	if cfg.BasePath == "" {
		cfg.BasePath = defaultConfig.BasePath
	}

	if cfg.FilePath == "" {
		cfg.FilePath = defaultConfig.FilePath
	}

	return cfg
}
