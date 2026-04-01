package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	claimTypeAccess  = "access"
	claimTypeRefresh = "refresh"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	AccessExp    time.Time
	RefreshExp   time.Time
}

type Claims struct {
	UserID      uuid.UUID `json:"uid"`
	RoleID      int32     `json:"rid"`
	Permissions []string  `json:"perms"`
	TokenType   string    `json:"typ"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTService(secret, issuer string, accessTTL, refreshTTL time.Duration) *JWTService {
	return &JWTService{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *JWTService) IssuePair(userID uuid.UUID, roleID int32, permissions []string) (*TokenPair, error) {
	jtiAccess := uuid.NewString()
	jtiRefresh := uuid.NewString()
	now := time.Now()

	accessClaims := Claims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: permissions,
		TokenType:   claimTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jtiAccess,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}
	refreshClaims := Claims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: permissions,
		TokenType:   claimTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jtiRefresh,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessStr, err := at.SignedString(s.secret)
	if err != nil {
		return nil, err
	}
	refreshStr, err := rt.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		AccessExp:    now.Add(s.accessTTL),
		RefreshExp:   now.Add(s.refreshTTL),
	}, nil
}

func (s *JWTService) ParseAccess(token string) (*Claims, error) {
	return s.parse(token, claimTypeAccess)
}

func (s *JWTService) ParseRefresh(token string) (*Claims, error) {
	return s.parse(token, claimTypeRefresh)
}

func (s *JWTService) parse(tokenString, wantType string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != wantType {
		return nil, errors.New("wrong token type")
	}
	return claims, nil
}
