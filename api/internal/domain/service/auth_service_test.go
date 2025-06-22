package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// Mock implementations
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) StoreSession(ctx context.Context, session *entity.SessionInfo) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepository) GetSession(ctx context.Context, jti string) (*entity.SessionInfo, error) {
	args := m.Called(ctx, jti)
	return args.Get(0).(*entity.SessionInfo), args.Error(1)
}

func (m *MockAuthRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*entity.SessionInfo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.SessionInfo), args.Error(1)
}

func (m *MockAuthRepository) UpdateSessionActivity(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteSession(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteSessionsByTokenFamily(ctx context.Context, tokenFamily string) error {
	args := m.Called(ctx, tokenFamily)
	return args.Error(0)
}

func (m *MockAuthRepository) AddToBlacklist(ctx context.Context, jti string, expiry time.Time) error {
	args := m.Called(ctx, jti, expiry)
	return args.Error(0)
}

func (m *MockAuthRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) CleanupExpiredBlacklist(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) IncrementRateLimit(ctx context.Context, key string, window time.Duration) error {
	args := m.Called(ctx, key, window)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthRepository) UpdateUserLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateAccessToken(user *entity.User, roles []string) (string, error) {
	args := m.Called(user, roles)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GenerateRefreshToken(userID, tokenFamily string) (string, error) {
	args := m.Called(userID, tokenFamily)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*entity.Claims), args.Error(1)
}

func (m *MockTokenService) ValidateRefreshToken(tokenString string) (*entity.RefreshClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*entity.RefreshClaims), args.Error(1)
}

func (m *MockTokenService) GenerateTokenPair(user *entity.User, roles []string, tokenFamily string) (*entity.TokenPair, error) {
	args := m.Called(user, roles, tokenFamily)
	return args.Get(0).(*entity.TokenPair), args.Error(1)
}

func (m *MockTokenService) RevokeToken(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockTokenService) RevokeTokenFamily(ctx context.Context, tokenFamily string) error {
	args := m.Called(ctx, tokenFamily)
	return args.Error(0)
}

func (m *MockTokenService) GetTokenClaims(tokenString string) (*entity.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*entity.Claims), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindWithOffsetPagination(ctx context.Context, limit, offset int, filter repository.PaginationFilter, sort repository.SortOption) (*repository.PaginationResult, error) {
	args := m.Called(ctx, limit, offset, filter, sort)
	return args.Get(0).(*repository.PaginationResult), args.Error(1)
}

func (m *MockUserRepository) FindWithCursorPagination(ctx context.Context, limit int, cursor string, filter repository.PaginationFilter, sort repository.SortOption) ([]*entity.User, error) {
	args := m.Called(ctx, limit, cursor, filter, sort)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func TestAuthenticationService_Login_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockAuthRepo := &MockAuthRepository{}
	mockTokenService := &MockTokenService{}
	mockUserRepo := &MockUserRepository{}

	authService := NewAuthenticationService(mockAuthRepo, mockTokenService, mockUserRepo)

	// Test data
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	user := &entity.User{
		ID:       "user-123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: string(hashedPassword),
	}

	loginReq := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	tokenPair := &entity.TokenPair{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	claims := &entity.Claims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	refreshClaims := &entity.RefreshClaims{
		UserID: "user-123",
		JTI:    "refresh-jti-123",
		Family: "family-123",
	}

	// Mock expectations
	mockAuthRepo.On("CheckRateLimit", ctx, "login_attempts:127.0.0.1", 5, 15*time.Minute).Return(true, nil)
	mockAuthRepo.On("GetUserByEmail", ctx, "test@example.com").Return(user, nil)
	mockTokenService.On("GenerateTokenPair", user, []string{"user"}, mock.AnythingOfType("string")).Return(tokenPair, nil)
	mockTokenService.On("GetTokenClaims", "access-token").Return(claims, nil)
	mockTokenService.On("ValidateRefreshToken", "refresh-token").Return(refreshClaims, nil)
	mockAuthRepo.On("StoreSession", ctx, mock.AnythingOfType("*entity.SessionInfo")).Return(nil)
	mockAuthRepo.On("UpdateUserLastLogin", ctx, "user-123").Return(nil)

	// Execute
	result, err := authService.Login(ctx, loginReq, "127.0.0.1", "test-agent")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)

	// Verify all expectations
	mockAuthRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAuthenticationService_Login_InvalidCredentials(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockAuthRepo := &MockAuthRepository{}
	mockTokenService := &MockTokenService{}
	mockUserRepo := &MockUserRepository{}

	authService := NewAuthenticationService(mockAuthRepo, mockTokenService, mockUserRepo)

	loginReq := &entity.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations
	mockAuthRepo.On("CheckRateLimit", ctx, "login_attempts:127.0.0.1", 5, 15*time.Minute).Return(true, nil)
	mockAuthRepo.On("GetUserByEmail", ctx, "nonexistent@example.com").Return((*entity.User)(nil), fmt.Errorf("user not found"))
	mockAuthRepo.On("IncrementRateLimit", ctx, "login_attempts:127.0.0.1", 15*time.Minute).Return(nil)

	// Execute
	result, err := authService.Login(ctx, loginReq, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credentials")

	// Verify expectations
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticationService_Login_RateLimited(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockAuthRepo := &MockAuthRepository{}
	mockTokenService := &MockTokenService{}
	mockUserRepo := &MockUserRepository{}

	authService := NewAuthenticationService(mockAuthRepo, mockTokenService, mockUserRepo)

	loginReq := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	// Mock expectations
	mockAuthRepo.On("CheckRateLimit", ctx, "login_attempts:127.0.0.1", 5, 15*time.Minute).Return(false, nil)

	// Execute
	result, err := authService.Login(ctx, loginReq, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "too many login attempts")

	// Verify expectations
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticationService_ValidateToken_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockAuthRepo := &MockAuthRepository{}
	mockTokenService := &MockTokenService{}
	mockUserRepo := &MockUserRepository{}

	authService := NewAuthenticationService(mockAuthRepo, mockTokenService, mockUserRepo)

	token := "valid-token"
	claims := &entity.Claims{
		UserID: "user-123",
		JTI:    "jti-123",
		Email:  "test@example.com",
	}

	// Mock expectations
	mockTokenService.On("ValidateAccessToken", token).Return(claims, nil)
	mockAuthRepo.On("IsBlacklisted", ctx, "jti-123").Return(false, nil)
	mockAuthRepo.On("UpdateSessionActivity", ctx, "jti-123").Return(nil)

	// Execute
	result, err := authService.ValidateToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-123", result.UserID)
	assert.Equal(t, "test@example.com", result.Email)

	// Verify expectations
	mockAuthRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestAuthenticationService_ValidateToken_Blacklisted(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockAuthRepo := &MockAuthRepository{}
	mockTokenService := &MockTokenService{}
	mockUserRepo := &MockUserRepository{}

	authService := NewAuthenticationService(mockAuthRepo, mockTokenService, mockUserRepo)

	token := "blacklisted-token"
	claims := &entity.Claims{
		UserID: "user-123",
		JTI:    "jti-123",
	}

	// Mock expectations
	mockTokenService.On("ValidateAccessToken", token).Return(claims, nil)
	mockAuthRepo.On("IsBlacklisted", ctx, "jti-123").Return(true, nil)

	// Execute
	result, err := authService.ValidateToken(ctx, token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token has been revoked")

	// Verify expectations
	mockAuthRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}