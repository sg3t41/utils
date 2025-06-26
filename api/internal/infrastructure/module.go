package infrastructure

import (
	"database/sql"

	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/infrastructure/persistence"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("infrastructure",
	fx.Provide(
		provideDB,
		provideUserRepository,
		provideAuthRepository,
		provideArticleRepository,
		provideLinkRepository,
	),
)

func provideDB(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	if cfg.UseMemoryDB {
		return nil, nil
	}
	return persistence.NewPostgresDB(cfg, logger)
}

func provideUserRepository(cfg *config.Config, db *sql.DB) repository.UserRepository {
	if cfg.UseMemoryDB {
		return persistence.NewMemoryUserRepository()
	}
	return persistence.NewPostgresUserRepository(db)
}

func provideAuthRepository(cfg *config.Config, userRepo repository.UserRepository) repository.AuthRepository {
	return persistence.NewMemoryAuthRepository(userRepo)
}

func provideArticleRepository(cfg *config.Config, db *sql.DB) repository.ArticleRepository {
	if cfg.UseMemoryDB {
		// TODO: Implement memory article repository if needed
		panic("Memory article repository not implemented")
	}
	return persistence.NewPostgresArticleRepository(db)
}

func provideLinkRepository(cfg *config.Config, db *sql.DB) repository.LinkRepository {
	if cfg.UseMemoryDB {
		// TODO: Implement memory link repository if needed
		panic("Memory link repository not implemented")
	}
	return persistence.NewPostgresLinkRepository(persistence.NewSqlxDB(db))
}