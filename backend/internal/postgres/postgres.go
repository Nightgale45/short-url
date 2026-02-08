package postgres

import (
	"context"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/logger"
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

func (db *DbPool) InsertUrl(ctx context.Context, url string, salt int64, passcode *string) int64 {
	sql := `INSERT INTO urls (original_url, salt, passcode) VALUES ($1, $2, $3) RETURNING id`

	var id int64

	db.dbPool.QueryRow(ctx, sql, url, salt, passcode).Scan(&id)
	return id
}

func (db *DbPool) Close() {
	db.dbPool.Close()
}
