package middleware

import (
	"net/http"

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
		// 最初に認証チェック
		authHandler := m.authMiddleware.RequireAuth()
		authHandler(c)
		
		// 認証処理で中断された場合はそのまま返す
		if c.IsAborted() {
			return
		}
		
		// ユーザーIDを取得
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
			c.Abort()
			return
		}
		
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーID形式が不正です"})
			c.Abort()
			return
		}
		
		// ユーザー情報を取得
		user, err := m.userRepo.FindByID(c.Request.Context(), userIDStr)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー情報が見つかりません"})
			c.Abort()
			return
		}
		
		// stユーザーかチェック
		if user.Name != "st" {
			c.JSON(http.StatusForbidden, gin.H{"error": "管理者権限が必要です"})
			c.Abort()
			return
		}
		
		c.Next()
	})
}