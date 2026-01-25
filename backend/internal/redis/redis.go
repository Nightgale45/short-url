package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientService struct {
	redisClient *redis.Client
}

var (
	ctx          = context.Background()
	RedisService = &RedisClientService{}
)

const CacheDuration = 6 * time.Hour

func InitializeRedis() error {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)

	RedisService.redisClient = RedisClient

	return nil
}

func (rds *RedisClientService) SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := rds.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration)
	if err != nil {
		fmt.Printf("Failed to save key url | Error: %v - shortUrl: %s - originalUrl: %s", err, shortUrl, originalUrl)
	}

}

func (rds *RedisClientService) GetOriginalUrl(shortUrl string) (string, error) {
	url, err := rds.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis: url is not present | shortUrl: %s | Error: %v", shortUrl, err))
	}

	return url, err
}
