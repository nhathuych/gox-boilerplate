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

type ArticleRepository struct {
	q sqlc.Querier
}

func NewArticleRepository(q sqlc.Querier) *ArticleRepository {
	return &ArticleRepository{q: q}
}

func (r *ArticleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Article, error) {
	row, err := r.q.GetArticleByID(ctx, pgxutil.ToPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return mapArticle(row), nil
}

func (r *ArticleRepository) List(ctx context.Context, limit, offset int32) ([]domain.Article, error) {
	rows, err := r.q.ListArticles(ctx, sqlc.ListArticlesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Article, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapArticle(row))
	}
	return out, nil
}

func (r *ArticleRepository) ListByAuthor(ctx context.Context, authorID uuid.UUID, limit, offset int32) ([]domain.Article, error) {
	rows, err := r.q.ListArticlesByAuthor(ctx, sqlc.ListArticlesByAuthorParams{
		AuthorID: pgxutil.ToPgUUID(authorID),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Article, 0, len(rows))
	for _, row := range rows {
		out = append(out, *mapArticle(row))
	}
	return out, nil
}

func (r *ArticleRepository) Create(ctx context.Context, a *domain.Article) error {
	row, err := r.q.CreateArticle(ctx, sqlc.CreateArticleParams{
		Title:    a.Title,
		Body:     a.Body,
		State:    string(a.State),
		AuthorID: pgxutil.ToPgUUID(a.AuthorID),
	})
	if err != nil {
		return err
	}
	mapped := mapArticle(row)
	*a = *mapped
	return nil
}

func (r *ArticleRepository) Update(ctx context.Context, a *domain.Article) error {
	row, err := r.q.UpdateArticle(ctx, sqlc.UpdateArticleParams{
		ID:    pgxutil.ToPgUUID(a.ID),
		Title: a.Title,
		Body:  a.Body,
		State: string(a.State),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	mapped := mapArticle(row)
	*a = *mapped
	return nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteArticle(ctx, pgxutil.ToPgUUID(id))
	if err != nil {
		return err
	}
	return nil
}

func mapArticle(row sqlc.Article) *domain.Article {
	id, _ := pgxutil.FromPgUUID(row.ID)
	author, _ := pgxutil.FromPgUUID(row.AuthorID)
	return &domain.Article{
		ID:        id,
		Title:     row.Title,
		Body:      row.Body,
		State:     domain.ArticleState(row.State),
		AuthorID:  author,
		CreatedAt: pgxutil.FromTimestamptz(row.CreatedAt),
		UpdatedAt: pgxutil.FromTimestamptz(row.UpdatedAt),
	}
}
