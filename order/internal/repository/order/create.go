package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order *domain.Order) error {
	repoOrder := converter.DomainOrderToRepoModel(order)

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("rollback error: %v\n", err)
		}
	}()

	query := `
        INSERT INTO orders (
            order_uuid,
            user_uuid,
            part_uuids,
            total_price,
            payment_method,
            order_status,
            transaction_uuid
        ) VALUES (
            :order_uuid,
            :user_uuid,
            :part_uuids,
            :total_price,
            :payment_method,
            :order_status,
            :transaction_uuid
        )
    `

	result, err := tx.NamedExecContext(ctx, query, repoOrder)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return model.ErrOrderAlreadyExist
			}
		}
		return fmt.Errorf("failed to insert order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows inserted")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
