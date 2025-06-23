package usecase

import (
	"context"
	"testing"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockArticleRepository は ArticleRepository のモック
type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) FindByID(ctx context.Context, id string) (*entity.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Article), args.Error(1)
}

func (m *MockArticleRepository) FindAll(ctx context.Context, filter repository.ArticleFilter) ([]*entity.Article, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Article), args.Int(1), args.Error(2)
}

func (m *MockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleRepository) FindByStatus(ctx context.Context, status entity.ArticleStatus, limit, offset int) ([]*entity.Article, int, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Article), args.Int(1), args.Error(2)
}

func (m *MockArticleRepository) FindByTag(ctx context.Context, tag string, limit, offset int) ([]*entity.Article, int, error) {
	args := m.Called(ctx, tag, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Article), args.Int(1), args.Error(2)
}

func (m *MockArticleRepository) FindByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Article, int, error) {
	args := m.Called(ctx, authorID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Article), args.Int(1), args.Error(2)
}

func TestCreateArticleUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		input       CreateArticleInput
		setupMock   func(*MockArticleRepository)
		expectError bool
	}{
		{
			name: "記事作成成功",
			input: CreateArticleInput{
				Title:    "テスト記事",
				Content:  "記事の内容",
				Summary:  "記事の概要",
				AuthorID: "test-author-id",
				Tags:     []string{"Go", "テスト"},
			},
			setupMock: func(repo *MockArticleRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Article")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "リポジトリエラーの場合",
			input: CreateArticleInput{
				Title:    "テスト記事",
				Content:  "記事の内容",
				AuthorID: "test-author-id",
			},
			setupMock: func(repo *MockArticleRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Article")).Return(assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			tt.setupMock(mockRepo)

			usecase := NewCreateArticleUseCase(mockRepo)
			output, err := usecase.Execute(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tt.input.Title, output.Article.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetArticleUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		input       GetArticleInput
		setupMock   func(*MockArticleRepository)
		expectError bool
	}{
		{
			name: "記事取得成功",
			input: GetArticleInput{
				ID: "test-id",
			},
			setupMock: func(repo *MockArticleRepository) {
				article := &entity.Article{
					ID:      "test-id",
					Title:   "テスト記事",
					Content: "記事の内容",
				}
				repo.On("FindByID", mock.Anything, "test-id").Return(article, nil)
			},
			expectError: false,
		},
		{
			name: "記事が見つからない場合エラー",
			input: GetArticleInput{
				ID: "not-found-id",
			},
			setupMock: func(repo *MockArticleRepository) {
				repo.On("FindByID", mock.Anything, "not-found-id").Return(nil, assert.AnError)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			tt.setupMock(mockRepo)

			usecase := NewGetArticleUseCase(mockRepo)
			output, err := usecase.Execute(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tt.input.ID, output.Article.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}