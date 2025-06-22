package usecase

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/service"
)

type CreateUserInput struct {
	Email string
	Name  string
}

type CreateUserOutput struct {
	User *entity.User
}

type CreateUserUseCase struct {
	userService *service.UserService
}

func NewCreateUserUseCase(userService *service.UserService) *CreateUserUseCase {
	return &CreateUserUseCase{
		userService: userService,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	user, err := uc.userService.CreateUser(ctx, input.Email, input.Name)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutput{
		User: user,
	}, nil
}