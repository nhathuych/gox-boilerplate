package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/nhathuych/gox-boilerplate/internal/auth"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/usecase"
)

type AuthHandler struct {
	authUC *usecase.AuthUsecase
	jwt    *auth.JWTService
	val    *validator.Validate
}

func NewAuthHandler(authUC *usecase.AuthUsecase, jwt *auth.JWTService) *AuthHandler {
	return &AuthHandler{authUC: authUC, jwt: jwt, val: validator.New()}
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type registerReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginReq true "credentials"
// @Success 200 {object} tokenResp
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, err := h.authUC.Login(c.Request.Context(), usecase.LoginInput{Email: req.Email, Password: req.Password})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	c.JSON(http.StatusOK, tokenResp{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.AccessExp.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

// Register godoc
// @Summary Register
// @Tags auth
// @Accept json
// @Produce json
// @Param body body registerReq true "payload"
// @Success 201 {object} tokenResp
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, err := h.authUC.Register(c.Request.Context(), usecase.RegisterInput{Email: req.Email, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "could not register"})
		return
	}
	c.JSON(http.StatusCreated, tokenResp{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.AccessExp.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

// Refresh godoc
// @Summary Refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body refreshReq true "payload"
// @Success 200 {object} tokenResp
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.val.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pair, err := h.authUC.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, tokenResp{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresAt:    pair.AccessExp.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

// Logout godoc
// @Summary Logout (blacklist access token)
// @Tags auth
// @Security BearerAuth
// @Success 204
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	hd := c.GetHeader("Authorization")
	if hd == "" || !strings.HasPrefix(hd, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
		return
	}
	raw := strings.TrimPrefix(strings.TrimSpace(hd), "Bearer ")
	claims, err := h.jwt.ParseAccess(raw)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	if claims.ExpiresAt == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token claims"})
		return
	}
	if err := h.authUC.LogoutAccess(c.Request.Context(), claims.ID, claims.ExpiresAt.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}
	c.Status(http.StatusNoContent)
}
