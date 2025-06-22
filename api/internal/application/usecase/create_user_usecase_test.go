package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/domain/service"
)

// MockUserRepository はテスト用のUserRepositoryのモック実装
type MockUserRepository struct {
	users          map[string]*entity.User
	emailToUserMap map[string]*entity.User
	createError      error
	findByIDError    error
	findByEmailError error
	findAllError     error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:          make(map[string]*entity.User),
		emailToUserMap: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.createError != nil {
		return m.createError
	}
	m.users[user.ID] = user
	m.emailToUserMap[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	if m.findByIDError != nil {
		return nil, m.findByIDError
	}
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findByEmailError != nil {
		return nil, m.findByEmailError
	}
	if user, exists := m.emailToUserMap[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	if m.findAllError != nil {
		return nil, m.findAllError
	}
	var users []*entity.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	m.emailToUserMap[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emailToUserMap, user.Email)
	}
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	count := 0
	for _, user := range m.users {
		if count >= offset && len(users) < limit {
			users = append(users, user)
		}
		count++
	}
	return users, nil
}

func (m *MockUserRepository) SetCreateError(err error) {
	m.createError = err
}

func (m *MockUserRepository) SetFindByIDError(err error) {
	m.findByIDError = err
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, id string) error {
	if user, exists := m.users[id]; exists {
		user.SoftDelete()
	}
	return nil
}

func (m *MockUserRepository) HardDelete(ctx context.Context, id string) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emailToUserMap, user.Email)
	}
	return nil
}

func (m *MockUserRepository) FindWithOffsetPagination(ctx context.Context, limit, offset int, filter repository.PaginationFilter, sort repository.SortOption) (*repository.PaginationResult, error) {
	users, _ := m.List(ctx, limit, offset)
	return &repository.PaginationResult{
		Users: users,
		Total: len(m.users),
	}, nil
}

func (m *MockUserRepository) FindWithCursorPagination(ctx context.Context, limit int, cursor string, filter repository.PaginationFilter, sort repository.SortOption) ([]*entity.User, error) {
	var users []*entity.User
	count := 0
	for _, user := range m.users {
		if count < limit {
			users = append(users, user)
		}
		count++
	}
	return users, nil
}

func (m *MockUserRepository) SetFindByEmailError(err error) {
	m.findByEmailError = err
}

func TestCreateUserUseCase_Execute(t *testing.T) {
	type fields struct {
		userService *service.UserService
	}
	type args struct {
		ctx   context.Context
		input CreateUserInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(*MockUserRepository)
		want    func(*CreateUserOutput) bool
		wantErr bool
	}{
		{
			name: "Should_CreateUser_When_ValidInputProvided",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email: "test@example.com",
					Name:  "Test User",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// メールが存在しないことを模擬
				mockRepo.SetFindByEmailError(errors.New("user not found"))
			},
			want: func(output *CreateUserOutput) bool {
				return output != nil &&
					output.User != nil &&
					output.User.Email == "test@example.com" &&
					output.User.Name == "Test User" &&
					output.User.ID != ""
			},
			wantErr: false,
		},
		{
			name: "Should_ReturnError_When_EmailAlreadyExists",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email: "existing@example.com",
					Name:  "Test User",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 既存ユーザーを作成
				existingUser, _ := entity.NewUser("existing@example.com", "Existing User")
				mockRepo.Create(context.Background(), existingUser)
			},
			want: func(output *CreateUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_ReturnError_When_EmailIsEmpty",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email: "",
					Name:  "Test User",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 特別な設定は不要
			},
			want: func(output *CreateUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_ReturnError_When_NameIsEmpty",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email: "test@example.com",
					Name:  "",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// メールが存在しないことを模擬
				mockRepo.SetFindByEmailError(errors.New("user not found"))
			},
			want: func(output *CreateUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_ReturnError_When_RepositoryCreateFails",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email: "test@example.com",
					Name:  "Test User",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// メールが存在しないことを模擬
				mockRepo.SetFindByEmailError(errors.New("user not found"))
				// リポジトリの作成でエラーを模擬
				mockRepo.SetCreateError(errors.New("database error"))
			},
			want: func(output *CreateUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := NewMockUserRepository()
			tt.setup(mockRepo)

			// UserServiceとCreateUserUseCaseの作成
			userService := service.NewUserService(mockRepo)
			uc := NewCreateUserUseCase(userService)

			// テスト実行
			got, err := uc.Execute(tt.args.ctx, tt.args.input)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果の検証
			if !tt.want(got) {
				t.Errorf("CreateUserUseCase.Execute() result validation failed")
			}
		})
	}
}