package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClaims_Valid(t *testing.T) {
	// Test valid claims
	validClaims := &Claims{
		UserID: "user-123",
		Email:  "test@example.com",
		JTI:    "jti-123",
	}

	err := validClaims.Valid()
	assert.NoError(t, err)

	// Test invalid claims - missing UserID
	invalidClaims := &Claims{
		Email: "test@example.com",
		JTI:   "jti-123",
	}

	err = invalidClaims.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID is required")

	// Test invalid claims - missing Email
	invalidClaims2 := &Claims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	err = invalidClaims2.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is required")

	// Test invalid claims - missing JTI
	invalidClaims3 := &Claims{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	err = invalidClaims3.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JTI is required")
}

func TestRefreshClaims_Valid(t *testing.T) {
	// Test valid refresh claims
	validClaims := &RefreshClaims{
		UserID: "user-123",
		JTI:    "jti-123",
		Family: "family-456",
	}

	err := validClaims.Valid()
	assert.NoError(t, err)

	// Test invalid claims - missing Family
	invalidClaims := &RefreshClaims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	err = invalidClaims.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token family is required")
}

func TestNewTokenPair(t *testing.T) {
	accessToken := "access-token-123"
	refreshToken := "refresh-token-456"
	duration := 15 * time.Minute

	tokenPair := NewTokenPair(accessToken, refreshToken, duration)

	assert.Equal(t, accessToken, tokenPair.AccessToken)
	assert.Equal(t, refreshToken, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Equal(t, int64(900), tokenPair.ExpiresIn) // 15 minutes = 900 seconds
}

func TestNewClaims(t *testing.T) {
	user := &User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	roles := []string{"user", "admin"}
	duration := 15 * time.Minute

	claims := NewClaims(user, roles, duration)

	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.Equal(t, roles, claims.Roles)
	assert.NotEmpty(t, claims.JTI)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.NotBefore)
}

func TestNewRefreshClaims(t *testing.T) {
	userID := "user-123"
	family := "family-456"
	duration := 7 * 24 * time.Hour

	claims := NewRefreshClaims(userID, family, duration)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, family, claims.Family)
	assert.NotEmpty(t, claims.JTI)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.NotBefore)
}

func TestSessionInfo(t *testing.T) {
	userID := "user-123"
	tokenFamily := "family-456"
	jti := "jti-789"
	refreshJTI := "refresh-jti-012"
	expiresAt := time.Now().Add(1 * time.Hour)
	ipAddress := "192.168.1.1"
	userAgent := "test-agent"

	session := NewSessionInfo(userID, tokenFamily, jti, refreshJTI, expiresAt, ipAddress, userAgent)

	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, tokenFamily, session.TokenFamily)
	assert.Equal(t, jti, session.JTI)
	assert.Equal(t, refreshJTI, session.RefreshJTI)
	assert.Equal(t, expiresAt, session.ExpiresAt)
	assert.Equal(t, ipAddress, session.IPAddress)
	assert.Equal(t, userAgent, session.UserAgent)
	assert.NotNil(t, session.IssuedAt)
	assert.NotNil(t, session.LastActivity)
}

func TestSessionInfo_UpdateActivity(t *testing.T) {
	session := &SessionInfo{
		LastActivity: time.Now().Add(-1 * time.Hour),
	}

	oldActivity := session.LastActivity
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	session.UpdateActivity()

	assert.True(t, session.LastActivity.After(oldActivity))
}

func TestSessionInfo_IsExpired(t *testing.T) {
	// Test non-expired session
	futureSession := &SessionInfo{
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(t, futureSession.IsExpired())

	// Test expired session
	expiredSession := &SessionInfo{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(t, expiredSession.IsExpired())
}