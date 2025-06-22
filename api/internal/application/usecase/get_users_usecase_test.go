package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sg3t41/api/internal/domain/entity"
)

func TestGetUsersUseCase_Execute(t *testing.T) {
	type fields struct {
		userRepository *MockUserRepository
	}
	type args struct {
		ctx   context.Context
		input GetUsersInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(*MockUserRepository)
		want    func(*GetUsersOutput) bool
		wantErr bool
	}{
		{
			name: "Should_GetAllUsers_When_NoLimitProvided",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit: 0, // 制限なし
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// テストユーザーを3人作成
				user1, _ := entity.NewUser("user1@example.com", "User 1")
				user2, _ := entity.NewUser("user2@example.com", "User 2")
				user3, _ := entity.NewUser("user3@example.com", "User 3")
				mockRepo.Create(context.Background(), user1)
				mockRepo.Create(context.Background(), user2)
				mockRepo.Create(context.Background(), user3)
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) == 3 &&
					output.Total == 3
			},
			wantErr: false,
		},
		{
			name: "Should_GetUsersWithPagination_When_LimitAndPageProvided",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit:     2,
					Page:      1,
					UseCursor: false,
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// テストユーザーを5人作成
				for i := 1; i <= 5; i++ {
					user, _ := entity.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("User %d", i))
					mockRepo.Create(context.Background(), user)
				}
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) == 2 &&
					output.Total == 5
			},
			wantErr: false,
		},
		{
			name: "Should_GetUsersWithCursorPagination_When_CursorEnabled",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit:     3,
					Cursor:    "cursor123",
					UseCursor: true,
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// テストユーザーを5人作成
				for i := 1; i <= 5; i++ {
					user, _ := entity.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("User %d", i))
					mockRepo.Create(context.Background(), user)
				}
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) <= 3 &&
					output.Total == 0 // カーソルベースの場合、Totalは0
			},
			wantErr: false,
		},
		{
			name: "Should_GetUsersWithFilter_When_SearchProvided",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit:  10,
					Page:   1,
					Search: "admin",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 検索対象のユーザーを作成
				admin, _ := entity.NewUser("admin@example.com", "Admin User")
				user, _ := entity.NewUser("user@example.com", "Normal User")
				mockRepo.Create(context.Background(), admin)
				mockRepo.Create(context.Background(), user)
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					output.Users != nil &&
					output.Total >= 0
			},
			wantErr: false,
		},
		{
			name: "Should_GetUsersWithSort_When_SortOptionProvided",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit: 10,
					Page:  1,
					Sort:  "created_at",
					Order: "desc",
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 作成時間が異なるユーザーを作成
				user1, _ := entity.NewUser("user1@example.com", "User 1")
				user2, _ := entity.NewUser("user2@example.com", "User 2")
				mockRepo.Create(context.Background(), user1)
				mockRepo.Create(context.Background(), user2)
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) == 2 &&
					output.Total == 2
			},
			wantErr: false,
		},
		{
			name: "Should_ReturnEmptyList_When_NoUsersExist",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit: 10,
					Page:  1,
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// ユーザーを作成しない
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) == 0 &&
					output.Total == 0
			},
			wantErr: false,
		},
		{
			name: "Should_ReturnError_When_RepositoryFails",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit: 0, // FindAllを呼ぶ
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// FindAllでエラーを発生させる
				mockRepo.findAllError = errors.New("database connection failed")
			},
			want: func(output *GetUsersOutput) bool {
				return output == nil
			},
			wantErr: true,
		},
		{
			name: "Should_HandleLargePageNumbers_When_PageExceedsData",
			fields: fields{
				userRepository: nil, // setupで設定
			},
			args: args{
				ctx: context.Background(),
				input: GetUsersInput{
					Limit: 10,
					Page:  100, // 大きなページ番号
				},
			},
			setup: func(mockRepo *MockUserRepository) {
				// 少数のユーザーを作成
				user, _ := entity.NewUser("user@example.com", "User")
				mockRepo.Create(context.Background(), user)
			},
			want: func(output *GetUsersOutput) bool {
				return output != nil &&
					len(output.Users) == 0 &&
					output.Total == 1
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := NewMockUserRepository()
			tt.setup(mockRepo)

			// GetUsersUseCaseの作成
			uc := NewGetUsersUseCase(mockRepo)

			// テスト実行
			got, err := uc.Execute(tt.args.ctx, tt.args.input)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsersUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果の検証
			if !tt.want(got) {
				t.Errorf("GetUsersUseCase.Execute() result validation failed")
				if got != nil {
					t.Errorf("Got output: Users count=%d, Total=%d", len(got.Users), got.Total)
				}
			}
		})
	}
}