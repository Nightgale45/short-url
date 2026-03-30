package main

import (
	"log"
	"net/http"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/handler"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.GetInstance().Info("Starting up the application")

	conf := config.LoadConf()
	redisClient := redis.InitializeRedis(&conf.RedisConf)
	postgresClient := postgres.InitDB(&conf.DatabaseConf)

	defer redisClient.Close()
	defer postgresClient.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{conf.BaseUrl},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/shorten", handler.Shorten(postgresClient, redisClient, conf.BaseUrl))
	v1.GET("/:id", handler.Redirect(redisClient, postgresClient, conf.BaseUrl))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
