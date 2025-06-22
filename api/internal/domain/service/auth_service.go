package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type AuthenticationService struct {
	authRepo     repository.AuthRepository
	tokenService repository.TokenService
	userRepo     repository.UserRepository
}

func NewAuthenticationService(
	authRepo repository.AuthRepository,
	tokenService repository.TokenService,
	userRepo repository.UserRepository,
) *AuthenticationService {
	return &AuthenticationService{
		authRepo:     authRepo,
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}

func (s *AuthenticationService) Login(ctx context.Context, req *entity.LoginRequest, ipAddress, userAgent string) (*entity.TokenPair, error) {
	// Rate limiting check
	rateLimitKey := fmt.Sprintf("login_attempts:%s", ipAddress)
	allowed, err := s.authRepo.CheckRateLimit(ctx, rateLimitKey, 5, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("rate limit check failed: %w", err)
	}
	if !allowed {
		return nil, errors.New("too many login attempts, please try again later")
	}

	// Get user by email
	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Increment rate limit counter on failed attempt
		s.authRepo.IncrementRateLimit(ctx, rateLimitKey, 15*time.Minute)
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Increment rate limit counter on failed attempt
		s.authRepo.IncrementRateLimit(ctx, rateLimitKey, 15*time.Minute)
		return nil, errors.New("invalid credentials")
	}

	// Generate token family for this session
	tokenFamily := uuid.New().String()
	
	// Generate token pair
	roles := []string{"user"} // Default role, could be fetched from user or database
	tokenPair, err := s.tokenService.GenerateTokenPair(user, roles, tokenFamily)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Extract claims to get JTIs
	accessClaims, err := s.tokenService.GetTokenClaims(tokenPair.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract access token claims: %w", err)
	}

	refreshClaims, err := s.tokenService.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract refresh token claims: %w", err)
	}

	// Store session information
	session := entity.NewSessionInfo(
		user.ID,
		tokenFamily,
		accessClaims.JTI,
		refreshClaims.JTI,
		refreshClaims.ExpiresAt.Time,
		ipAddress,
		userAgent,
	)

	if err := s.authRepo.StoreSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Update user's last login time
	if err := s.authRepo.UpdateUserLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Failed to update last login for user %s: %v\n", user.ID, err)
	}

	return tokenPair, nil
}

func (s *AuthenticationService) RefreshToken(ctx context.Context, req *entity.RefreshRequest, ipAddress, userAgent string) (*entity.TokenPair, error) {
	// Validate refresh token
	refreshClaims, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if refresh token is blacklisted
	isBlacklisted, err := s.authRepo.IsBlacklisted(ctx, refreshClaims.JTI)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if isBlacklisted {
		// Token reuse detected - revoke entire token family
		s.tokenService.RevokeTokenFamily(ctx, refreshClaims.Family)
		return nil, errors.New("token reuse detected - all sessions revoked")
	}

	// Get session info
	session, err := s.authRepo.GetSession(ctx, refreshClaims.JTI)
	if err != nil {
		// Token reuse detected - revoke entire token family
		s.tokenService.RevokeTokenFamily(ctx, refreshClaims.Family)
		return nil, errors.New("invalid session")
	}

	// Check if session is expired
	if session.IsExpired() {
		s.authRepo.DeleteSession(ctx, refreshClaims.JTI)
		return nil, errors.New("session expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, refreshClaims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Blacklist the old refresh token
	if err := s.tokenService.RevokeToken(ctx, refreshClaims.JTI); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new token pair with same family
	roles := []string{"user"} // Should be fetched from user or database
	tokenPair, err := s.tokenService.GenerateTokenPair(user, roles, refreshClaims.Family)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Extract new claims
	newAccessClaims, err := s.tokenService.GetTokenClaims(tokenPair.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract new access token claims: %w", err)
	}

	newRefreshClaims, err := s.tokenService.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract new refresh token claims: %w", err)
	}

	// Update session info
	session.JTI = newAccessClaims.JTI
	session.RefreshJTI = newRefreshClaims.JTI
	session.ExpiresAt = newRefreshClaims.ExpiresAt.Time
	session.UpdateActivity()

	if err := s.authRepo.StoreSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Delete old session
	if err := s.authRepo.DeleteSession(ctx, refreshClaims.JTI); err != nil {
		// Log error but don't fail the refresh
		fmt.Printf("Failed to delete old session %s: %v\n", refreshClaims.JTI, err)
	}

	return tokenPair, nil
}

func (s *AuthenticationService) Logout(ctx context.Context, accessToken string, req *entity.LogoutRequest) error {
	// Extract claims from access token
	claims, err := s.tokenService.GetTokenClaims(accessToken)
	if err != nil {
		return fmt.Errorf("invalid access token: %w", err)
	}

	// Blacklist the access token
	if err := s.tokenService.RevokeToken(ctx, claims.JTI); err != nil {
		return fmt.Errorf("failed to revoke access token: %w", err)
	}

	// If refresh token is provided, revoke the specific session
	if req.RefreshToken != "" {
		refreshClaims, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
		if err == nil {
			// Blacklist refresh token and delete session
			s.tokenService.RevokeToken(ctx, refreshClaims.JTI)
			s.authRepo.DeleteSession(ctx, refreshClaims.JTI)
		}
	} else {
		// Revoke all sessions for the user
		if err := s.authRepo.DeleteSessionsByUserID(ctx, claims.UserID); err != nil {
			return fmt.Errorf("failed to revoke all sessions: %w", err)
		}
	}

	return nil
}

func (s *AuthenticationService) ValidateToken(ctx context.Context, tokenString string) (*entity.Claims, error) {
	// Validate token structure and signature
	claims, err := s.tokenService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	isBlacklisted, err := s.authRepo.IsBlacklisted(ctx, claims.JTI)
	if err != nil {
		return nil, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if isBlacklisted {
		return nil, errors.New("token has been revoked")
	}

	// Update session activity if session exists
	if err := s.authRepo.UpdateSessionActivity(ctx, claims.JTI); err != nil {
		// Log error but don't fail validation
		fmt.Printf("Failed to update session activity for token %s: %v\n", claims.JTI, err)
	}

	return claims, nil
}

func (s *AuthenticationService) RevokeAllSessions(ctx context.Context, userID string) error {
	return s.authRepo.DeleteSessionsByUserID(ctx, userID)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}