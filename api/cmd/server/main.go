package main

import (
	"context"
	"database/sql"
	"github.com/sg3t41/api/internal/application"
	"github.com/sg3t41/api/internal/domain"
	"github.com/sg3t41/api/internal/infrastructure"
	"github.com/sg3t41/api/internal/interfaces"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		config.Module,
		domain.Module,
		application.Module,
		infrastructure.Module,
		interfaces.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(logger *zap.Logger, lifecycle fx.Lifecycle, router interfaces.Router, db *sql.DB, cfg *config.Config) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Starting application")
					
					// Database migrations are run manually
					logger.Info("Database migrations should be run manually")
					
					go router.Run()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Stopping application")
					return nil
				},
			})
		}),
	).Run()
}
