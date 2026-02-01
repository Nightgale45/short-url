package handler

import (
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Receive a url and create a shorten url to return
func Shorten(db *pgxpool.Pool, redis *redis.RedisClientService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.JSON(200, &models.ShortKey{
			ShortUrl:    "short key",
			OriginalUrl: "original key",
			CreatedAt:   time.Now(),
		})
	}
}
