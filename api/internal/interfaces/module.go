package interfaces

import (
	"github.com/sg3t41/api/internal/interfaces/handler"
	"github.com/sg3t41/api/internal/interfaces/middleware"
	"github.com/sg3t41/api/internal/interfaces/router"
	"go.uber.org/fx"
)

type Router interface {
	Run() error
}

var Module = fx.Module("interfaces",
	fx.Provide(
		handler.NewUserHandler,
		handler.NewAuthHandler,
		handler.NewArticleHandler,
		handler.NewUploadHandler,
		handler.NewLineHandler,
		fx.Annotate(
			handler.NewLineBotHandler,
			fx.As(new(handler.LineBotHandler)),
		),
		middleware.NewAuthMiddleware,
		middleware.NewAdminMiddleware,
		fx.Annotate(
			router.NewRouter,
			fx.As(new(Router)),
		),
	),
)