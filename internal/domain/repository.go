package domain

import (
	"context"

	"github.com/google/uuid"
)

type ArticleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Article, error)
	List(ctx context.Context, limit, offset int32) ([]Article, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int32) ([]Article, error)
	Create(ctx context.Context, a *Article) error
	Update(ctx context.Context, a *Article) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User, passwordHash string) error
}

type RBACRepository interface {
	ListPermissionsByRoleID(ctx context.Context, roleID int32) ([]string, error)
}
