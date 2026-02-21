package postgres

import (
	"context"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/Nightgale45/short-url/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbPool struct {
	dbPool *pgxpool.Pool
}

var log = logger.GetInstance()

func InitDB(dbConf *config.DatabaseConfig) *DbPool {
	conf, err := pgxpool.ParseConfig(dbConf.Url)
	if err != nil {
		log.Error("POSTGRES: Cannot create db config", "Error", err)
		panic(err)
	}

	conf.MaxConns = int32(dbConf.MaxConns)
	conf.MinConns = int32(dbConf.MinConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		log.Error("POSTGRES: Cannot connect to db", "Error", err)
		panic(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Error("POSTGRES: Cannot ping database", "Error", err)
		panic(err)
	}

	log.Info("POSTGRES: Successful ping of db")

	return &DbPool{
		dbPool: pool,
	}
}

func (db *DbPool) InsertUrlData(ctx context.Context, data models.UrlData) int64 {
	sql := `INSERT INTO url_data (original_url, created_at, counter, passcode) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64

	db.dbPool.QueryRow(ctx, sql, data.OriginalUrl, data.CreatedAt, data.Counter, data.Passcode).Scan(&id)
	return id
}

func (db *DbPool) UpdateData(ctx context.Context, id int64, count int) {
	sql := `UPDATE url_data SET count = $1 WHERE id = $2`

	db.dbPool.QueryRow(ctx, sql, id, count)
}

func (db *DbPool) QueryRow(ctx context.Context, id int64) (*models.RedirectData, error) {
	sql := `SELECT original_url, salt, counter FROM url_data WHERE id = $1`

	var url string
	var salt int64
	var counter int

	err := db.dbPool.QueryRow(ctx, sql, id).Scan(&url, &salt, &counter)
	if err != nil {
		return nil, err
	}

	return &models.RedirectData{
		OriginalUrl: url,
		Salt:        salt,
		Counter:     counter,
	}, nil
}

func (db *DbPool) Close() {
	db.dbPool.Close()
}
