package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/infrastructure/persistence"
)

func main() {
	repo := persistence.NewMemoryUserRepository()
	ctx := context.Background()

	fmt.Println("Creating test users...")
	users := createTestUsers(repo, 15)
	fmt.Printf("Created %d users\n", len(users))

	fmt.Println("\n=== Testing Offset Pagination ===")
	testOffsetPagination(repo, ctx)

	fmt.Println("\n=== Testing Cursor Pagination ===")
	testCursorPagination(repo, ctx)

	fmt.Println("\n=== Testing Filters ===")
	testFilters(repo, ctx)
}

func createTestUsers(repo repository.UserRepository, count int) []*entity.User {
	ctx := context.Background()
	users := make([]*entity.User, count)

	for i := 0; i < count; i++ {
		user, err := entity.NewUser(
			fmt.Sprintf("user%d@example.com", i),
			fmt.Sprintf("User %d", i),
		)
		if err != nil {
			log.Fatal(err)
		}

		user.CreatedAt = time.Now().Add(-time.Duration(count-i) * time.Minute)
		user.UpdatedAt = user.CreatedAt

		err = repo.Create(ctx, user)
		if err != nil {
			log.Fatal(err)
		}

		users[i] = user
	}

	return users
}

func testOffsetPagination(repo repository.UserRepository, ctx context.Context) {
	filter := repository.PaginationFilter{}
	sort := repository.SortOption{Field: "created_at", Order: "desc"}

	fmt.Println("First page (limit=5, offset=0):")
	result, err := repo.FindWithOffsetPagination(ctx, 5, 0, filter, sort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d users out of %d total\n", len(result.Users), result.Total)
	for i, user := range result.Users {
		fmt.Printf("  %d. %s (%s)\n", i+1, user.Name, user.Email)
	}

	fmt.Println("\nSecond page (limit=5, offset=5):")
	result, err = repo.FindWithOffsetPagination(ctx, 5, 5, filter, sort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d users out of %d total\n", len(result.Users), result.Total)
	for i, user := range result.Users {
		fmt.Printf("  %d. %s (%s)\n", i+1, user.Name, user.Email)
	}
}

func testCursorPagination(repo repository.UserRepository, ctx context.Context) {
	filter := repository.PaginationFilter{}
	sort := repository.SortOption{Field: "created_at", Order: "desc"}

	fmt.Println("First page (limit=3):")
	users, err := repo.FindWithCursorPagination(ctx, 3, "", filter, sort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d users\n", len(users))
	for i, user := range users {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, user.Name, user.Email, user.CreatedAt.Format(time.RFC3339))
	}

	if len(users) > 0 {
		cursor := users[len(users)-1].CreatedAt.Format(time.RFC3339)
		fmt.Printf("\nNext page with cursor %s:\n", cursor)
		users, err = repo.FindWithCursorPagination(ctx, 3, cursor, filter, sort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  Found %d users\n", len(users))
		for i, user := range users {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, user.Name, user.Email, user.CreatedAt.Format(time.RFC3339))
		}
	}
}

func testFilters(repo repository.UserRepository, ctx context.Context) {
	filter := repository.PaginationFilter{Search: "User 1"}
	sort := repository.SortOption{Field: "name", Order: "asc"}

	fmt.Println("Search filter 'User 1':")
	result, err := repo.FindWithOffsetPagination(ctx, 10, 0, filter, sort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d users out of %d total\n", len(result.Users), result.Total)
	for i, user := range result.Users {
		fmt.Printf("  %d. %s (%s)\n", i+1, user.Name, user.Email)
	}

	fmt.Println("\nSort by name ascending:")
	filter = repository.PaginationFilter{}
	sort = repository.SortOption{Field: "name", Order: "asc"}

	result, err = repo.FindWithOffsetPagination(ctx, 5, 0, filter, sort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Found %d users\n", len(result.Users))
	for i, user := range result.Users {
		fmt.Printf("  %d. %s (%s)\n", i+1, user.Name, user.Email)
	}
}