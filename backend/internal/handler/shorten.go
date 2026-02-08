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

// Receive a url and create a shorten url to return
func Shorten(db *postgres.DbPool, redis *redis.RedisClientService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var shortenRequest models.ShortenRequest

		// should use should bind to consume the request and assugn to the var
		if err := ctx.ShouldBindJSON(&shortenRequest); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if validateUrl(shortenRequest.OriginalUrl) {
			var encryptPass *string

			if shortenRequest.Passcode != nil {
				hashed, err := bcrypt.GenerateFromPassword([]byte(*shortenRequest.Passcode), 10)
				if err != nil {
					ctx.JSON(400, gin.H{"error": "Error encoding passcode"})
					return
				}

				hashStr := string(hashed)
				encryptPass = &hashStr
			}

			randNum, err := rand.Int(rand.Reader, big.NewInt(90_000_000))
			if err != nil {
				ctx.JSON(400, gin.H{"error": "Error generating salt"})
				return
			}

			salt := randNum.Int64() + int64(10_000_000)
			dbId := db.InsertUrl(ctx.Request.Context(), shortenRequest.OriginalUrl, salt, encryptPass)

			codec.Base62Encoder(dbId, salt)

		} else {
			ctx.JSON(400, gin.H{"error": "Url is invalid"})
		}
	}
}

func validateUrl(userUrl string) bool {
	// clean the url
	cleanUrl := strings.TrimSpace(userUrl)

	u, err := url.Parse(cleanUrl)
	if err != nil {
		logger.GetInstance().Error("SHORTEN: error parsing user url", "Error", err)
		return false
	}

	scheme := u.Scheme
	if !strings.HasPrefix(scheme, "http") && !strings.HasPrefix(scheme, "https") {
		return false
	}

	match, err := regexp.MatchString(pattern, u.Host)
	if err != nil {
		logger.GetInstance().Error("SHORTEN: error with regex url", "Error", err)
		return false
	}

	return match
}
