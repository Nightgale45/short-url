package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestShortenWithMalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/shorten", Shorten(testhelper.DefaultMockDB(), testhelper.DefaultMockRedis()))

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{invalid json`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShortenWithDBInsertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	db := testhelper.DefaultMockDB()
	db.InsertErr = errors.New("db insert failed")

	r.POST("/shorten", Shorten(db, testhelper.DefaultMockRedis()))

	bodyBytes, err := json.Marshal(defaultMockData())
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShortenRedisFailureContinues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	redis := &testhelper.MockRedisService{SaveErr: errors.New("redis down")}
	r.POST("/shorten", Shorten(testhelper.DefaultMockDB(), redis))

	bodyBytes, err := json.Marshal(defaultMockData())
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
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
	assert.NotEmpty(t, response.ShortenUrl)
	assert.Equal(t, "http://google.com", response.OriginalUrl)

}
