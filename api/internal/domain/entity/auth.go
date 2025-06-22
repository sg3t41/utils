package entity

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string   `json:"sub"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	JTI    string   `json:"jti"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"sub"`
	JTI    string `json:"jti"`
	Family string `json:"family"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type JWTConfig struct {
	PrivateKey           *rsa.PrivateKey
	PublicKey            *rsa.PublicKey
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
}

func NewTokenPair(accessToken, refreshToken string, expiresIn time.Duration) *TokenPair {
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(expiresIn.Seconds()),
	}
}

func NewClaims(user *User, roles []string, duration time.Duration) *Claims {
	now := time.Now()
	return &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Roles:  roles,
		JTI:    uuid.New().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

func NewRefreshClaims(userID string, family string, duration time.Duration) *RefreshClaims {
	now := time.Now()
	return &RefreshClaims{
		UserID: userID,
		JTI:    uuid.New().String(),
		Family: family,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

func (c *Claims) Valid() error {
	if c.UserID == "" {
		return errors.New("user ID is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.JTI == "" {
		return errors.New("JTI is required")
	}
	return nil
}

func (rc *RefreshClaims) Valid() error {
	if rc.UserID == "" {
		return errors.New("user ID is required")
	}
	if rc.JTI == "" {
		return errors.New("JTI is required")
	}
	if rc.Family == "" {
		return errors.New("token family is required")
	}
	return nil
}

type SessionInfo struct {
	UserID       string    `json:"user_id"`
	TokenFamily  string    `json:"token_family"`
	JTI          string    `json:"jti"`
	RefreshJTI   string    `json:"refresh_jti"`
	IssuedAt     time.Time `json:"issued_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActivity time.Time `json:"last_activity"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

func NewSessionInfo(userID, tokenFamily, jti, refreshJTI string, expiresAt time.Time, ipAddress, userAgent string) *SessionInfo {
	now := time.Now()
	return &SessionInfo{
		UserID:       userID,
		TokenFamily:  tokenFamily,
		JTI:          jti,
		RefreshJTI:   refreshJTI,
		IssuedAt:     now,
		ExpiresAt:    expiresAt,
		LastActivity: now,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}
}

func (s *SessionInfo) UpdateActivity() {
	s.LastActivity = time.Now()
}

func (s *SessionInfo) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}