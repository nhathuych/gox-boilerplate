package middleware

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const (
	userIDKey ctxKey = "userID"
	permsKey  ctxKey = "permissions"
	jtiKey    ctxKey = "jti"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	return v, ok
}

func PermissionsFromContext(ctx context.Context) ([]string, bool) {
	v, ok := ctx.Value(permsKey).([]string)
	return v, ok
}

func JTIFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(jtiKey).(string)
	return v, ok
}

func WithUser(ctx context.Context, id uuid.UUID, perms []string, jti string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, id)
	ctx = context.WithValue(ctx, permsKey, perms)
	ctx = context.WithValue(ctx, jtiKey, jti)
	return ctx
}
