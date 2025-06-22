package persistence

import (
	"context"
	"sync"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type MemoryAuthRepository struct {
	sessions   map[string]*entity.SessionInfo
	blacklist  map[string]time.Time
	rateLimits map[string]int
	mutex      sync.RWMutex
	userRepo   repository.UserRepository
}

func NewMemoryAuthRepository(userRepo repository.UserRepository) *MemoryAuthRepository {
	return &MemoryAuthRepository{
		sessions:   make(map[string]*entity.SessionInfo),
		blacklist:  make(map[string]time.Time),
		rateLimits: make(map[string]int),
		userRepo:   userRepo,
	}
}

func (m *MemoryAuthRepository) StoreSession(ctx context.Context, session *entity.SessionInfo) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[session.RefreshJTI] = session
	return nil
}

func (m *MemoryAuthRepository) GetSession(ctx context.Context, sessionID string) (*entity.SessionInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	session, exists := m.sessions[sessionID]
	if !exists {
		return nil, nil
	}
	return session, nil
}

func (m *MemoryAuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.sessions, sessionID)
	return nil
}

func (m *MemoryAuthRepository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for sessionID, session := range m.sessions {
		if session.UserID == userID {
			delete(m.sessions, sessionID)
		}
	}
	return nil
}

func (m *MemoryAuthRepository) GetUserSessions(ctx context.Context, userID string) ([]*entity.SessionInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var sessions []*entity.SessionInfo
	for _, session := range m.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (m *MemoryAuthRepository) UpdateSessionExpiry(ctx context.Context, sessionID string, expiresAt time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if session, exists := m.sessions[sessionID]; exists {
		session.ExpiresAt = expiresAt
	}
	return nil
}

// Interface methods to satisfy AuthRepository
func (m *MemoryAuthRepository) GetSessionsByUserID(ctx context.Context, userID string) ([]*entity.SessionInfo, error) {
	return m.GetUserSessions(ctx, userID)
}

func (m *MemoryAuthRepository) UpdateSessionActivity(ctx context.Context, jti string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if session, exists := m.sessions[jti]; exists {
		session.UpdateActivity()
	}
	return nil
}

func (m *MemoryAuthRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	return m.DeleteAllUserSessions(ctx, userID)
}

func (m *MemoryAuthRepository) DeleteSessionsByTokenFamily(ctx context.Context, tokenFamily string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for sessionID, session := range m.sessions {
		if session.TokenFamily == tokenFamily {
			delete(m.sessions, sessionID)
		}
	}
	return nil
}

func (m *MemoryAuthRepository) AddToBlacklist(ctx context.Context, jti string, expiry time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.blacklist[jti] = expiry
	return nil
}

func (m *MemoryAuthRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	expiry, exists := m.blacklist[jti]
	if !exists {
		return false, nil
	}
	if time.Now().After(expiry) {
		delete(m.blacklist, jti)
		return false, nil
	}
	return true, nil
}

func (m *MemoryAuthRepository) CleanupExpiredBlacklist(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	now := time.Now()
	for jti, expiry := range m.blacklist {
		if now.After(expiry) {
			delete(m.blacklist, jti)
		}
	}
	return nil
}

func (m *MemoryAuthRepository) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	count, exists := m.rateLimits[key]
	if !exists {
		return true, nil
	}
	return count < limit, nil
}

func (m *MemoryAuthRepository) IncrementRateLimit(ctx context.Context, key string, window time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.rateLimits[key]++
	return nil
}

func (m *MemoryAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.userRepo.FindByEmail(ctx, email)
}

func (m *MemoryAuthRepository) UpdateUserLastLogin(ctx context.Context, userID string) error {
	// This is a simple implementation, in a real app you'd update the user record
	return nil
}