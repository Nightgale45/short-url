package config

import (
	"os"
	"strconv"

	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/joho/godotenv"
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

func getEnvInt(key string, defaultVal int) int {
	envVal := os.Getenv(key)

	val, err := strconv.Atoi(envVal)
	if err != nil {
		logger.GetInstance().Error("CONFIG: Cannot conver env vale to int",
			"envKey", key)
		panic(err)
	}

	if val == 0 {
		logger.GetInstance().Warn("CONFIG: value is zero env key using default value",
			"envKey", key,
			"defaultVal", defaultVal)
		return defaultVal
	}

	return val
}
