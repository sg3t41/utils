package service

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sg3t41/api/internal/domain/entity"
)

func generateTestKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func TestJWTTokenService_GenerateAccessToken_Success(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		AccessTokenDuration: 15 * time.Minute,
		Issuer:              "test-issuer",
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	user := &entity.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	roles := []string{"user", "admin"}

	// Execute
	token, err := tokenService.GenerateAccessToken(user, roles)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the generated token
	claims, err := tokenService.ValidateAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.NotEmpty(t, claims.JTI)
}

func TestJWTTokenService_GenerateRefreshToken_Success(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               "test-issuer",
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	userID := "user-123"
	tokenFamily := "family-456"

	// Execute
	token, err := tokenService.GenerateRefreshToken(userID, tokenFamily)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the generated token
	claims, err := tokenService.ValidateRefreshToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "family-456", claims.Family)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.NotEmpty(t, claims.JTI)
}

func TestJWTTokenService_GenerateTokenPair_Success(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
		Issuer:               "test-issuer",
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	user := &entity.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	roles := []string{"user"}
	tokenFamily := "family-456"

	// Execute
	tokenPair, err := tokenService.GenerateTokenPair(user, roles, tokenFamily)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Equal(t, int64(900), tokenPair.ExpiresIn) // 15 minutes = 900 seconds

	// Validate access token
	accessClaims, err := tokenService.ValidateAccessToken(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", accessClaims.UserID)

	// Validate refresh token
	refreshClaims, err := tokenService.ValidateRefreshToken(tokenPair.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", refreshClaims.UserID)
	assert.Equal(t, "family-456", refreshClaims.Family)
}

func TestJWTTokenService_ValidateAccessToken_InvalidToken(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	// Execute
	claims, err := tokenService.ValidateAccessToken("invalid-token")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTTokenService_ValidateAccessToken_WrongSigningMethod(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	// Create a token with wrong signing method (HMAC instead of RSA)
	wrongToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	// Execute
	claims, err := tokenService.ValidateAccessToken(wrongToken)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

func TestJWTTokenService_GetTokenClaims_WithoutValidation(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		AccessTokenDuration: 15 * time.Minute,
		Issuer:              "test-issuer",
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	user := &entity.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	roles := []string{"user"}

	// Generate a token
	token, err := tokenService.GenerateAccessToken(user, roles)
	assert.NoError(t, err)

	// Execute - get claims without full validation (useful for expired tokens)
	claims, err := tokenService.GetTokenClaims(token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
}

func TestJWTTokenService_ValidateExpiredToken(t *testing.T) {
	// Setup
	privateKey, publicKey, err := generateTestKeys()
	assert.NoError(t, err)

	config := &entity.JWTConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		AccessTokenDuration: -1 * time.Hour, // Expired 1 hour ago
		Issuer:              "test-issuer",
	}

	mockAuthRepo := &MockAuthRepository{}
	tokenService := NewJWTTokenService(config, mockAuthRepo)

	user := &entity.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	roles := []string{"user"}

	// Generate an expired token
	token, err := tokenService.GenerateAccessToken(user, roles)
	assert.NoError(t, err)

	// Execute
	claims, err := tokenService.ValidateAccessToken(token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestEntity_Claims_Valid(t *testing.T) {
	// Test valid claims
	validClaims := &entity.Claims{
		UserID: "user-123",
		Email:  "test@example.com",
		JTI:    "jti-123",
	}

	err := validClaims.Valid()
	assert.NoError(t, err)

	// Test invalid claims - missing UserID
	invalidClaims := &entity.Claims{
		Email: "test@example.com",
		JTI:   "jti-123",
	}

	err = invalidClaims.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")

	// Test invalid claims - missing Email
	invalidClaims2 := &entity.Claims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	err = invalidClaims2.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")

	// Test invalid claims - missing JTI
	invalidClaims3 := &entity.Claims{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	err = invalidClaims3.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JTI is required")
}

func TestEntity_RefreshClaims_Valid(t *testing.T) {
	// Test valid refresh claims
	validClaims := &entity.RefreshClaims{
		UserID: "user-123",
		JTI:    "jti-123",
		Family: "family-456",
	}

	err := validClaims.Valid()
	assert.NoError(t, err)

	// Test invalid claims - missing Family
	invalidClaims := &entity.RefreshClaims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	err = invalidClaims.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token family is required")
}