package worker

import (
	"context"
	kwekkerprotobufs "github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	"go.uber.org/zap"
)

func (w *Worker) handleCreateKwek(createKwek *kwekkerprotobufs.CreateKwek) {
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

func (w *Worker) handleUpdateKwek(updateKwek *kwekkerprotobufs.UpdateKwek) {
	w.logger.Debug("Handling update kwek request", "kwek", updateKwek)

	_, err := w.dbconn.Exec(
		context.Background(),
		`UPDATE "Kweks" SET "Text" = $1 WHERE "Guid" = $2`,
		updateKwek.GetText(),
		updateKwek.GetKwekGuid(),
	)

	if err != nil {
		w.logger.Error("Failed to update kwek in database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully updated kwek in database")
}

func (w *Worker) handleDeleteKwek(deleteKwek *kwekkerprotobufs.DeleteKwek) {
	w.logger.Debug("Handling delete kwek request", "kwek", deleteKwek)

	_, err := w.dbconn.Exec(
		context.Background(),
		`DELETE FROM "Kweks" WHERE "Guid" = $1`,
		deleteKwek.GetKwekGuid(),
	)

	if err != nil {
		w.logger.Error("Failed to delete kwek in database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully deleted kwek in database")
}
