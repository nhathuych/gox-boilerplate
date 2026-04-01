package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhathuych/gox-boilerplate/internal/auth"
	"github.com/nhathuych/gox-boilerplate/internal/config"
	"github.com/nhathuych/gox-boilerplate/internal/delivery/http/handler"
	"github.com/nhathuych/gox-boilerplate/internal/domain"
	"github.com/nhathuych/gox-boilerplate/internal/middleware"
	"github.com/nhathuych/gox-boilerplate/internal/repository"
	"github.com/nhathuych/gox-boilerplate/internal/repository/sqlc"
	"github.com/nhathuych/gox-boilerplate/internal/usecase"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"time"
)

func ProvidersModule() fx.Option {
	return fx.Module("providers",
		fx.Provide(
			newQuerier,
			newArticleRepository,
			newUserRepository,
			newRBACRepository,
			newJWTService,
			newTokenBlacklist,
			usecase.NewUnitOfWork,
			usecase.NewAuthUsecase,
			usecase.NewArticleUsecase,
			handler.NewAuthHandler,
			handler.NewArticleHandler,
			newAuthMiddleware,
			newRateLimiter,
		),
	)
}

func newQuerier(pool *pgxpool.Pool) sqlc.Querier {
	return sqlc.New(pool)
}

func newArticleRepository(q sqlc.Querier) domain.ArticleRepository {
	return repository.NewArticleRepository(q)
}

func newUserRepository(q sqlc.Querier) domain.UserRepository {
	return repository.NewUserRepository(q)
}

func newRBACRepository(q sqlc.Querier) domain.RBACRepository {
	return repository.NewRBACRepository(q)
}

func newJWTService(cfg *config.Config) *auth.JWTService {
	return auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
}

func newTokenBlacklist(cfg *config.Config, rdb *redis.Client) *auth.TokenBlacklist {
	return auth.NewTokenBlacklist(rdb, cfg.JWT.BlacklistTTL)
}

func newAuthMiddleware(jwt *auth.JWTService, bl *auth.TokenBlacklist) gin.HandlerFunc {
	return middleware.Auth(jwt, bl)
}

func newRateLimiter(rdb *redis.Client) middleware.RateLimiterHandler {
	// Simple default: 30 requests per minute per user/ip.
	return middleware.NewRedisRateLimiter(rdb, 30, time.Minute)
}
