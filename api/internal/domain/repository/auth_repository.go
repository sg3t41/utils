package repository

import (
	"context"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
)

type AuthRepository interface {
	// セッション管理
	StoreSession(ctx context.Context, session *entity.SessionInfo) error
	GetSession(ctx context.Context, jti string) (*entity.SessionInfo, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]*entity.SessionInfo, error)
	UpdateSessionActivity(ctx context.Context, jti string) error
	DeleteSession(ctx context.Context, jti string) error
	DeleteSessionsByUserID(ctx context.Context, userID string) error
	DeleteSessionsByTokenFamily(ctx context.Context, tokenFamily string) error
	
	// トークンブラックリスト
	AddToBlacklist(ctx context.Context, jti string, expiry time.Time) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
	CleanupExpiredBlacklist(ctx context.Context) error
	
	// レート制限
	CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
	IncrementRateLimit(ctx context.Context, key string, window time.Duration) error
	
	// ユーザー認証
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUserLastLogin(ctx context.Context, userID string) error
}

type TokenService interface {
	GenerateAccessToken(user *entity.User, roles []string) (string, error)
	GenerateRefreshToken(userID, tokenFamily string) (string, error)
	ValidateAccessToken(tokenString string) (*entity.Claims, error)
	ValidateRefreshToken(tokenString string) (*entity.RefreshClaims, error)
	GenerateTokenPair(user *entity.User, roles []string, tokenFamily string) (*entity.TokenPair, error)
	RevokeToken(ctx context.Context, jti string) error
	RevokeTokenFamily(ctx context.Context, tokenFamily string) error
	GetTokenClaims(tokenString string) (*entity.Claims, error)
}

type AuthService interface {
	Login(ctx context.Context, req *entity.LoginRequest, ipAddress, userAgent string) (*entity.TokenPair, error)
	RefreshToken(ctx context.Context, req *entity.RefreshRequest, ipAddress, userAgent string) (*entity.TokenPair, error)
	Logout(ctx context.Context, accessToken string, req *entity.LogoutRequest) error
	ValidateToken(ctx context.Context, tokenString string) (*entity.Claims, error)
	RevokeAllSessions(ctx context.Context, userID string) error
}