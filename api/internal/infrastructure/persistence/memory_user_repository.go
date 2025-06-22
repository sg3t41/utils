package persistence

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
	index map[string]*entity.User // email index
}

func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*entity.User),
		index: make(map[string]*entity.User),
	}
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.ID] = user
	r.index[user.Email] = user
	return nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.index[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	r.users[user.ID] = user
	r.index[user.Email] = user
	return nil
}

func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	delete(r.index, user.Email)
	return nil
}

func (r *MemoryUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *MemoryUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	start := offset
	if start > len(users) {
		return []*entity.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func (r *MemoryUserRepository) SoftDelete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.New("user not found")
	}

	if user.IsDeleted() {
		return errors.New("user already deleted")
	}

	user.SoftDelete()
	return nil
}

func (r *MemoryUserRepository) HardDelete(ctx context.Context, id string) error {
	return r.Delete(ctx, id)
}

func (r *MemoryUserRepository) FindWithOffsetPagination(ctx context.Context, limit, offset int, filter repository.PaginationFilter, sortOption repository.SortOption) (*repository.PaginationResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0)
	for _, user := range r.users {
		if r.matchesFilter(user, filter) {
			users = append(users, user)
		}
	}

	r.sortUsers(users, sortOption)

	total := len(users)
	start := offset
	if start > total {
		return &repository.PaginationResult{
			Users: []*entity.User{},
			Total: total,
		}, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	return &repository.PaginationResult{
		Users: users[start:end],
		Total: total,
	}, nil
}

func (r *MemoryUserRepository) FindWithCursorPagination(ctx context.Context, limit int, cursor string, filter repository.PaginationFilter, sortOption repository.SortOption) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0)
	for _, user := range r.users {
		if r.matchesFilter(user, filter) {
			users = append(users, user)
		}
	}

	r.sortUsers(users, sortOption)

	if cursor != "" {
		cursorTime, err := time.Parse(time.RFC3339, cursor)
		if err != nil {
			return nil, err
		}

		filteredUsers := make([]*entity.User, 0)
		for _, user := range users {
			if user.CreatedAt.Before(cursorTime) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	if len(users) > limit {
		return users[:limit+1], nil
	}

	return users, nil
}

func (r *MemoryUserRepository) matchesFilter(user *entity.User, filter repository.PaginationFilter) bool {
	if filter.Search != "" {
		searchLower := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(user.Name), searchLower) &&
			!strings.Contains(strings.ToLower(user.Email), searchLower) {
			return false
		}
	}

	if filter.Status != "" {
		switch filter.Status {
		case "active":
			if user.IsDeleted() {
				return false
			}
		case "deleted":
			if !user.IsDeleted() {
				return false
			}
		}
	}

	if filter.CreatedFrom != "" {
		createdFrom, err := time.Parse("2006-01-02", filter.CreatedFrom)
		if err == nil && user.CreatedAt.Before(createdFrom) {
			return false
		}
	}

	if filter.CreatedTo != "" {
		createdTo, err := time.Parse("2006-01-02", filter.CreatedTo)
		if err == nil && user.CreatedAt.After(createdTo.Add(24*time.Hour)) {
			return false
		}
	}

	return true
}

func (r *MemoryUserRepository) sortUsers(users []*entity.User, sortOption repository.SortOption) {
	field := sortOption.Field
	if field == "" {
		field = "created_at"
	}

	order := sortOption.Order
	if order == "" {
		order = "desc"
	}

	sort.Slice(users, func(i, j int) bool {
		var compare bool
		switch field {
		case "id":
			compare = users[i].ID < users[j].ID
		case "name":
			compare = users[i].Name < users[j].Name
		case "email":
			compare = users[i].Email < users[j].Email
		case "created_at":
			compare = users[i].CreatedAt.Before(users[j].CreatedAt)
		case "updated_at":
			compare = users[i].UpdatedAt.Before(users[j].UpdatedAt)
		default:
			compare = users[i].CreatedAt.Before(users[j].CreatedAt)
		}

		if order == "desc" {
			return !compare
		}
		return compare
	})
}