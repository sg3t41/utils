package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type JWTTokenService struct {
	config     *entity.JWTConfig
	authRepo   repository.AuthRepository
}

func NewJWTTokenService(config *entity.JWTConfig, authRepo repository.AuthRepository) *JWTTokenService {
	return &JWTTokenService{
		config:   config,
		authRepo: authRepo,
	}
}

func (s *JWTTokenService) GenerateAccessToken(user *entity.User, roles []string) (string, error) {
	claims := entity.NewClaims(user, roles, s.config.AccessTokenDuration)
	claims.Issuer = s.config.Issuer
	
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.config.PrivateKey)
}

func (s *JWTTokenService) GenerateRefreshToken(userID, tokenFamily string) (string, error) {
	claims := entity.NewRefreshClaims(userID, tokenFamily, s.config.RefreshTokenDuration)
	claims.Issuer = s.config.Issuer
	
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.config.PrivateKey)
}

func (s *JWTTokenService) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.config.PublicKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*entity.Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

func (s *JWTTokenService) ValidateRefreshToken(tokenString string) (*entity.RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.config.PublicKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*entity.RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid refresh token")
}

func (s *JWTTokenService) GenerateTokenPair(user *entity.User, roles []string, tokenFamily string) (*entity.TokenPair, error) {
	if tokenFamily == "" {
		tokenFamily = uuid.New().String()
	}
	
	accessToken, err := s.GenerateAccessToken(user, roles)
	if err != nil {
		return nil, err
	}
	
	refreshToken, err := s.GenerateRefreshToken(user.ID, tokenFamily)
	if err != nil {
		return nil, err
	}
	
	return entity.NewTokenPair(accessToken, refreshToken, s.config.AccessTokenDuration), nil
}

func (s *JWTTokenService) RevokeToken(ctx context.Context, jti string) error {
	expiry := time.Now().Add(s.config.AccessTokenDuration)
	return s.authRepo.AddToBlacklist(ctx, jti, expiry)
}

func (s *JWTTokenService) RevokeTokenFamily(ctx context.Context, tokenFamily string) error {
	return s.authRepo.DeleteSessionsByTokenFamily(ctx, tokenFamily)
}

func (s *JWTTokenService) GetTokenClaims(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.config.PublicKey, nil
	}, jwt.WithoutClaimsValidation())
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*entity.Claims); ok {
		return claims, nil
	}
	
	return nil, errors.New("invalid token format")
}

func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// This would typically load from secure storage in production
	// For now, we'll use a dummy implementation
	return nil, nil, errors.New("key pair generation not implemented - should load from secure storage")
}