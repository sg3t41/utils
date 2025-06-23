package usecase

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
	"github.com/sg3t41/api/internal/interfaces/dto"
)

type CreateLineUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateLineUserUseCase(userRepo repository.UserRepository) *CreateLineUserUseCase {
	return &CreateLineUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *CreateLineUserUseCase) Execute(ctx context.Context, lineUserID, displayName, profileImage string) (*entity.User, error) {
	// 既存ユーザーをチェック
	existingUser, err := uc.userRepo.FindByLineUserID(ctx, lineUserID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return existingUser, nil
	}

	// 新規ユーザーを作成・保存
	return uc.createNewUser(ctx, lineUserID, displayName, profileImage)
}

// createNewUser 新規ユーザーの作成
func (uc *CreateLineUserUseCase) createNewUser(ctx context.Context, lineUserID, displayName, profileImage string) (*entity.User, error) {
	user, err := entity.NewLineUser(lineUserID, displayName, profileImage)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DTOからの変換用ヘルパー
func (uc *CreateLineUserUseCase) ExecuteFromProfile(ctx context.Context, profile *dto.LineProfile) (*entity.User, error) {
	return uc.Execute(ctx, profile.UserID, profile.DisplayName, profile.PictureURL)
}