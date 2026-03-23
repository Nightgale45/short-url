package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Nightgale45/short-url/internal/codec"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func Redirect(redis redis.RedisService, db postgres.DbService, baseUrl string) gin.HandlerFunc {
	notFoundUrl := baseUrl + "/not-found"

	return func(ginCtx *gin.Context) {
		key := ginCtx.Param("id")
		ctx := ginCtx.Request.Context()

		data, redisErr := redis.GetOriginalUrl(ctx, key)
		if redisErr != nil {
			logger.GetInstance().Warn("REDIRECT: redis degraded, falling back to db", "error", redisErr)
		}

		if data != nil {
			if !isSafeUrl(data.Data.OriginalUrl) {
				ginCtx.Redirect(http.StatusTemporaryRedirect, notFoundUrl)
				return
			}
			ginCtx.Redirect(http.StatusFound, data.Data.OriginalUrl)
			return
		}

		id, salt := codec.Base62Decoder(key)

		if id == -1 {
			ginCtx.Redirect(http.StatusTemporaryRedirect, notFoundUrl)
			return
		}

		url, dbSalt, err := db.QueryRow(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ginCtx.Redirect(http.StatusTemporaryRedirect, notFoundUrl)
			} else {
				log.Error("REDIRECT: db error", "error", err)
				ginCtx.Status(http.StatusInternalServerError)
			}
			return
		}

		if dbSalt != salt {
			ginCtx.Redirect(http.StatusTemporaryRedirect, notFoundUrl)
			return
		}

		if !isSafeUrl(url) {
			ginCtx.Redirect(http.StatusTemporaryRedirect, notFoundUrl)
			return
		}

		ginCtx.Redirect(http.StatusFound, url)
	}
}

func isSafeUrl(u string) bool {
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}
