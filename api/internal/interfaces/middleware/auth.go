package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type User struct {
	ID      string
	IsAdmin bool
}

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	authService repository.AuthService
}

func NewAuthMiddleware(authService repository.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization token required",
			})
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			status := http.StatusUnauthorized
			errorCode := "invalid_token"
			message := "Invalid or expired token"

			if strings.Contains(err.Error(), "revoked") {
				errorCode = "token_revoked"
				message = "Token has been revoked"
			} else if strings.Contains(err.Error(), "expired") {
				errorCode = "token_expired"
				message = "Token has expired"
			}

			c.JSON(status, ErrorResponse{
				Error:   errorCode,
				Message: message,
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		
		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First require authentication
		token := extractTokenFromHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "Authorization token required",
			})
			c.Abort()
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_token",
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, userRole := range claims.Roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "insufficient_permissions",
				Message: "Insufficient permissions to access this resource",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		claims, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		
		c.Next()
	}
}

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

// Legacy functions for backward compatibility
func LegacyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func validateToken(token string) (*User, error) {
	switch token {
	case "admin-token":
		return &User{ID: "admin", IsAdmin: true}, nil
	case "user1-token":
		return &User{ID: "user1", IsAdmin: false}, nil
	case "user2-token":
		return &User{ID: "user2", IsAdmin: false}, nil
	default:
		return nil, http.ErrNotSupported
	}
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}