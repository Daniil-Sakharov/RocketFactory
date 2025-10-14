package order

import (
	"github.com/jmoiron/sqlx"

	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{
		db: db,
	}
}
