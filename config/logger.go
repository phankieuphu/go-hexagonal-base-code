package config

type Logger struct {
	Level     string // debug, info, warn, error
	Format    string // json, text
	AddSource bool
}

func loadLoggerConfig() Logger {
	return Logger{
		Level:     GetEnv("LOG_LEVEL", "info"),
		Format:    GetEnv("LOG_FORMAT", "json"),
		AddSource: getEnvBool("LOG_ADD_SOURCE", false),
	}
}
