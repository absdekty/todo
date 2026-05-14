package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RestPort string
	AppMode  string
	DBPath   string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Shutdown     time.Duration
}

func Init() *Config {
	godotenv.Load(".env")
	godotenv.Overload(".local.env")

	return &Config{
		RestPort: getEnvString("REST_PORT", "8080"),
		AppMode:  getEnvString("APP_MODE", "prod"),
		DBPath:   getEnvString("DB_PATH", "./tasks.db"),

		ReadTimeout:  getEnvDuration("READ_TIMEOUT", time.Second*10),
		WriteTimeout: getEnvDuration("WRITE_TIMEOUT", time.Second*10),
		IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", time.Second*30),
		Shutdown:     getEnvDuration("SHUTDOWN_TIMEOUT", time.Second*30),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
