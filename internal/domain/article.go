package domain

import (
	"time"

	"github.com/google/uuid"
)

type ArticleState string

const (
	ArticleStateDraft     ArticleState = "draft"
	ArticleStatePublished ArticleState = "published"
)

type Article struct {
	ID        uuid.UUID
	Title     string
	Body      string
	State     ArticleState
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
