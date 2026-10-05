package config

import "account-service/pkg/logger"

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

// ToOptions converts the Logger config into pkg/logger.Options for logger.Init.
func (l Logger) ToOptions() logger.Options {
	format := logger.FormatJSON
	if logger.Format(l.Format) == logger.FormatText {
		format = logger.FormatText
	}

	return logger.Options{
		Level:     l.Level,
		Format:    format,
		AddSource: l.AddSource,
	}
}
