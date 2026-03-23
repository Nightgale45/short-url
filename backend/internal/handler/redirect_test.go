package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

const testBaseUrl = "http://localhost"

func TestRedirectWithRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(testhelper.DefaultMockRedis(), testhelper.DefaultMockDB(), testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "http://mock-url.com", w.Header().Get("Location"))
}

func TestRedirectWithDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// "9L7Co" encodes id=1, salt=10_000_000 -- matches DefaultMockDB
	r.GET("/:id", Redirect(&testhelper.MockRedisService{
		Cache: nil,
	}, testhelper.DefaultMockDB(), testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/9L7Co", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "http://mock-url.com", w.Header().Get("Location"))
}

func TestUrlIncorrect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(&testhelper.MockRedisService{
		Cache: nil,
	}, testhelper.DefaultMockDB(), testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/@", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, testBaseUrl+"/not-found", w.Header().Get("Location"))
}

func TestRedirectSaltMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// "9L7Co" decodes to id=1, salt=10_000_000 but DB returns a different salt
	db := &testhelper.MockDbService{
		Url:  "http://mock-url.com",
		Salt: int64(99999),
	}
	r.GET("/:id", Redirect(&testhelper.MockRedisService{Cache: nil}, db, testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/9L7Co", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, testBaseUrl+"/not-found", w.Header().Get("Location"))
}

func TestRedirectUnsafeUrlFromCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	redis := &testhelper.MockRedisService{
		Cache: &models.CacheData{
			ShortenKey: "mock key",
			Data: models.UrlData{
				OriginalUrl: "javascript:alert(1)",
				CreatedAt:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				Salt:        int64(123),
			},
		},
	}
	r.GET("/:id", Redirect(redis, testhelper.DefaultMockDB(), testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, testBaseUrl+"/not-found", w.Header().Get("Location"))
}

func TestRedirectUnsafeUrlFromDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// "9L7Co" decodes to id=1, salt=10_000_000
	db := &testhelper.MockDbService{
		Url:  "javascript:alert(1)",
		Salt: int64(10_000_000),
	}
	r.GET("/:id", Redirect(&testhelper.MockRedisService{Cache: nil}, db, testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/9L7Co", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, testBaseUrl+"/not-found", w.Header().Get("Location"))
}

func TestRedirectDBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := testhelper.DefaultMockDB()
	db.Err = errors.New("connection timeout")

	r.GET("/:id", Redirect(&testhelper.MockRedisService{Cache: nil}, db, testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/9L7Co", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDBUrlMiss(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := testhelper.DefaultMockDB()
	db.Err = pgx.ErrNoRows

	r.GET("/:id", Redirect(&testhelper.MockRedisService{Cache: nil}, db, testBaseUrl))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, testBaseUrl+"/not-found", w.Header().Get("Location"))
}
