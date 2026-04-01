package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nhathuych/gox-boilerplate/internal/config"
	"github.com/nhathuych/gox-boilerplate/internal/delivery/http/handler"
	"github.com/nhathuych/gox-boilerplate/internal/middleware"
)

func NewRouter(
	cfg *config.Config,
	authMw gin.HandlerFunc,
	authH *handler.AuthHandler,
	articleH *handler.ArticleHandler,
	rateLimiterMw middleware.RateLimiterHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestTimeout(cfg.Server.RequestTimeout))
	r.Use(middleware.PrometheusMetrics())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(gin.HandlerFunc(rateLimiterMw))
	{
		api.POST("/auth/login", authH.Login)
		api.POST("/auth/register", authH.Register)
		api.POST("/auth/refresh", authH.Refresh)

		protected := api.Group("")
		protected.Use(authMw)
		{
			protected.POST("/auth/logout", authH.Logout)

			protected.POST("/articles", articleH.Create)
			protected.GET("/articles", articleH.List)
			protected.GET("/articles/:id", articleH.Get)
			protected.PUT("/articles/:id", articleH.Update)
			protected.DELETE("/articles/:id", middleware.RequirePermission("article:delete"), articleH.Delete)
			protected.POST("/articles/:id/publish", middleware.RequirePermission("article:publish"), articleH.Publish)
		}
	}

	return r
}
