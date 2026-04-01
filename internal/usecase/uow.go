package usecase

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err := fn(q); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
