package usecase

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nhathuych/gox-boilerplate/internal/auth"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
)

type AuthUsecase struct {
	users     domain.UserRepository
	rbac      domain.RBACRepository
	jwt       *auth.JWTService
	blacklist *auth.TokenBlacklist
}

func NewAuthUsecase(
	users domain.UserRepository,
	rbac domain.RBACRepository,
	jwt *auth.JWTService,
	blacklist *auth.TokenBlacklist,
) *AuthUsecase {
	return &AuthUsecase{
		users:     users,
		rbac:      rbac,
		jwt:       jwt,
		blacklist: blacklist,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (u *AuthUsecase) Login(ctx context.Context, in LoginInput) (*auth.TokenPair, error) {
	user, err := u.users.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidPassword
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, domain.ErrInvalidPassword
	}
	perms, err := u.rbac.ListPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}
	return u.jwt.IssuePair(user.ID, user.RoleID, perms)
}

type RegisterInput struct {
	Email    string
	Password string
}

const roleUserID = int32(2)

func (u *AuthUsecase) Register(ctx context.Context, in RegisterInput) (*auth.TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	du := &domain.User{
		Email:  in.Email,
		RoleID: roleUserID,
	}
	if err := u.users.Create(ctx, du, string(hash)); err != nil {
		return nil, err
	}
	full, err := u.users.GetByID(ctx, du.ID)
	if err != nil {
		return nil, err
	}
	perms, err := u.rbac.ListPermissionsByRoleID(ctx, full.RoleID)
	if err != nil {
		return nil, err
	}
	return u.jwt.IssuePair(full.ID, full.RoleID, perms)
}

func (u *AuthUsecase) Refresh(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	claims, err := u.jwt.ParseRefresh(refreshToken)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	user, err := u.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	perms, err := u.rbac.ListPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}
	return u.jwt.IssuePair(user.ID, user.RoleID, perms)
}

func (u *AuthUsecase) LogoutAccess(ctx context.Context, jti string, exp time.Time) error {
	if jti == "" {
		return nil
	}
	return u.blacklist.Add(ctx, jti, exp)
}
