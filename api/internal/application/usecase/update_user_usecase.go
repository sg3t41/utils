package usecase

import (
	"context"
	"errors"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type UpdateUserUseCase struct {
	userRepository repository.UserRepository
}

func NewUpdateUserUseCase(userRepository repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepository: userRepository,
	}
}

type UpdateUserInput struct {
	ID      string
	Name    *string
	Email   *string
	Version int
}

type UpdateUserOutput struct {
	User *entity.User
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email is already taken")
	ErrVersionConflict   = errors.New("version conflict - user has been updated by another process")
)

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UpdateUserOutput, error) {
	user, err := uc.userRepository.FindByID(ctx, input.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if user.Version != input.Version {
		return nil, ErrVersionConflict
	}

	if input.Email != nil && *input.Email != user.Email {
		existingUser, err := uc.userRepository.FindByEmail(ctx, *input.Email)
		if err == nil && existingUser.ID != user.ID {
			return nil, ErrEmailAlreadyTaken
		}
	}

	if err := user.Update(input.Name, input.Email); err != nil {
		return nil, err
	}

	err = uc.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutput{
		User: user,
	}, nil
}