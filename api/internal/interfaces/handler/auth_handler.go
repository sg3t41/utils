package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type AuthHandler struct {
	authService repository.AuthService
}

func NewAuthHandler(authService repository.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Attempt login
	tokenPair, err := h.authService.Login(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		errorCode := "invalid_credentials"
		
		// Check for specific error types
		if strings.Contains(err.Error(), "too many") {
			status = http.StatusTooManyRequests
			errorCode = "rate_limit_exceeded"
		}

		c.JSON(status, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	// Extract user info from token
	claims, err := h.authService.ValidateToken(c.Request.Context(), tokenPair.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate valid token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: UserInfo{
			ID:    claims.UserID,
			Email: claims.Email,
			Name:  claims.Name,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req entity.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid request format",
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Refresh token
	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		errorCode := "invalid_token"
		
		if strings.Contains(err.Error(), "reuse detected") {
			errorCode = "token_reuse_detected"
		} else if strings.Contains(err.Error(), "session expired") {
			errorCode = "session_expired"
		}

		c.JSON(status, ErrorResponse{
			Error:   errorCode,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    tokenPair.TokenType,
		"expires_in":    tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract access token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Authorization header required",
		})
		return
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "invalid_token_format",
			Message: "Bearer token required",
		})
		return
	}

	accessToken := tokenParts[1]

	// Parse logout request (refresh token is optional)
	var req entity.LogoutRequest
	c.ShouldBindJSON(&req) // Don't error if body is empty

	// Perform logout
	if err := h.authService.Logout(c.Request.Context(), accessToken, &req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "logout_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *AuthHandler) RevokeAllSessions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	if err := h.authService.RevokeAllSessions(c.Request.Context(), userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "revocation_failed",
			Message: "Failed to revoke all sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All sessions revoked successfully",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user claims from context (set by auth middleware)
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	userClaims := claims.(*entity.Claims)
	
	c.JSON(http.StatusOK, UserInfo{
		ID:    userClaims.UserID,
		Email: userClaims.Email,
		Name:  userClaims.Name,
	})
}