package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type RedisAuthRepository struct {
	client *redis.Client
	userRepo repository.UserRepository
}

func NewRedisAuthRepository(client *redis.Client, userRepo repository.UserRepository) *RedisAuthRepository {
	return &RedisAuthRepository{
		client: client,
		userRepo: userRepo,
	}
}

// セッション管理
func (r *RedisAuthRepository) StoreSession(ctx context.Context, session *entity.SessionInfo) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Store session by JTI
	sessionKey := fmt.Sprintf("session:%s", session.JTI)
	refreshKey := fmt.Sprintf("session:%s", session.RefreshJTI)
	userSessionsKey := fmt.Sprintf("user_sessions:%s", session.UserID)
	familyKey := fmt.Sprintf("token_family:%s", session.TokenFamily)

	pipe := r.client.Pipeline()
	
	// Store session data
	pipe.Set(ctx, sessionKey, data, time.Until(session.ExpiresAt))
	pipe.Set(ctx, refreshKey, data, time.Until(session.ExpiresAt))
	
	// Add to user's session list
	pipe.SAdd(ctx, userSessionsKey, session.JTI)
	pipe.Expire(ctx, userSessionsKey, 24*time.Hour*30) // 30 days
	
	// Add to token family
	pipe.SAdd(ctx, familyKey, session.JTI)
	pipe.Expire(ctx, familyKey, time.Until(session.ExpiresAt))

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisAuthRepository) GetSession(ctx context.Context, jti string) (*entity.SessionInfo, error) {
	sessionKey := fmt.Sprintf("session:%s", jti)
	data, err := r.client.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	var session entity.SessionInfo
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (r *RedisAuthRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*entity.SessionInfo, error) {
	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID)
	sessionJTIs, err := r.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return nil, err
	}

	var sessions []*entity.SessionInfo
	for _, jti := range sessionJTIs {
		session, err := r.GetSession(ctx, jti)
		if err != nil {
			// Skip expired or invalid sessions
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *RedisAuthRepository) UpdateSessionActivity(ctx context.Context, jti string) error {
	session, err := r.GetSession(ctx, jti)
	if err != nil {
		return err
	}

	session.UpdateActivity()
	return r.StoreSession(ctx, session)
}

func (r *RedisAuthRepository) DeleteSession(ctx context.Context, jti string) error {
	session, err := r.GetSession(ctx, jti)
	if err != nil {
		return nil // Session doesn't exist, nothing to delete
	}

	pipe := r.client.Pipeline()
	
	// Delete session keys
	sessionKey := fmt.Sprintf("session:%s", jti)
	refreshKey := fmt.Sprintf("session:%s", session.RefreshJTI)
	pipe.Del(ctx, sessionKey, refreshKey)
	
	// Remove from user sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%s", session.UserID)
	pipe.SRem(ctx, userSessionsKey, jti)
	
	// Remove from token family
	familyKey := fmt.Sprintf("token_family:%s", session.TokenFamily)
	pipe.SRem(ctx, familyKey, jti)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisAuthRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	sessions, err := r.GetSessionsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := r.DeleteSession(ctx, session.JTI); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete session %s: %v\n", session.JTI, err)
		}
	}

	// Clear user sessions set
	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID)
	return r.client.Del(ctx, userSessionsKey).Err()
}

func (r *RedisAuthRepository) DeleteSessionsByTokenFamily(ctx context.Context, tokenFamily string) error {
	familyKey := fmt.Sprintf("token_family:%s", tokenFamily)
	sessionJTIs, err := r.client.SMembers(ctx, familyKey).Result()
	if err != nil {
		return err
	}

	for _, jti := range sessionJTIs {
		if err := r.DeleteSession(ctx, jti); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete session %s from family: %v\n", jti, err)
		}
	}

	return r.client.Del(ctx, familyKey).Err()
}

// トークンブラックリスト
func (r *RedisAuthRepository) AddToBlacklist(ctx context.Context, jti string, expiry time.Time) error {
	blacklistKey := fmt.Sprintf("blacklist:%s", jti)
	return r.client.Set(ctx, blacklistKey, "revoked", time.Until(expiry)).Err()
}

func (r *RedisAuthRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	blacklistKey := fmt.Sprintf("blacklist:%s", jti)
	_, err := r.client.Get(ctx, blacklistKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisAuthRepository) CleanupExpiredBlacklist(ctx context.Context) error {
	// Redis handles TTL automatically, so this is a no-op
	return nil
}

// レート制限
func (r *RedisAuthRepository) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)
	
	current, err := r.client.Get(ctx, rateLimitKey).Result()
	if err == redis.Nil {
		return true, nil // No previous attempts
	}
	if err != nil {
		return false, err
	}

	attempts, err := strconv.Atoi(current)
	if err != nil {
		return false, err
	}

	return attempts < limit, nil
}

func (r *RedisAuthRepository) IncrementRateLimit(ctx context.Context, key string, window time.Duration) error {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)
	
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, window)
	
	_, err := pipe.Exec(ctx)
	return err
}

// ユーザー認証
func (r *RedisAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return r.userRepo.FindByEmail(ctx, email)
}

func (r *RedisAuthRepository) UpdateUserLastLogin(ctx context.Context, userID string) error {
	// This could be implemented to update user's last login in the main database
	return nil
}