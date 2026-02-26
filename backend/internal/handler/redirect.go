package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nightgale45/short-url/internal/codec"
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
)

func Redirect(redis *redis.RedisClientService, db *postgres.DbPool) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		key := ginCtx.Param("id")
		ctx := ginCtx.Request.Context()

		cache := redis.GetOriginalUrl(ctx, key)

		if cache != nil {
			var data models.CacheData
			json.Unmarshal(cache, &data)

			ginCtx.Redirect(http.StatusAccepted, data.Data.OriginalUrl)

		} else {

			id, salt := codec.Base62Decoder(key)

			url, dbSalt, err := db.QueryRow(ctx, id)
			if err != nil {
				ginCtx.JSON(400, "Data not available")
				return
			}

			if dbSalt == salt {
				ginCtx.Redirect(http.StatusPermanentRedirect, url)
			}
		}

	}
}
