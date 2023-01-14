package worker

import (
	"context"
	"fmt"
	userproto "github.com/googolplex-s6/kwekker-protobufs/v3/user"
	"go.uber.org/zap"
)

func (w *Worker) handleCreateUser(createUser *userproto.CreateUser) {
	w.logger.Debug("Handling create user request", "user", createUser)

	_, err := w.dbconn.Exec(
		context.Background(),
		`INSERT INTO "Users" ("ProviderId", "Username", "Email", "DisplayName", "AvatarUrl")
			 VALUES ($1, $2, $3, $4, $5)`,
		createUser.GetUserId(),
		createUser.GetUsername(),
		createUser.GetEmail(),
		createUser.GetDisplayName(),
		createUser.GetAvatarUrl(),
	)

	if err != nil {
		w.logger.Error("Failed to insert user into database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully inserted user into database")
}

func (w *Worker) handleUpdateUser(updateUser *userproto.UpdateUser) {
	w.logger.Debug("Handling update user request", "user", updateUser)

	updatedFields := make(map[string]string, 0)

	if updateUser.GetUsername() != "" {
		updatedFields["Username"] = updateUser.GetUsername()
	}

	if updateUser.GetEmail() != "" {
		updatedFields["Email"] = updateUser.GetEmail()
	}

	if updateUser.GetDisplayName() != "" {
		updatedFields["DisplayName"] = updateUser.GetDisplayName()
	}

	if updateUser.GetAvatarUrl() != "" {
		updatedFields["AvatarUrl"] = updateUser.GetAvatarUrl()
	}

	if len(updatedFields) == 0 {
		w.logger.Debug("No fields to update")
		return
	}

	query := `UPDATE "Users" SET `
	values := []any{updateUser.GetUserId()}

	i := 2
	for field, value := range updatedFields {
		query += fmt.Sprintf(`"%s" = $%d,`, field, i)
		values = append(values, value)
		i++
	}

	query = fmt.Sprintf(`%s WHERE "ProviderId" = $1`, query[:len(query)-1])

	_, err := w.dbconn.Exec(
		context.Background(),
		query,
		values...,
	)

	if err != nil {
		w.logger.Error("Failed to update kwek in database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully updated kwek in database")
}

func (w *Worker) handleDeleteUser(deleteUser *userproto.DeleteUser) {
	w.logger.Debug("Handling delete user request", "user", deleteUser)

	_, err := w.dbconn.Exec(
		context.Background(),
		`DELETE FROM "Users" WHERE "ProviderId" = $1`,
		deleteUser.GetUserId(),
	)

	if err != nil {
		w.logger.Error("Failed to delete user in database", zap.Error(err))
		return
	}

	w.logger.Debug("Successfully deleted user in database")
}
