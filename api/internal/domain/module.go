package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/domain/service"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/fx"
)

var Module = fx.Module("domain",
	fx.Provide(
		service.NewUserService,
		fx.Annotate(
			service.NewAuthenticationService,
			fx.As(new(repository.AuthService)),
		),
		fx.Annotate(
			service.NewJWTTokenService,
			fx.As(new(repository.TokenService)),
		),
		provideJWTConfig,
	),
)

func provideJWTConfig(cfg *config.Config) (*entity.JWTConfig, error) {
	// Generate RSA keys for JWT signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &entity.JWTConfig{
		PrivateKey:           privateKey,
		PublicKey:            &privateKey.PublicKey,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               "utils-api",
	}, nil
}