package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

// User ユーザー情報構造体
type User struct {
	ID      string
	IsAdmin bool
}

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware 認証ミドルウェア構造体
type AuthMiddleware struct {
	authService repository.AuthService
}

// NewAuthMiddleware 認証ミドルウェアの作成
func NewAuthMiddleware(authService repository.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 認証を必須とするミドルウェア
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			HandleUnauthorizedError(c, ErrMsgTokenRequired)
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			HandleAuthError(c, err)
			return
		}

		setUserContext(c, claims)
		c.Next()
	}
}

// RequireRole 特定のロールを必須とするミドルウェア
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			HandleUnauthorizedError(c, ErrMsgTokenRequired)
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			HandleAuthError(c, err)
			return
		}

		if !hasRole(claims.Roles, role) {
			HandleForbiddenError(c, ErrMsgInsufficientPerms)
			return
		}

		setUserContext(c, claims)
		c.Next()
	}
}

// OptionalAuth オプションの認証ミドルウェア
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}

		setUserContext(c, claims)
		c.Next()
	}
}

// extractTokenFromHeader リクエストヘッダーからトークンを抽出
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ""
	}

	return tokenParts[1]
}

// setUserContext ユーザー情報をコンテキストに設定
func setUserContext(c *gin.Context, claims *entity.Claims) {
	c.Set("user", claims)
	c.Set("user_id", claims.UserID)
	c.Set("user_email", claims.Email)
	c.Set("user_roles", claims.Roles)
}

// hasRole ユーザーが特定のロールを持っているかチェック
func hasRole(userRoles []string, requiredRole string) bool {
	for _, role := range userRoles {
		if role == requiredRole {
			return true
		}
	}
	return false
}


func GetUserFromContext(c *gin.Context) *User {
	if user, exists := c.Get("user"); exists {
		if claims, ok := user.(*entity.Claims); ok {
			return &User{
				ID:      claims.UserID,
				IsAdmin: contains(claims.Roles, "admin"),
			}
		}
		return user.(*User)
	}
	return nil
}

func UserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(UserContextKey).(*User); ok {
		return user
	}
	return nil
}

// contains スライス内に特定の要素が含まれているかチェック
func contains(slice []string, item string) bool {
	return hasRole(slice, item)
}