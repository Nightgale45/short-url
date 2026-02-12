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
	godotenv.Load()
	return &Config{
		Env: getEnv("ENVIRONMENT", "development"),
		DatabaseConf: DatabaseConfig{
			Url:      getEnv("DATABASE_URL", "localhost"),
			MaxConns: getEnvInt("DATABASE_MAX_CONN", 5),
			MinConns: getEnvInt("DATABASE_MIN_CONN", 1),
		},
		RedisConf: RedisConfig{
			Addr:     getEnv("REDIS_URL", "locahost:6379"),
			Password: getEnv("REDIS_PASSWORD", "redispassword"),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		BaseUrl: getEnv("BASE_URL", "http://localhost:8080"),
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
		log.Error("CONFIG: Cannot conver env vale to int",
			"envKey", key)
		panic(err)
	}

	if val == 0 {
		log.Warn("CONFIG: value is zero env key using default value",
			"envKey", key,
			"defaultVal", defaultVal)
		return defaultVal
	}

	return val
}
