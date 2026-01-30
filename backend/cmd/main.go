package main

import (
	"log"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting up the application")
	conf := config.LoadConf()

	_, err := redis.InitializeRedis(&conf.RedisConf)
	if err != nil {
		log.Panicf("Redis: error to connect - %v", err)
	}

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.Run(":8080")
}
