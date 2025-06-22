package persistence

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresUserRepository) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	
	repo := NewPostgresUserRepository(db).(*PostgresUserRepository)
	return db, mock, repo
}

func TestPostgresUserRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    *entity.User
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "successful user creation",
			user: &entity.User{
				ID:        "test-id",
				Name:      "Test User",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test-id", "Test User", "test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "database error during creation",
			user: &entity.User{
				ID:        "test-id",
				Name:      "Test User",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test-id", "Test User", "test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			err := repo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_FindByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		setup    func(mock sqlmock.Sqlmock)
		wantUser *entity.User
		wantErr  bool
		wantNil  bool
	}{
		{
			name: "successful find by ID",
			id:   "test-id",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
					AddRow("test-id", "Test User", "test@example.com", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE id = \\$1").
					WithArgs("test-id").
					WillReturnRows(rows)
			},
			wantUser: &entity.User{
				ID:    "test-id",
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "user not found",
			id:   "non-existent-id",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE id = \\$1").
					WithArgs("non-existent-id").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  false,
			wantNil:  true,
		},
		{
			name: "database error",
			id:   "test-id",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE id = \\$1").
					WithArgs("test-id").
					WillReturnError(sql.ErrConnDone)
			},
			wantUser: nil,
			wantErr:  true,
			wantNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			user, err := repo.FindByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, user)
			} else if tt.wantUser != nil {
				require.NotNil(t, user)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Name, user.Name)
				assert.Equal(t, tt.wantUser.Email, user.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_FindByEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		setup    func(mock sqlmock.Sqlmock)
		wantUser *entity.User
		wantErr  bool
		wantNil  bool
	}{
		{
			name:  "successful find by email",
			email: "test@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
					AddRow("test-id", "Test User", "test@example.com", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			wantUser: &entity.User{
				ID:    "test-id",
				Name:  "Test User",
				Email: "test@example.com",
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name:  "user not found by email",
			email: "nonexistent@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE email = \\$1").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  false,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			user, err := repo.FindByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, user)
			} else if tt.wantUser != nil {
				require.NotNil(t, user)
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Name, user.Name)
				assert.Equal(t, tt.wantUser.Email, user.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(mock sqlmock.Sqlmock)
		wantUsers []*entity.User
		wantErr   bool
	}{
		{
			name: "successful find all with multiple users",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
					AddRow("id1", "User 1", "user1@example.com", time.Now(), time.Now()).
					AddRow("id2", "User 2", "user2@example.com", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			wantUsers: []*entity.User{
				{ID: "id1", Name: "User 1", Email: "user1@example.com"},
				{ID: "id2", Name: "User 2", Email: "user2@example.com"},
			},
			wantErr: false,
		},
		{
			name: "successful find all with empty result",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			wantUsers: []*entity.User{},
			wantErr:   false,
		},
		{
			name: "database error during find all",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC").
					WillReturnError(sql.ErrConnDone)
			},
			wantUsers: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			users, err := repo.FindAll(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantUsers), len(users))
				for i, expectedUser := range tt.wantUsers {
					assert.Equal(t, expectedUser.ID, users[i].ID)
					assert.Equal(t, expectedUser.Name, users[i].Name)
					assert.Equal(t, expectedUser.Email, users[i].Email)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		user    *entity.User
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "successful user update",
			user: &entity.User{
				ID:    "test-id",
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\$2, email = \\$3, updated_at = \\$4 WHERE id = \\$1").
					WithArgs("test-id", "Updated User", "updated@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found during update",
			user: &entity.User{
				ID:    "non-existent-id",
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\$2, email = \\$3, updated_at = \\$4 WHERE id = \\$1").
					WithArgs("non-existent-id", "Updated User", "updated@example.com", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "database error during update",
			user: &entity.User{
				ID:    "test-id",
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\$2, email = \\$3, updated_at = \\$4 WHERE id = \\$1").
					WithArgs("test-id", "Updated User", "updated@example.com", sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			err := repo.Update(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "successful user deletion",
			id:   "test-id",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
					WithArgs("test-id").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found during deletion",
			id:   "non-existent-id",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
					WithArgs("non-existent-id").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "database error during deletion",
			id:   "test-id",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
					WithArgs("test-id").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			err := repo.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresUserRepository_List(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		offset    int
		setup     func(mock sqlmock.Sqlmock)
		wantUsers []*entity.User
		wantErr   bool
	}{
		{
			name:   "successful list with pagination",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
					AddRow("id1", "User 1", "user1@example.com", time.Now(), time.Now()).
					AddRow("id2", "User 2", "user2@example.com", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantUsers: []*entity.User{
				{ID: "id1", Name: "User 1", Email: "user1@example.com"},
				{ID: "id2", Name: "User 2", Email: "user2@example.com"},
			},
			wantErr: false,
		},
		{
			name:   "successful list with offset",
			limit:  5,
			offset: 10,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
					WithArgs(5, 10).
					WillReturnRows(rows)
			},
			wantUsers: []*entity.User{},
			wantErr:   false,
		},
		{
			name:   "database error during list",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \\$1 OFFSET \\$2").
					WithArgs(10, 0).
					WillReturnError(sql.ErrConnDone)
			},
			wantUsers: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, repo := setupMockDB(t)
			defer db.Close()

			tt.setup(mock)

			users, err := repo.List(context.Background(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantUsers), len(users))
				for i, expectedUser := range tt.wantUsers {
					assert.Equal(t, expectedUser.ID, users[i].ID)
					assert.Equal(t, expectedUser.Name, users[i].Name)
					assert.Equal(t, expectedUser.Email, users[i].Email)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}