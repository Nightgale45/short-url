package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiter())
	r.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	return r
}

func doRequest(r *gin.Engine, ip string) int {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", ip)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

// Burst capacity is 20, so the first 20 requests must pass.
func TestRateLimiter_AllowsWithinBurst(t *testing.T) {
	r := newTestRouter()
	for i := range 20 {
		code := doRequest(r, "1.2.3.4")
		assert.Equal(t, http.StatusOK, code, "request %d should pass", i+1)
	}
}

// After exhausting the burst of 20, the next request must be rejected.
func TestRateLimiter_BlocksAfterBurstExhausted(t *testing.T) {
	r := newTestRouter()
	for range 20 {
		doRequest(r, "1.2.3.4")
	}
	code := doRequest(r, "1.2.3.4")
	assert.Equal(t, http.StatusTooManyRequests, code)
}

// Each IP gets its own independent limiter.
func TestRateLimiter_IndependentLimitersPerIP(t *testing.T) {
	r := newTestRouter()
	for range 20 {
		doRequest(r, "1.2.3.4")
	}

	// A different IP should still have its full burst available.
	code := doRequest(r, "5.6.7.8")
	assert.Equal(t, http.StatusOK, code)
}

// The 429 response body must contain the expected error field.
func TestRateLimiter_ErrorResponseBody(t *testing.T) {
	r := newTestRouter()
	for range 20 {
		doRequest(r, "1.2.3.4")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "rate limit exceeded")
}
