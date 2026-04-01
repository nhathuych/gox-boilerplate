package pgxutil

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidUUID = errors.New("invalid uuid")

func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func FromPgUUID(p pgtype.UUID) (uuid.UUID, error) {
	if !p.Valid {
		return uuid.Nil, ErrInvalidUUID
	}
	return uuid.UUID(p.Bytes), nil
}
