package persistence

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryUserRepository_FindWithOffsetPagination(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	users := createTestUsers(t, repo, 25)

	t.Run("basic pagination", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 0, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 10, len(result.Users))
		assert.Equal(t, 25, result.Total)
	})

	t.Run("second page", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 10, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 10, len(result.Users))
		assert.Equal(t, 25, result.Total)
	})

	t.Run("last page", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 20, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 5, len(result.Users))
		assert.Equal(t, 25, result.Total)
	})

	t.Run("search filter", func(t *testing.T) {
		filter := repository.PaginationFilter{Search: "user1"}
		sort := repository.SortOption{Field: "name", Order: "asc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 0, filter, sort)
		require.NoError(t, err)
		assert.True(t, len(result.Users) > 0)
		for _, user := range result.Users {
			assert.Contains(t, user.Name, "user1")
		}
	})

	t.Run("sort by name ascending", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "name", Order: "asc"}

		result, err := repo.FindWithOffsetPagination(ctx, 5, 0, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 5, len(result.Users))

		for i := 1; i < len(result.Users); i++ {
			assert.True(t, result.Users[i-1].Name <= result.Users[i].Name)
		}
	})

	t.Run("status filter - active only", func(t *testing.T) {
		err := repo.SoftDelete(ctx, users[0].ID)
		require.NoError(t, err)

		filter := repository.PaginationFilter{Status: "active"}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 0, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 24, result.Total)
		for _, user := range result.Users {
			assert.False(t, user.IsDeleted())
		}
	})

	t.Run("status filter - deleted only", func(t *testing.T) {
		filter := repository.PaginationFilter{Status: "deleted"}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		result, err := repo.FindWithOffsetPagination(ctx, 10, 0, filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.True(t, result.Users[0].IsDeleted())
	})
}

func TestMemoryUserRepository_FindWithCursorPagination(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	createTestUsers(t, repo, 15)

	t.Run("first page", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		users, err := repo.FindWithCursorPagination(ctx, 5, "", filter, sort)
		require.NoError(t, err)
		assert.Equal(t, 5, len(users))
	})

	t.Run("with cursor", func(t *testing.T) {
		filter := repository.PaginationFilter{}
		sort := repository.SortOption{Field: "created_at", Order: "desc"}

		cursor := time.Now().Add(-5 * time.Minute).Format(time.RFC3339)
		users, err := repo.FindWithCursorPagination(ctx, 5, cursor, filter, sort)
		require.NoError(t, err)
		assert.True(t, len(users) <= 5)
	})

	t.Run("search with cursor", func(t *testing.T) {
		filter := repository.PaginationFilter{Search: "user"}
		sort := repository.SortOption{Field: "name", Order: "asc"}

		users, err := repo.FindWithCursorPagination(ctx, 3, "", filter, sort)
		require.NoError(t, err)
		assert.True(t, len(users) <= 4)
		for _, user := range users {
			assert.Contains(t, user.Name, "user")
		}
	})
}

func TestMemoryUserRepository_SortUsers(t *testing.T) {
	repo := &MemoryUserRepository{}

	now := time.Now()
	users := []*entity.User{
		{ID: "3", Name: "Charlie", Email: "charlie@example.com", CreatedAt: now.Add(-1 * time.Hour)},
		{ID: "1", Name: "Alice", Email: "alice@example.com", CreatedAt: now.Add(-3 * time.Hour)},
		{ID: "2", Name: "Bob", Email: "bob@example.com", CreatedAt: now.Add(-2 * time.Hour)},
	}

	t.Run("sort by name ascending", func(t *testing.T) {
		testUsers := make([]*entity.User, len(users))
		copy(testUsers, users)

		sort := repository.SortOption{Field: "name", Order: "asc"}
		repo.sortUsers(testUsers, sort)

		assert.Equal(t, "Alice", testUsers[0].Name)
		assert.Equal(t, "Bob", testUsers[1].Name)
		assert.Equal(t, "Charlie", testUsers[2].Name)
	})

	t.Run("sort by created_at descending", func(t *testing.T) {
		testUsers := make([]*entity.User, len(users))
		copy(testUsers, users)

		sort := repository.SortOption{Field: "created_at", Order: "desc"}
		repo.sortUsers(testUsers, sort)

		assert.Equal(t, "Charlie", testUsers[0].Name)
		assert.Equal(t, "Bob", testUsers[1].Name)
		assert.Equal(t, "Alice", testUsers[2].Name)
	})
}

func TestMemoryUserRepository_MatchesFilter(t *testing.T) {
	repo := &MemoryUserRepository{}

	user := &entity.User{
		ID:        "1",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC),
	}

	t.Run("search filter matches name", func(t *testing.T) {
		filter := repository.PaginationFilter{Search: "john"}
		assert.True(t, repo.matchesFilter(user, filter))
	})

	t.Run("search filter matches email", func(t *testing.T) {
		filter := repository.PaginationFilter{Search: "example"}
		assert.True(t, repo.matchesFilter(user, filter))
	})

	t.Run("search filter does not match", func(t *testing.T) {
		filter := repository.PaginationFilter{Search: "alice"}
		assert.False(t, repo.matchesFilter(user, filter))
	})

	t.Run("active status filter", func(t *testing.T) {
		filter := repository.PaginationFilter{Status: "active"}
		assert.True(t, repo.matchesFilter(user, filter))

		user.SoftDelete()
		assert.False(t, repo.matchesFilter(user, filter))
	})

	t.Run("date range filter", func(t *testing.T) {
		filter := repository.PaginationFilter{
			CreatedFrom: "2023-06-01",
			CreatedTo:   "2023-06-30",
		}
		assert.True(t, repo.matchesFilter(user, filter))

		filter = repository.PaginationFilter{
			CreatedFrom: "2023-07-01",
			CreatedTo:   "2023-07-31",
		}
		assert.False(t, repo.matchesFilter(user, filter))
	})
}

func createTestUsers(t *testing.T, repo repository.UserRepository, count int) []*entity.User {
	ctx := context.Background()
	users := make([]*entity.User, count)

	for i := 0; i < count; i++ {
		user, err := entity.NewUser(
			fmt.Sprintf("user%d@example.com", i),
			fmt.Sprintf("user%d", i),
		)
		require.NoError(t, err)

		user.CreatedAt = time.Now().Add(-time.Duration(count-i) * time.Minute)
		user.UpdatedAt = user.CreatedAt

		err = repo.Create(ctx, user)
		require.NoError(t, err)

		users[i] = user
	}

	return users
}