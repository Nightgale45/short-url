package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Env          string
	DatabaseConf DatabaseConfig
	RedisConf    RedisConfig
}

type DatabaseConfig struct {
	Url      string
	MaxConns int
	MinConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadConf() *Config {
	godotenv.Load()
	return &Config{
		Env: getEnv("ENVIRONMENT", "development"),
		DatabaseConf: DatabaseConfig{
			Url:      getEnv("DATABASE_URL", "localhost"),
			MaxConns: 5,
			MinConns: 1,
		},
		RedisConf: RedisConfig{
			Addr:     getEnv("REDIS_URL", "locahost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
	}
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
