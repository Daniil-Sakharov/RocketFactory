package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (*entity.User, error) {
	query := `
	SELECT 
	    user_uuid,
		login,
		password_hash,
		email,
		notification_methods,
		created_at,
		updated_at
	FROM users
	WHERE user_uuid = $1;
`
	var repoUser repoModel.User
	err := r.db.QueryRowxContext(ctx, query, uuid).StructScan(&repoUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return converter.RepositoryToEntityUser(&repoUser), nil
}
