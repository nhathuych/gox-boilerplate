package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/pkg/pgxutil"
	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
)

type UserRepository struct {
	q sqlc.Querier
}

func NewUserRepository(q sqlc.Querier) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row, err := r.q.GetUserByID(ctx, pgxutil.ToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return mapUserRow(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return mapUserEmailRow(row), nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User, passwordHash string) error {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        u.Email,
		PasswordHash: passwordHash,
		RoleID:       u.RoleID,
	})
	if err != nil {
		return err
	}
	id, _ := pgxutil.FromPgUUID(row.ID)
	u.ID = id
	u.PasswordHash = row.PasswordHash
	u.CreatedAt = pgxutil.FromTimestamptz(row.CreatedAt)
	u.UpdatedAt = pgxutil.FromTimestamptz(row.UpdatedAt)
	return nil
}

func mapUserRow(row sqlc.GetUserByIDRow) *domain.User {
	id, _ := pgxutil.FromPgUUID(row.ID)
	return &domain.User{
		ID:           id,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleID:       row.RoleID,
		RoleName:     row.RoleName,
		CreatedAt:    pgxutil.FromTimestamptz(row.CreatedAt),
		UpdatedAt:    pgxutil.FromTimestamptz(row.UpdatedAt),
	}
}

func mapUserEmailRow(row sqlc.GetUserByEmailRow) *domain.User {
	id, _ := pgxutil.FromPgUUID(row.ID)
	return &domain.User{
		ID:           id,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleID:       row.RoleID,
		RoleName:     row.RoleName,
		CreatedAt:    pgxutil.FromTimestamptz(row.CreatedAt),
		UpdatedAt:    pgxutil.FromTimestamptz(row.UpdatedAt),
	}
}
