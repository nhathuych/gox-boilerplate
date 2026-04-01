package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/repository"
	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
)

func HasPermission(perms []string, need string) bool {
	for _, p := range perms {
		if p == need {
			return true
		}
	}
	return false
}

type ArticleUsecase struct {
	uow  *UnitOfWork
	repo domain.ArticleRepository
}

func NewArticleUsecase(uow *UnitOfWork, repo domain.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{uow: uow, repo: repo}
}

func (u *ArticleUsecase) Create(ctx context.Context, actorID uuid.UUID, perms []string, title, body string) (*domain.Article, error) {
	if !HasPermission(perms, "article:create") {
		return nil, domain.ErrForbidden
	}
	out := &domain.Article{
		Title:    title,
		Body:     body,
		State:    domain.ArticleStateDraft,
		AuthorID: actorID,
	}
	err := u.uow.WithTransaction(ctx, func(q *sqlc.Queries) error {
		repo := repository.NewArticleRepository(q)
		return repo.Create(ctx, out)
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (u *ArticleUsecase) GetByID(ctx context.Context, actorID uuid.UUID, perms []string, id uuid.UUID) (*domain.Article, error) {
	if !HasPermission(perms, "article:read") {
		return nil, domain.ErrForbidden
	}
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if HasPermission(perms, "article:delete") {
		return a, nil
	}
	if a.AuthorID == actorID {
		return a, nil
	}
	if a.State == domain.ArticleStatePublished {
		return a, nil
	}
	return nil, domain.ErrForbidden
}

func (u *ArticleUsecase) List(ctx context.Context, actorID uuid.UUID, perms []string, limit, offset int32, mineOnly bool) ([]domain.Article, error) {
	if !HasPermission(perms, "article:read") {
		return nil, domain.ErrForbidden
	}
	if mineOnly {
		return u.repo.ListByAuthor(ctx, actorID, limit, offset)
	}
	if !HasPermission(perms, "article:delete") {
		return u.repo.ListByAuthor(ctx, actorID, limit, offset)
	}
	return u.repo.List(ctx, limit, offset)
}

type UpdateArticleInput struct {
	Title string
	Body  string
	State domain.ArticleState
}

func (u *ArticleUsecase) Update(ctx context.Context, actorID uuid.UUID, perms []string, id uuid.UUID, in UpdateArticleInput) (*domain.Article, error) {
	if !HasPermission(perms, "article:update") {
		return nil, domain.ErrForbidden
	}
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	isAdmin := HasPermission(perms, "article:delete")
	if !isAdmin && a.AuthorID != actorID {
		return nil, domain.ErrForbidden
	}
	if !isAdmin {
		in.State = a.State
	} else if in.State == domain.ArticleStatePublished && !HasPermission(perms, "article:publish") {
		return nil, domain.ErrForbidden
	}
	a.Title = in.Title
	a.Body = in.Body
	a.State = in.State

	err = u.uow.WithTransaction(ctx, func(q *sqlc.Queries) error {
		repo := repository.NewArticleRepository(q)
		return repo.Update(ctx, a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (u *ArticleUsecase) Delete(ctx context.Context, actorID uuid.UUID, perms []string, id uuid.UUID) error {
	if !HasPermission(perms, "article:delete") {
		return domain.ErrForbidden
	}
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.uow.WithTransaction(ctx, func(q *sqlc.Queries) error {
		repo := repository.NewArticleRepository(q)
		return repo.Delete(ctx, id)
	})
}

func (u *ArticleUsecase) Publish(ctx context.Context, actorID uuid.UUID, perms []string, id uuid.UUID) (*domain.Article, error) {
	_ = actorID
	if !HasPermission(perms, "article:publish") {
		return nil, domain.ErrForbidden
	}
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.State = domain.ArticleStatePublished
	err = u.uow.WithTransaction(ctx, func(q *sqlc.Queries) error {
		repo := repository.NewArticleRepository(q)
		return repo.Update(ctx, a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}
