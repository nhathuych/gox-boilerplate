package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/repository"
	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
	"github.com/nhathuych/gox-boilerplate/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndGetByEmail(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		// Default docker-compose DSN.
		dsn = "postgres://postgres:postgres@localhost:5432/gox?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Migrate schema for integration test.
	migrationsDir := filepath.Join("..", "db", "migration")
	require.NoError(t, testutil.MigrateUp(dsn, migrationsDir))

	pool, err := testutil.ConnectPool(ctx, dsn)
	require.NoError(t, err)
	defer pool.Close()

	q := sqlc.New(pool)
	repo := repository.NewUserRepository(q)

	u := &domain.User{
		Email:  "it-user-" + uuid.NewString() + "@example.com",
		RoleID: 2, // seeded "user"
	}

	require.NoError(t, repo.Create(ctx, u, "test-password-hash"))

	got, err := repo.GetByEmail(ctx, u.Email)
	require.NoError(t, err)
	require.Equal(t, u.Email, got.Email)
	require.Equal(t, int32(2), got.RoleID)
	require.Equal(t, "user", got.RoleName)
	require.Equal(t, "test-password-hash", got.PasswordHash)
}
