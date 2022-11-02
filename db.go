package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

type DB struct {
	logger *zap.SugaredLogger
	config util.PostgresConfig
}

func NewDB(logger *zap.SugaredLogger, config util.PostgresConfig) *DB {
	return &DB{
		logger: logger,
		config: config,
	}
}

func (db *DB) Connect() *pgx.Conn {
	db.logger.Debug("Connecting to DB")

	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			db.config.Username,
			db.config.Password,
			db.config.Host,
			db.config.Port,
			db.config.Database,
		),
	)

	if err != nil {
		db.logger.Fatal("Failed to connect to DB", zap.Error(err))
	}

	db.logger.Debug("Connected to DB")

	return conn
}
