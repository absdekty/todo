package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	RestPort string
	AppMode  string
	DBPath   string
}

func Init() *Config {
	godotenv.Load(".env")
	godotenv.Overload(".local.env")

	return &Config{
		RestPort: getEnvString("REST_PORT", "8080"),
		AppMode:  getEnvString("APP_MODE", "prod"),
		DBPath:   getEnvString("DB_PATH", "./tasks.db"),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
