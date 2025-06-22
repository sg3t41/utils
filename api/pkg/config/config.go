package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	GinMode       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	UseMemoryDB   bool
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		GinMode:       getEnv("GIN_MODE", "debug"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "utils_user"),
		DBPassword:    getEnv("DB_PASSWORD", "utils_password"),
		DBName:        getEnv("DB_NAME", "utils_db"),
		UseMemoryDB:   getEnv("USE_MEMORY_DB", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}