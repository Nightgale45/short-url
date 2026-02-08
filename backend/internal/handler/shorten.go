package handler

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"regexp"
	"strings"

	"github.com/Nightgale45/short-url/internal/codec"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/postgres"
	"github.com/Nightgale45/short-url/internal/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// allow for subdomain but does not allow the beginning or end of domain with hyphen
const pattern = `^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`

var log = logger.GetInstance()

// Receive a url and create a shorten url to return
func Shorten(db *postgres.DbPool, redis *redis.RedisClientService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var shortenRequest models.ShortenRequest
		ctx := ginCtx.Request.Context()

		// should use should bind to consume the request and assugn to the var
		err := ginCtx.ShouldBindJSON(&shortenRequest)
		if err != nil {
			errMsg := "Malform request"
			log.Error("SHORTEN: Error binding json request to struct", "Error", err)
			ginCtx.JSON(400, generateResponse(shortenRequest.OriginalUrl, nil, &errMsg))
			return
		}

		if !validateUrl(shortenRequest.OriginalUrl) {
			errMsg := "Error url is invalid"
			ginCtx.JSON(400, generateResponse(shortenRequest.OriginalUrl, nil, &errMsg))
		}

		// encryptPass could be nil
		var encryptPass *string

		if shortenRequest.Passcode != nil {

			hashed, err := bcrypt.GenerateFromPassword([]byte(*shortenRequest.Passcode), 10)
			if err != nil {
				errMsg := "Error generating url"
				log.Error("SHORTEN: Error hashing passcode", "Error", err)
				ginCtx.JSON(500, generateResponse(shortenRequest.OriginalUrl, nil, &errMsg))
				return
			}

			hashStr := string(hashed)
			encryptPass = &hashStr
		}

		randNum, err := rand.Int(rand.Reader, big.NewInt(90_000_000))

		if err != nil {
			errMsg := "Error generating url"
			log.Error("SHORTEN: Error generating salt number", "Error", err)
			ginCtx.JSON(500, generateResponse(shortenRequest.OriginalUrl, nil, &errMsg))
			return
		}

		salt := randNum.Int64() + int64(10_000_000)
		dbId := db.InsertUrl(ctx, shortenRequest.OriginalUrl, salt, encryptPass)

		shortenKey := codec.Base62Encoder(dbId, salt)
		redis.SaveUrlMapping(ctx, shortenKey, shortenRequest.OriginalUrl)

		ginCtx.JSON(200, generateResponse(shortenRequest.OriginalUrl, &shortenKey, nil))
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
	if !strings.HasPrefix(scheme, "http") && !strings.HasPrefix(scheme, "https") {
		return false
	}

	match, err := regexp.MatchString(pattern, u.Host)
	if err != nil {
		log.Info("SHORTEN: error with regex url", "Error", err)
		return false
	}

	return match
}

func generateResponse(originalUrl string, shortenUrl *string, errorMessage *string) models.ShortenResponse {
	return models.ShortenResponse{
		OriginalUrl: originalUrl,
		ShortenUrl:  shortenUrl,
		Error:       errorMessage,
	}
}
