package usecase

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/service"
)

type GetUserInput struct {
	ID string
}

type GetUserOutput struct {
	User *entity.User
}

type GetUserUseCase struct {
	userService *service.UserService
}

func NewGetUserUseCase(userService *service.UserService) *GetUserUseCase {
	return &GetUserUseCase{
		userService: userService,
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := uc.userService.GetUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{
		User: user,
	}, nil
}