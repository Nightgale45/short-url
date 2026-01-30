package redis

import (
	"context"
	"log"
	"time"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClientService struct {
	redisClient *redis.Client
}

const CacheDuration = 6 * time.Hour

func InitializeRedis(conf *config.RedisConfig) (*RedisClientService, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	pong, err := rc.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Printf("\nRedis started successfully: pong message = {%s}\n", pong)

	return &RedisClientService{redisClient: rc}, nil
}

func (rcs *RedisClientService) SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := rcs.redisClient.Set(context.Background(), shortUrl, originalUrl, CacheDuration)
	if err != nil {
		log.Printf("Failed to save key url | Error: %v - shortUrl: %s - originalUrl: %s", err, shortUrl, originalUrl)
	}

}

func (rcs *RedisClientService) GetOriginalUrl(shortUrl string) (string, error) {
	url, err := rcs.redisClient.Get(context.Background(), shortUrl).Result()
	if err != nil {
		panic(log.Sprintf("Redis: url is not present | shortUrl: %s | Error: %v", shortUrl, err))
	}

	return url, err
}
