package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nightgale45/short-url/internal/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestRedirectWithRedis(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(testhelper.DefaultMockRedis(), testhelper.DefaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
	assert.Equal(t, "/mock-url", w.Header().Get("Location"))
}

func TestRedirectWithDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(&testhelper.MockRedisService{
		Cache: nil,
	}, testhelper.DefaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/FXsk", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPermanentRedirect, w.Code)
	assert.Equal(t, "/mock-url", w.Header().Get("Location"))
}

func TestUrlIncorrect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/:id", Redirect(&testhelper.MockRedisService{
		Cache: nil,
	}, testhelper.DefaultMockDB()))

	req := httptest.NewRequest(http.MethodGet, "/@", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://localhost:80/not-found", w.Header().Get("Location"))
}

func TestDBUrlMiss(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := testhelper.DefaultMockDB()
	db.Err = pgx.ErrNoRows

	r.GET("/:id", Redirect(&testhelper.MockRedisService{Cache: nil}, db))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(t, "http://localhost:80/not-found", w.Header().Get("Location"))
}
