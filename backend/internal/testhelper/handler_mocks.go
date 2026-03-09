package testhelper

import (
	"context"
	"time"

	"github.com/Nightgale45/short-url/internal/models"
)

//============= Redis services =================j

func DefaultMockRedis() *MockRedisService {
	return &MockRedisService{
		Cache: &models.CacheData{
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
	Cache *models.CacheData
}

func (m *MockRedisService) GetOriginalUrl(ctx context.Context, shortUrl string) *models.CacheData {
	return m.Cache
}

func (m *MockRedisService) SaveUrlMapping(ctx context.Context, shortUrl string, urlData []byte) {

}

func (m *MockRedisService) Close() {}

//============== DB services ====================

type MockDbService struct {
	Url  string
	Salt int64
	Err  error
	Id   int64
}

func (db *MockDbService) QueryRow(ctx context.Context, id int64) (string, int64, error) {
	return db.Url, db.Salt, db.Err
}

func (db *MockDbService) InsertUrlData(ctx context.Context, data models.UrlData) int64 {
	return db.Id
}

func (db *MockDbService) Close() {}

func DefaultMockDB() *MockDbService {
	return &MockDbService{
		Url:  "mock-url",
		Salt: int64(10_000_000),
		Err:  nil,
		Id:   int64(1),
	}
}
