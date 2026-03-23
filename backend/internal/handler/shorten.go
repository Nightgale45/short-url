package handler

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Nightgale45/short-url/internal/codec"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
)

// allow for subdomain but does not allow the beginning or end of domain with hyphen
const pattern = `^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`

var (
	log         = logger.GetInstance()
	hostPattern = regexp.MustCompile(pattern)
)

// Receive a url and create a shorten url to return
func Shorten(db postgres.DbService, redis redis.RedisService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var shortenRequest models.ShortenRequest
		ctx := ginCtx.Request.Context()

		// should use should bind to consume the request and assugn to the var
		err := ginCtx.ShouldBindJSON(&shortenRequest)
		if err != nil {
			log.Error("SHORTEN: Error binding json request to struct", "Error", err)
			ginCtx.JSON(http.StatusBadRequest, generateResponse(shortenRequest.OriginalUrl, ""))
			return
		}

		if !validateUrl(shortenRequest.OriginalUrl) {
			ginCtx.JSON(http.StatusBadRequest, generateResponse(shortenRequest.OriginalUrl, ""))
			return
		}

		randNum, err := rand.Int(rand.Reader, big.NewInt(90_000_000))
		if err != nil {
			log.Error("SHORTEN: Error generating salt number", "Error", err)
			ginCtx.JSON(http.StatusInternalServerError, generateResponse(shortenRequest.OriginalUrl, ""))
			return
		}

		salt := randNum.Int64() + int64(10_000_000)

		urlData := models.UrlData{
			OriginalUrl: shortenRequest.OriginalUrl,
			CreatedAt:   time.Now().UTC(),
			Salt:        salt,
		}

		dbId, err := db.InsertUrlData(ctx, urlData)
		if err != nil {
			log.Error("SHORTEN: Error inserting url data", "Error", err)
			ginCtx.JSON(http.StatusInternalServerError, generateResponse(shortenRequest.OriginalUrl, ""))
			return
		}

		shortenKey := codec.Base62Encoder(dbId, salt)

		jResp, err := json.Marshal(models.CacheData{
			ShortenKey: shortenKey,
			Data:       urlData,
		})

		if err != nil {
			log.Error("SHORTEN: Error converting response to json", "Error", err)
			ginCtx.JSON(http.StatusInternalServerError, generateResponse(shortenRequest.OriginalUrl, ""))
			return
		}

		if err := redis.SaveUrlMapping(ctx, shortenKey, jResp); err != nil {
			log.Warn("SHORTEN: Failed to populate cache, continuing", "Error", err)
		}

		ginCtx.JSON(http.StatusOK, generateResponse(shortenRequest.OriginalUrl, shortenKey))
	}
}

func validateUrl(userUrl string) bool {
	// clean the url
	cleanUrl := strings.TrimSpace(userUrl)

	u, err := url.Parse(cleanUrl)
	if err != nil {
		log.Info("SHORTEN: error parsing user url", "Error", err)
		return false
	}

	scheme := u.Scheme
	if scheme != "http" && scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return hostPattern.MatchString(u.Host)
}

func generateResponse(originalUrl string, shortenUrl string) models.ShortenResponse {
	return models.ShortenResponse{
		OriginalUrl: originalUrl,
		ShortenUrl:  shortenUrl,
	}
}
