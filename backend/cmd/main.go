package main

import (
	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/handler"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.GetInstance().Info("Starting up the application")

	conf := config.LoadConf()
	redis := redis.InitializeRedis(&conf.RedisConf)
	postgres := postgres.InitDB(&conf.DatabaseConf)

	defer redis.Close()
	defer postgres.Close()
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/shorten", handler.Shorten(postgres, redis))

	r.Run(":8080")
}
