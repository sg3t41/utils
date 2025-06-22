package usecase

import (
	"context"
	"errors"

	"github.com/sg3t41/api/internal/domain/repository"
)

type DeleteUserUseCase struct {
	userRepository repository.UserRepository
}

func NewDeleteUserUseCase(userRepository repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

type DeleteUserInput struct {
	ID   string
	Hard bool
}

type DeleteUserOutput struct{}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	if input.ID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.userRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if user.IsDeleted() {
		return nil, errors.New("user is already deleted")
	}

	if input.Hard {
		err = uc.userRepository.HardDelete(ctx, input.ID)
	} else {
		err = uc.userRepository.SoftDelete(ctx, input.ID)
	}

	if err != nil {
		return nil, err
	}

	return &DeleteUserOutput{}, nil
}