package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, user *entity.User) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				if !errors.Is(rbErr, sql.ErrTxDone) {
					err = fmt.Errorf("tx rollback failed: %w, original error: %w", rbErr, err)
				}
			}
		}
	}()

	repoUser := converter.EntityToRepositoryUser(user)

	var exist bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE user_uuid = $1)`

	err = tx.QueryRowContext(ctx, checkQuery, repoUser.UserUUID).Scan(&exist)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !exist {
		return model.ErrUserNotFound
	}

	updateQuery := `
		UPDATE users
		SET
			login = :login,
			password_hash = :password_hash,
			email = :email,
			notification_methods = :notification_methods,
			updated_at = :updated_at
		WHERE user_uuid = :user_uuid
`

	result, err := tx.NamedExecContext(ctx, updateQuery, repoUser)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "users_email_key" {
					return model.ErrEmailAlreadyExists
				}
				if pqErr.Constraint == "users_login_key" {
					return model.ErrLoginAlreadyExists
				}
				return fmt.Errorf("unique constraint violation: %w", err)
			case "23503":
				return fmt.Errorf("foreign key constraint violation: %w", err)
			case "23502":
				return fmt.Errorf("not null constraint violation: %w", err)
			}
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
