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
	BaseUrl      string
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

var log = logger.GetInstance()

func LoadConf() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn("CONFIG: .env file not found or could not be loaded", "error", err)
	}
	return &Config{
		Env: getEnv("ENVIRONMENT", "development"),
		DatabaseConf: DatabaseConfig{
			Url:      getEnv("DATABASE_URL", "localhost"),
			MaxConns: getEnvInt("DATABASE_MAX_CONN", 5),
			MinConns: getEnvInt("DATABASE_MIN_CONN", 1),
		},
		RedisConf: RedisConfig{
			Addr:     getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", "redispassword"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		BaseUrl: getEnv("BASE_URL", "http://localhost"),
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

	if envVal == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(envVal)
	if err != nil {
		log.Warn("CONFIG: cannot convert env value to int, using default",
			"envKey", key,
			"defaultVal", defaultVal)
		return defaultVal
	}

	return val
}
