package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/middleware"
	"github.com/nhathuych/gox-boilerplate/internal/usecase"
)

type ArticleHandler struct {
	uc  *usecase.ArticleUsecase
	val *validator.Validate
}

func NewArticleHandler(uc *usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{uc: uc, val: validator.New()}
}

type createArticleReq struct {
	Title string `json:"title" validate:"required,min=1,max=500"`
	Body  string `json:"body" validate:"required"`
}

type updateArticleReq struct {
	Title string `json:"title" validate:"required,min=1,max=500"`
	Body  string `json:"body" validate:"required"`
	State string `json:"state" validate:"required,oneof=draft published"`
}

type articleResp struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toArticleResp(a *domain.Article) articleResp {
	return articleResp{
		ID:        a.ID.String(),
		Title:     a.Title,
		Body:      a.Body,
		State:     string(a.State),
		AuthorID:  a.AuthorID.String(),
		CreatedAt: a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt: a.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func actor(c *gin.Context) (uuid.UUID, []string, error) {
	uid, ok := middleware.UserIDFromContext(c.Request.Context())
	if !ok {
		return uuid.Nil, nil, errors.New("missing user")
	}
	perms, ok := middleware.PermissionsFromContext(c.Request.Context())
	if !ok {
		return uuid.Nil, nil, errors.New("missing permissions")
	}
	return uid, perms, nil
}

// Create godoc
// @Summary Create article (draft)
// @Tags articles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body createArticleReq true "payload"
// @Success 201 {object} articleResp
// @Router /articles [post]
func (h *ArticleHandler) Create(c *gin.Context) {
	var req createArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	a, err := h.uc.Create(c.Request.Context(), uid, perms, req.Title, req.Body)
	if err != nil {
		writeArticleErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, toArticleResp(a))
}

// Get godoc
// @Summary Get article by ID
// @Tags articles
// @Security BearerAuth
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} articleResp
// @Router /articles/{id} [get]
func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	a, err := h.uc.GetByID(c.Request.Context(), uid, perms, id)
	if err != nil {
		writeArticleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toArticleResp(a))
}

// List godoc
// @Summary List articles
// @Tags articles
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param mine query bool false "Only my articles"
// @Success 200 {array} articleResp
// @Router /articles [get]
func (h *ArticleHandler) List(c *gin.Context) {
	limit := int32(20)
	offset := int32(0)
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n > 0 && n <= 100 {
			limit = int32(n)
		}
	}
	if v := c.Query("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 0 {
			offset = int32(n)
		}
	}
	mine := c.Query("mine") == "true"
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	list, err := h.uc.List(c.Request.Context(), uid, perms, limit, offset, mine)
	if err != nil {
		writeArticleErr(c, err)
		return
	}
	out := make([]articleResp, 0, len(list))
	for i := range list {
		out = append(out, toArticleResp(&list[i]))
	}
	c.JSON(http.StatusOK, out)
}

// Update godoc
// @Summary Update article
// @Tags articles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param body body updateArticleReq true "payload"
// @Success 200 {object} articleResp
// @Router /articles/{id} [put]
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	a, err := h.uc.Update(c.Request.Context(), uid, perms, id, usecase.UpdateArticleInput{
		Title: req.Title,
		Body:  req.Body,
		State: domain.ArticleState(req.State),
	})
	if err != nil {
		writeArticleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toArticleResp(a))
}

// Delete godoc
// @Summary Delete article (admin)
// @Tags articles
// @Security BearerAuth
// @Param id path string true "Article ID"
// @Success 204
// @Router /articles/{id} [delete]
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.Delete(c.Request.Context(), uid, perms, id); err != nil {
		writeArticleErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Publish godoc
// @Summary Publish article (admin)
// @Tags articles
// @Security BearerAuth
// @Param id path string true "Article ID"
// @Success 200 {object} articleResp
// @Router /articles/{id}/publish [post]
func (h *ArticleHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	uid, perms, err := actor(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	a, err := h.uc.Publish(c.Request.Context(), uid, perms, id)
	if err != nil {
		writeArticleErr(c, err)
		return
	}
	c.JSON(http.StatusOK, toArticleResp(a))
}

func writeArticleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
