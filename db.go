package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"kwekker-worker/util"
	"time"
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

	var conn *pgx.Conn

	for i := 0; i < 5; i++ {
		var err error
		conn, err = pgx.Connect(
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

		if err == nil {
			break
		}

		db.logger.Debug("Retrying database connection", "error", err)
		time.Sleep(5 * time.Second)
	}

	if conn == nil {
		db.logger.Fatal("Failed to connect to DB")
	}

	db.logger.Debug("Connected to DB")

	return conn
}
