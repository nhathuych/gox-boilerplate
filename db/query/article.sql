-- name: GetArticleByID :one
SELECT
  id,
  title,
  body,
  state,
  author_id,
  created_at,
  updated_at
FROM articles
WHERE
  id = $1;

-- name: ListArticles :many
SELECT
  id,
  title,
  body,
  state,
  author_id,
  created_at,
  updated_at
FROM articles
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListArticlesByAuthor :many
SELECT
  id,
  title,
  body,
  state,
  author_id,
  created_at,
  updated_at
FROM articles
WHERE
  author_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateArticle :one
INSERT INTO articles (title, body, state, author_id)
VALUES ($1, $2, $3, $4)
RETURNING
  id,
  title,
  body,
  state,
  author_id,
  created_at,
  updated_at;

-- name: UpdateArticle :one
UPDATE articles
SET
  title = $2,
  body = $3,
  state = $4,
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  id,
  title,
  body,
  state,
  author_id,
  created_at,
  updated_at;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = $1;
