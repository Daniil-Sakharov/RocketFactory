package user

import (
	"github.com/jmoiron/sqlx"

	def "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
)

var _ def.UsersRepository = (*repository)(nil)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}
