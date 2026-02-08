package redis

import (
	"context"
	"time"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/redis/go-redis/v9"
)

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

	pong, err := rc.Ping(context.Background()).Result()
	if err != nil {
		log.Error("REDIS: Error to connect", "Error", err)
		panic(err)
	}

	log.Info("REDIS: started successfully",
		"pong message", pong)

	return &RedisClientService{redisClient: rc}
}

func (rcs *RedisClientService) SaveUrlMapping(ctx context.Context, shortUrl string, originalUrl string) {
	err := rcs.redisClient.Set(ctx, shortUrl, originalUrl, cacheDuration)
	if err != nil {
		log.Error("REDIS: Failed to save key url",
			"error", err,
			"shortUrl", shortUrl,
			"originalUrl", originalUrl)
	}

}

func (rcs *RedisClientService) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	url, err := rcs.redisClient.Get(ctx, shortUrl).Result()

	if err == redis.Nil {
		log.Info("REDIS: key does not exist",
			"shortUrl", shortUrl)
		return "", nil

	} else if err != nil {
		log.Error("REDIS: redis Get failed",
			"shortUrl", shortUrl,
			"Error", err)

		return "", nil
	}

	return url, err
}

func (rcs *RedisClientService) Close() {
	rcs.redisClient.Close()
}
