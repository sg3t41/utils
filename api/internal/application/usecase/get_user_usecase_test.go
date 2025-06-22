package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/service"
)

func TestGetUserUseCase_Execute(t *testing.T) {
	type fields struct {
		userService *service.UserService
	}
	type args struct {
		ctx   context.Context
		input GetUserInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(*MockUserRepository)
		want    func(*GetUserOutput) bool
		wantErr bool
	}{
		{
			name: "Should_GetUser_When_ValidIDProvided",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "test-user-id",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// テストユーザーを作成
				testUser, _ := entity.NewUser("test@example.com", "Test User")
				testUser.ID = "test-user-id"
				mockRepo.Create(context.Background(), testUser)
			},
			want: func(output *GetUserOutput) bool {
				return output != nil &&
					output.User != nil &&
					output.User.ID == "test-user-id" &&
					output.User.Email == "test@example.com" &&
					output.User.Name == "Test User"
			},
			wantErr: false,
		},
		{
			name: "Should_ReturnError_When_UserNotFound",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "non-existent-id",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 特別な設定は不要（ユーザーが存在しない状態）
			},
			want: func(output *GetUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_ReturnError_When_EmptyIDProvided",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 特別な設定は不要
			},
			want: func(output *GetUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_ReturnError_When_RepositoryFails",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "test-user-id",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// リポジトリでエラーを模擬
				mockRepo.SetFindByIDError(errors.New("database connection failed"))
			},
			want: func(output *GetUserOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_GetUser_When_UserExists_ButHasSpecialCharacters",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "user-with-special-chars",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 特殊文字を含む名前のユーザーを作成
				testUser, _ := entity.NewUser("special@example.com", "テスト ユーザー")
				testUser.ID = "user-with-special-chars"
				mockRepo.Create(context.Background(), testUser)
			},
			want: func(output *GetUserOutput) bool {
				return output != nil &&
					output.User != nil &&
					output.User.ID == "user-with-special-chars" &&
					output.User.Email == "special@example.com" &&
					output.User.Name == "テスト ユーザー"
			},
			wantErr: false,
		},
		{
			name: "Should_GetUser_When_UserExists_AndHasBeenUpdated",
			fields: fields{
				userService: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUserInput{
					ID: "updated-user-id",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// ユーザーを作成して更新
				testUser, _ := entity.NewUser("original@example.com", "Original Name")
				testUser.ID = "updated-user-id"
				mockRepo.Create(context.Background(), testUser)
				
				// 名前を更新
				testUser.UpdateName("Updated Name")
				mockRepo.Update(context.Background(), testUser)
			},
			want: func(output *GetUserOutput) bool {
				return output != nil &&
					output.User != nil &&
					output.User.ID == "updated-user-id" &&
					output.User.Name == "Updated Name" &&
					output.User.Version == 2 // 更新されているのでバージョンが2
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := NewMockUserRepository()
			tt.setup(mockRepo)

			// UserServiceとGetUserUseCaseの作成
			userService := service.NewUserService(mockRepo)
			uc := NewGetUserUseCase(userService)

			// テスト実行
			got, err := uc.Execute(tt.args.ctx, tt.args.input)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果の検証
			if !tt.want(got) {
				t.Errorf("GetUserUseCase.Execute() result validation failed")
				if got != nil && got.User != nil {
					t.Errorf("Got user: ID=%s, Email=%s, Name=%s, Version=%d", 
						got.User.ID, got.User.Email, got.User.Name, got.User.Version)
				}
			}
		})
	}
}