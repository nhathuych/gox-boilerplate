package pgxutil

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromTimestamptz(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
