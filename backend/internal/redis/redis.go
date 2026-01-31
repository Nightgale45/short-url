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

const CacheDuration = 6 * time.Hour

func InitializeRedis(conf *config.RedisConfig) *RedisClientService {
	rc := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	pong, err := rc.Ping(context.Background()).Result()
	if err != nil {
		logger.GetInstance().Error("REDIS: Error to connect", "Error", err)
		panic(err)
	}

	logger.GetInstance().Info("REDIS: started successfully",
		"pong message", pong)

	return &RedisClientService{redisClient: rc}
}

func (rcs *RedisClientService) SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := rcs.redisClient.Set(context.Background(), shortUrl, originalUrl, CacheDuration)
	if err != nil {
		logger.GetInstance().Error("REDIS: Failed to save key url",
			"error", err,
			"shortUrl", shortUrl,
			"originalUrl", originalUrl)
	}

}

func (rcs *RedisClientService) GetOriginalUrl(shortUrl string) (string, error) {
	url, err := rcs.redisClient.Get(context.Background(), shortUrl).Result()

	if err == redis.Nil {
		logger.GetInstance().Info("REDIS: key does not exist",
			"shortUrl", shortUrl)
		return "", nil

	} else if err != nil {
		logger.GetInstance().Error("REDIS: GetOrginalUrl failed",
			"shortUrl", shortUrl,
			"Error", err)

		return "", nil
	}

	return url, err
}
