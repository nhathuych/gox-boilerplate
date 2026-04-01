package repository

import (
	"context"

	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
)

type RBACRepository struct {
	q sqlc.Querier
}

func NewRBACRepository(q sqlc.Querier) *RBACRepository {
	return &RBACRepository{q: q}
}

func (r *RBACRepository) ListPermissionsByRoleID(ctx context.Context, roleID int32) ([]string, error) {
	return r.q.ListPermissionsByRoleID(ctx, roleID)
}
