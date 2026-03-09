package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nightgale45/short-url/internal/models"
	"github.com/Nightgale45/short-url/internal/testhelper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func defaultMockData() *models.ShortenRequest {
	return &models.ShortenRequest{
		OriginalUrl: "http://google.com",
	}
}

func TestUrlValidationWithIncorrectUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shorten", Shorten(testhelper.DefaultMockDB(), testhelper.DefaultMockRedis()))

	body := &models.ShortenRequest{
		OriginalUrl: "test",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShortenWithCorrectUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shorten", Shorten(testhelper.DefaultMockDB(), testhelper.DefaultMockRedis()))

	bodyBytes, err := json.Marshal(defaultMockData())
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ShortenResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, *response.ShortenUrl)
	assert.Equal(t, "http://google.com", response.OriginalUrl)

}
