package config

import "os"

type Config struct {
	Port      string
	GinMode   string
	LogLevel  string
	LogFormat string
}

func New() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		GinMode:   getEnv("GIN_MODE", "release"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "text"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
