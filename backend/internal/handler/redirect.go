package handler

import (
	"net/http"

	"github.com/Nightgale45/short-url/internal/codec"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
)

func Redirect(redis redis.RedisService, db postgres.DbService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		key := ginCtx.Param("id")
		ctx := ginCtx.Request.Context()

		data := redis.GetOriginalUrl(ctx, key)

		if data != nil {
			ginCtx.Redirect(http.StatusPermanentRedirect, data.Data.OriginalUrl)
			return
		} else {

			id, salt := codec.Base62Decoder(key)

			if id != -1 {
				url, dbSalt, err := db.QueryRow(ctx, id)

				if err != nil {
					ginCtx.Redirect(http.StatusTemporaryRedirect, "http://localhost:80/not-found")
					return
				}

				if dbSalt == salt {
					ginCtx.Redirect(http.StatusPermanentRedirect, url)
				}
			} else {
				// Add a env for base url
				ginCtx.Redirect(http.StatusTemporaryRedirect, "http://localhost:80/not-found")
				return
			}
		}

	}
}
