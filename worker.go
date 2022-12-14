package main

import (
	"context"
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

type Worker struct {
	logger *zap.SugaredLogger
	config util.Config
	dbconn *pgx.Conn
}

func NewWorker(logger *zap.SugaredLogger, config util.Config) *Worker {
	return &Worker{
		logger: logger,
		config: config,
	}
}

func (w *Worker) Initialize() {
	createKwekChannel := make(chan *kwekker_protobufs.CreateKwek)
	updateKwekChannel := make(chan *kwekker_protobufs.UpdateKwek)
	deleteKwekChannel := make(chan *kwekker_protobufs.DeleteKwek)

	rabbitMQWorker := NewRabbitMQWorker(w.logger, w.config.RabbitMQ)
	go rabbitMQWorker.Listen(createKwekChannel, updateKwekChannel, deleteKwekChannel)

	db := NewDB(w.logger, w.config.Postgres)
	w.dbconn = db.Connect()
	defer w.dbconn.Close(context.Background())

	for {
		select {
		case createKwek := <-createKwekChannel:
			w.logger.Debug("Received create kwek request", "kwek", createKwek)
			w.handleCreateKwek(createKwek)
		case updateKwek := <-updateKwekChannel:
			w.logger.Debug("Received update kwek request", "kwek", updateKwek)
			w.handleUpdateKwek(updateKwek)
		case deleteKwek := <-deleteKwekChannel:
			w.logger.Debug("Received delete kwek request", "kwek", deleteKwek)
			w.handleDeleteKwek(deleteKwek)
		}
	}
}

func (w *Worker) handleCreateKwek(createKwek *kwekker_protobufs.CreateKwek) {
	w.logger.Debug("Handling create kwek request", "kwek", createKwek)

	_, err := w.dbconn.Exec(
		context.Background(),
		`INSERT INTO "Kweks" ("Guid", "UserId", "Text", "PostedAt")
			 VALUES ($1, (SELECT "Id" FROM "Users" WHERE "ProviderId" = $2), $3, $4)`,
		createKwek.GetKwekGuid(),
		createKwek.GetUserId(),
		createKwek.GetText(),
		createKwek.GetPostedAt().AsTime(),
	)

	if err != nil {
		w.logger.Error("Failed to insert kwek into database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully inserted kwek into database")
}

func (w *Worker) handleUpdateKwek(updateKwek *kwekker_protobufs.UpdateKwek) {
	w.logger.Debug("Handling update kwek request", "kwek", updateKwek)

	_, err := w.dbconn.Exec(
		context.Background(),
		`UPDATE "Kweks" SET "Text" = $1 WHERE "Guid" = $2`,
		updateKwek.GetText(),
		updateKwek.GetKwekGuid(),
	)

	if err != nil {
		w.logger.Error("Failed to update kwek into database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully updated kwek into database")
}

func (w *Worker) handleDeleteKwek(deleteKwek *kwekker_protobufs.DeleteKwek) {
	w.logger.Debug("Handling delete kwek request", "kwek", deleteKwek)

	_, err := w.dbconn.Exec(
		context.Background(),
		`DELETE FROM "Kweks" WHERE "Guid" = $1`,
		deleteKwek.GetKwekGuid(),
	)

	if err != nil {
		w.logger.Error("Failed to delete kwek into database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully deleted kwek into database")
}
