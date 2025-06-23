package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/interfaces/dto"
	"github.com/sg3t41/api/internal/interfaces/handler"
	"github.com/sg3t41/api/internal/interfaces/middleware"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/zap"
)

type Router struct {
	engine             *gin.Engine
	logger             *zap.Logger
	config             *config.Config
	userHandler        *handler.UserHandler
	authHandler        *handler.AuthHandler
	articleHandler     *handler.ArticleHandler
	uploadHandler      *handler.UploadHandler
	lineHandler        *handler.LineHandler
	authMiddleware     *middleware.AuthMiddleware
	validationMiddleware *middleware.ValidationMiddleware
}

func NewRouter(
	logger *zap.Logger,
	config *config.Config,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	articleHandler *handler.ArticleHandler,
	uploadHandler *handler.UploadHandler,
	lineHandler *handler.LineHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	gin.SetMode(config.GinMode)
	engine := gin.New()

	engine.Use(middleware.Logger(logger))
	engine.Use(middleware.Recovery(logger))
	engine.Use(middleware.CORS())

	validationMiddleware := middleware.NewValidationMiddleware()

	return &Router{
		engine:             engine,
		logger:             logger,
		config:             config,
		userHandler:        userHandler,
		authHandler:        authHandler,
		articleHandler:     articleHandler,
		uploadHandler:      uploadHandler,
		lineHandler:        lineHandler,
		authMiddleware:     authMiddleware,
		validationMiddleware: validationMiddleware,
	}
}

func (r *Router) SetupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "utils-api",
		})
	})

	v1 := r.engine.Group("/api/v1")
	{
		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/logout", r.authHandler.Logout)
			
			// LINE認証エンドポイント
			line := auth.Group("/line")
			{
				line.GET("/url", r.lineHandler.GetAuthURL)
				line.POST("/callback", r.lineHandler.CallbackPost)
			}
			
			// Protected auth endpoints
			authProtected := auth.Use(r.authMiddleware.RequireAuth())
			{
				authProtected.GET("/profile", r.authHandler.GetProfile)
				authProtected.POST("/revoke-all", r.authHandler.RevokeAllSessions)
			}
		}

		// User endpoints
		users := v1.Group("/users")
		{
			users.GET("", r.validationMiddleware.ValidateQuery(&dto.ListUsersQuery{}), r.userHandler.GetUsers)
			users.POST("", r.validationMiddleware.ValidateJSON(&dto.CreateUserRequest{}), r.userHandler.CreateUser)
			users.GET("/:id", r.validationMiddleware.ValidateQuery(&dto.GetUserQuery{}), r.userHandler.GetUser)
			users.DELETE("/:id", r.userHandler.DeleteUser) // 開発用: 認証不要で削除可能
			
			// Protected user endpoints
			authenticated := users.Use(r.authMiddleware.RequireAuth())
			{
				authenticated.PATCH("/:id", r.validationMiddleware.ValidateJSON(&dto.UpdateUserRequest{}), r.userHandler.UpdateUser)
				authenticated.PATCH("/:id/password", r.validationMiddleware.ValidateJSON(&dto.UpdatePasswordRequest{}), r.userHandler.UpdatePassword)
			}
		}

		// Article endpoints
		articles := v1.Group("/articles")
		{
			// Public endpoints
			articles.GET("", r.validationMiddleware.ValidateQuery(&dto.ListArticlesQuery{}), r.articleHandler.GetArticles)
			articles.GET("/:id", r.validationMiddleware.ValidateQuery(&dto.GetArticleQuery{}), r.articleHandler.GetArticle)
			
			// Development: Temporarily allow CRUD without auth for testing
			articles.POST("", r.validationMiddleware.ValidateJSON(&dto.CreateArticleRequest{}), r.articleHandler.CreateArticle)
			articles.PUT("/:id", r.validationMiddleware.ValidateJSON(&dto.UpdateArticleRequest{}), r.articleHandler.UpdateArticle)
			articles.DELETE("/:id", r.articleHandler.DeleteArticle)
			articles.POST("/:id/publish", r.validationMiddleware.ValidateJSON(&dto.PublishArticleRequest{}), r.articleHandler.PublishArticle)
			articles.POST("/:id/unpublish", r.articleHandler.UnpublishArticle)
			
			// Protected article endpoints (commented out for development)
			// articlesAuth := articles.Use(r.authMiddleware.RequireAuth())
			// {
			// 	articlesAuth.PUT("/:id", r.validationMiddleware.ValidateJSON(&dto.UpdateArticleRequest{}), r.articleHandler.UpdateArticle)
			// 	articlesAuth.DELETE("/:id", r.articleHandler.DeleteArticle)
			// 	articlesAuth.POST("/:id/publish", r.validationMiddleware.ValidateJSON(&dto.PublishArticleRequest{}), r.articleHandler.PublishArticle)
			// 	articlesAuth.POST("/:id/unpublish", r.articleHandler.UnpublishArticle)
			// }
		}

		// Upload endpoints
		upload := v1.Group("/upload")
		{
			// Development: Temporarily allow upload without auth for testing
			upload.POST("/image", r.uploadHandler.UploadImage)
			upload.DELETE("/image", r.uploadHandler.DeleteImage)
		}

		// Static file serving for uploaded images
		v1.GET("/uploads/*path", r.uploadHandler.ServeImage)
	}
}

func (r *Router) Run() error {
	r.SetupRoutes()
	r.logger.Info("Starting server", zap.String("addr", r.config.ServerAddress))
	return r.engine.Run(r.config.ServerAddress)
}