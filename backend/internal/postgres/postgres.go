package postgres

import (
	"context"

	"github.com/Nightgale45/short-url/internal/config"
	"github.com/Nightgale45/short-url/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dbConf *config.DatabaseConfig) *pgxpool.Pool {
	conf, err := pgxpool.ParseConfig(dbConf.Url)
	if err != nil {
		logger.GetInstance().Error("POSTGRES: Cannot create db config", "Error", err)
		panic(err)
	}

	conf.MaxConns = int32(dbConf.MaxConns)
	conf.MinConns = int32(dbConf.MinConns)

	dbpool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		logger.GetInstance().Error("POSTGRES: Cannot connect to db", "Error", err)
		panic(err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.GetInstance().Error("POSTGRES: Cannot ping database", "Error", err)
		panic(err)
	}

	logger.GetInstance().Info("POSTGRES: Successful ping of db")

	return dbpool
}
