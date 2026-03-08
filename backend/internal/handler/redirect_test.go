package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nightgale45/short-url/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func defaultMockRedis() *MockRedisService {
	return &MockRedisService{
		cache: &models.CacheData{
			ShortenKey: "mock key",
			Data: models.UrlData{
				OriginalUrl: "mock-url",
				CreatedAt:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				Salt:        int64(123),
			},
		},
	}
}

type MockRedisService struct {
	cache *models.CacheData
}

func (m *MockRedisService) GetOriginalUrl(ctx context.Context, shortUrl string) *models.CacheData {
	return m.cache
}

func (m *MockRedisService) SaveUrlMapping(ctx context.Context, shortUrl string, urlData []byte) {

}

func (m *MockRedisService) Close() {}

type MockDbService struct {
	url  string
	salt int64
	err  error
	id   int64
}

func (db *MockDbService) QueryRow(ctx context.Context, id int64) (string, int64, error) {
	return db.url, db.salt, db.err
}

func (db *MockDbService) InsertUrlData(ctx context.Context, data models.UrlData) int64 {
	return db.id
}

func (db *MockDbService) Close() {}

func defaultMockDB() *MockDbService {
	return &MockDbService{
		url:  "mock-url",
		salt: int64(10_000_000),
		err:  nil,
		id:   int64(1),
	}
}

func TestRedirectWithRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(defaultMockRedis(), defaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
	assert.Equal(t, "/mock-url", w.Header().Get("Location"))
}

func TestRedirectWithDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(&MockRedisService{
		cache: nil,
	}, defaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
	assert.Equal(t, "/mock-url", w.Header().Get("Location"))
}

func TestUrlIncorrect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(&MockRedisService{
		cache: nil,
	}, defaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/@", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://localhost:80/not-found", w.Header().Get("Location"))
}

func TestDBUrlMiss(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := defaultMockDB()
	db.err = pgx.ErrNoRows

	r.GET("/:id", Redirect(&MockRedisService{cache: nil}, db))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://localhost:80/not-found", w.Header().Get("Location"))
}
