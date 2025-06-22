package persistence_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"
)

type PostgresUserRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo *persistence.PostgresUserRepository
}

func (suite *PostgresUserRepositoryTestSuite) SetupSuite() {
	// This would typically connect to a test database
	// For now, we'll skip if no test DB is available
	suite.T().Skip("Integration test requires database setup")
	
	// Example setup:
	// db, err := sql.Open("postgres", "postgres://user:pass@localhost/testdb?sslmode=disable")
	// require.NoError(suite.T(), err)
	// suite.db = db
	// suite.repo = persistence.NewPostgresUserRepository(db)
}

func (suite *PostgresUserRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *PostgresUserRepositoryTestSuite) SetupTest() {
	// Clean up test data before each test
	if suite.db != nil {
		_, err := suite.db.Exec("TRUNCATE TABLE users CASCADE")
		require.NoError(suite.T(), err)
	}
}

func (suite *PostgresUserRepositoryTestSuite) TestSoftDelete_Success() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Create a test user
	user, err := entity.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)
	user.Password = "hashedpassword"
	
	err = suite.repo.Create(ctx, user)
	require.NoError(suite.T(), err)
	
	// Soft delete the user
	err = suite.repo.SoftDelete(ctx, user.ID)
	assert.NoError(suite.T(), err)
	
	// Verify user cannot be found with normal queries
	foundUser, err := suite.repo.FindByID(ctx, user.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
	
	// Verify user is marked as deleted in database
	var deletedAt sql.NullTime
	err = suite.db.QueryRowContext(ctx, "SELECT deleted_at FROM users WHERE id = $1", user.ID).Scan(&deletedAt)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), deletedAt.Valid)
	assert.WithinDuration(suite.T(), time.Now(), deletedAt.Time, 5*time.Second)
}

func (suite *PostgresUserRepositoryTestSuite) TestHardDelete_Success() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Create a test user
	user, err := entity.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)
	user.Password = "hashedpassword"
	
	err = suite.repo.Create(ctx, user)
	require.NoError(suite.T(), err)
	
	// Hard delete the user
	err = suite.repo.HardDelete(ctx, user.ID)
	assert.NoError(suite.T(), err)
	
	// Verify user is completely removed from database
	var count int
	err = suite.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *PostgresUserRepositoryTestSuite) TestSoftDelete_UserNotFound() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Try to soft delete non-existent user
	err := suite.repo.SoftDelete(ctx, "non-existent-id")
	assert.Equal(suite.T(), sql.ErrNoRows, err)
}

func (suite *PostgresUserRepositoryTestSuite) TestHardDelete_UserNotFound() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Try to hard delete non-existent user
	err := suite.repo.HardDelete(ctx, "non-existent-id")
	assert.Equal(suite.T(), sql.ErrNoRows, err)
}

func (suite *PostgresUserRepositoryTestSuite) TestSoftDelete_AlreadyDeleted() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Create a test user
	user, err := entity.NewUser("test@example.com", "Test User")
	require.NoError(suite.T(), err)
	user.Password = "hashedpassword"
	
	err = suite.repo.Create(ctx, user)
	require.NoError(suite.T(), err)
	
	// Soft delete the user
	err = suite.repo.SoftDelete(ctx, user.ID)
	require.NoError(suite.T(), err)
	
	// Try to soft delete again
	err = suite.repo.SoftDelete(ctx, user.ID)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
}

func (suite *PostgresUserRepositoryTestSuite) TestFindByID_ExcludesSoftDeleted() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Create two test users
	user1, err := entity.NewUser("test1@example.com", "Test User 1")
	require.NoError(suite.T(), err)
	user1.Password = "hashedpassword"
	
	user2, err := entity.NewUser("test2@example.com", "Test User 2")
	require.NoError(suite.T(), err)
	user2.Password = "hashedpassword"
	
	err = suite.repo.Create(ctx, user1)
	require.NoError(suite.T(), err)
	
	err = suite.repo.Create(ctx, user2)
	require.NoError(suite.T(), err)
	
	// Soft delete user1
	err = suite.repo.SoftDelete(ctx, user1.ID)
	require.NoError(suite.T(), err)
	
	// FindByID should not return soft deleted user
	foundUser1, err := suite.repo.FindByID(ctx, user1.ID)
	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), foundUser1)
	
	// FindByID should still return active user
	foundUser2, err := suite.repo.FindByID(ctx, user2.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser2)
	assert.Equal(suite.T(), user2.ID, foundUser2.ID)
}

func (suite *PostgresUserRepositoryTestSuite) TestFindAll_ExcludesSoftDeleted() {
	if suite.db == nil {
		suite.T().Skip("Database not available")
	}
	
	ctx := context.Background()
	
	// Create two test users
	user1, err := entity.NewUser("test1@example.com", "Test User 1")
	require.NoError(suite.T(), err)
	user1.Password = "hashedpassword"
	
	user2, err := entity.NewUser("test2@example.com", "Test User 2")
	require.NoError(suite.T(), err)
	user2.Password = "hashedpassword"
	
	err = suite.repo.Create(ctx, user1)
	require.NoError(suite.T(), err)
	
	err = suite.repo.Create(ctx, user2)
	require.NoError(suite.T(), err)
	
	// Soft delete user1
	err = suite.repo.SoftDelete(ctx, user1.ID)
	require.NoError(suite.T(), err)
	
	// FindAll should only return active user
	users, err := suite.repo.FindAll(ctx)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), user2.ID, users[0].ID)
}

func TestPostgresUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresUserRepositoryTestSuite))
}