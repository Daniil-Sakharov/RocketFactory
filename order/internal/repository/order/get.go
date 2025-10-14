package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (*domain.Order, error) {
	query := `
		SELECT
    		order_uuid,
    		user_uuid,
    		part_uuids,
    		total_price,
    		transaction_uuid,
    		payment_method,
    		order_status
		FROM orders
		WHERE order_uuid = $1;
`

	var repoOrder repoModel.Order
	err := r.db.QueryRowxContext(ctx, query, orderUUID).StructScan(&repoOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return converter.RepoOrderToDomainModel(&repoOrder), nil
}
