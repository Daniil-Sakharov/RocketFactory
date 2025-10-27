package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order *domain.Order) error {
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

	updateQuery := `
        UPDATE orders
        SET 
            order_status = :order_status,
            payment_method = :payment_method,
            transaction_uuid = :transaction_uuid,
            updated_at = NOW()
        WHERE order_uuid = :order_uuid
    `

	result, err := tx.NamedExecContext(ctx, updateQuery, repoOrder)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get Rows Affect")
	}
	if rowsAffected == 0 {
		var currentStatus string
		checkErr := tx.GetContext(ctx, &currentStatus,
			`SELECT order_status FROM orders WHERE order_uuid = $1`,
			order.OrderUUID,
		)
		if errors.Is(checkErr, sql.ErrNoRows) {
			return model.ErrOrderNotFound
		}

		if checkErr != nil {
			return fmt.Errorf("failed to check order status: %w", checkErr)
		}
		switch currentStatus {
		case "PAID":
			return model.ErrOrderAlreadyPaid
		case "CANCELLED":
			return model.ErrOrderAlreadyCancelled
		default:
			return fmt.Errorf("unexpected order status: %s", currentStatus)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}
