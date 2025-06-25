package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/repository"
)

// AdminMiddleware st ユーザーのみアクセス可能な管理者ミドルウェア
type AdminMiddleware struct {
	authMiddleware *AuthMiddleware
	userRepo       repository.UserRepository
}

// NewAdminMiddleware creates a new AdminMiddleware
func NewAdminMiddleware(authMiddleware *AuthMiddleware, userRepo repository.UserRepository) *AdminMiddleware {
	return &AdminMiddleware{
		authMiddleware: authMiddleware,
		userRepo:       userRepo,
	}
}

// RequireAdmin stユーザーのみを許可するミドルウェア
func (m *AdminMiddleware) RequireAdmin() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHandler := m.authMiddleware.RequireAuth()
		authHandler(c)
		
		if c.IsAborted() {
			return
		}
		
		if !m.isAdminUser(c) {
			return
		}
		
		c.Next()
	})
}

// isAdminUser ユーザーが管理者かチェック
func (m *AdminMiddleware) isAdminUser(c *gin.Context) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		HandleUnauthorizedError(c, ErrMsgAuthRequired)
		return false
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		HandleUnauthorizedError(c, ErrMsgInvalidUserID)
		return false
	}
	
	user, err := m.userRepo.FindByID(c.Request.Context(), userIDStr)
	if err != nil || user == nil {
		HandleUnauthorizedError(c, ErrMsgUserNotFound)
		return false
	}
	
	if user.Name != "st" {
		HandleForbiddenError(c, ErrMsgAdminRequired)
		return false
	}
	
	return true
}