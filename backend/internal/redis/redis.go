package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisService interface {
	SaveUrlMapping(ctx context.Context, shortUrl string, urlData []byte) error
	GetOriginalUrl(ctx context.Context, shortUrl string) (*models.CacheData, error)
	Close()
}

type RedisClientService struct {
	redisClient *redis.Client
}

const cacheDuration = 24 * time.Hour

var log = logger.GetInstance()

func InitializeRedis(conf *config.RedisConfig) *RedisClientService {
	rc := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pong, err := rc.Ping(pingCtx).Result()
	if err != nil {
		log.Error("REDIS: Error to connect", "Error", err)
		panic(err)
	}

	log.Info("REDIS: started successfully",
		"pong message", pong)

	return &RedisClientService{redisClient: rc}
}

func (rcs *RedisClientService) SaveUrlMapping(ctx context.Context, shortUrl string, urlData []byte) error {
	err := rcs.redisClient.Set(ctx, shortUrl, urlData, cacheDuration).Err()
	if err != nil {
		log.Error("REDIS: Failed to save key url",
			"error", err,
			"shortUrl", shortUrl,
			"originalUrl", urlData)
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

func (rcs *RedisClientService) GetOriginalUrl(ctx context.Context, shortUrl string) (*models.CacheData, error) {
	jsonData, err := rcs.redisClient.Get(ctx, shortUrl).Result()

	if err == redis.Nil {
		log.Info("REDIS: key does not exist", "shortUrl", shortUrl)
		return nil, nil
	}

	if err != nil {
		log.Error("REDIS: redis Get failed", "shortUrl", shortUrl, "Error", err)
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var data models.CacheData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Error("REDIS: failed to unmarshal cache data", "error", err)
		return nil, fmt.Errorf("unmarshal cache data: %w", err)
	}

	return &data, nil
}

func (rcs *RedisClientService) Close() {
	rcs.redisClient.Close()
}
