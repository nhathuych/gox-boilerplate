package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockArticleRepo struct {
	mock.Mock
}

func (m *mockArticleRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *mockArticleRepo) List(ctx context.Context, limit, offset int32) ([]domain.Article, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Article), args.Error(1)
}

func (m *mockArticleRepo) ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int32) ([]domain.Article, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Article), args.Error(1)
}

func (m *mockArticleRepo) Create(ctx context.Context, a *domain.Article) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *mockArticleRepo) Update(ctx context.Context, a *domain.Article) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *mockArticleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestArticleUsecase_Delete_ForbidWhenMissingPermission(t *testing.T) {
	repo := &mockArticleRepo{}

	// uow is nil because we expect permission to fail before any transaction/repo call.
	uc := usecase.NewArticleUsecase((*usecase.UnitOfWork)(nil), repo)

	err := uc.Delete(context.Background(), uuid.New(), []string{}, uuid.New())
	require.ErrorIs(t, err, domain.ErrForbidden)

	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// Compile-time guard: mock implements the repository interface.
var _ domain.ArticleRepository = (*mockArticleRepo)(nil)
