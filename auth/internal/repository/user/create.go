package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	repoUser := converter.EntityToRepositoryUser(user)

	query := `
		INSERT INTO users (
			user_uuid,
			login,
			password_hash,
			email,
			notification_methods,
			created_at,
			updated_at
		) VALUES (
			:user_uuid,
			:login,
			:password_hash,
			:email,
			:notification_methods,
			:created_at,
			:updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, repoUser)
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
				if pqErr.Constraint == "users_pkey" {
					return model.ErrUserAlreadyExists
				}
				return fmt.Errorf("unique constraint violation: %s", pqErr.Constraint)

			case "23502":
				return fmt.Errorf("required field is null: %s", pqErr.Column)

			case "23503":
				return fmt.Errorf("foreign key violation: %s", pqErr.Constraint)
			}
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
