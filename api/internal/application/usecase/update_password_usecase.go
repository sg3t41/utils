package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/sg3t41/api/internal/domain/repository"
)

type UpdatePasswordUseCase struct {
	userRepository repository.UserRepository
}

func NewUpdatePasswordUseCase(userRepository repository.UserRepository) *UpdatePasswordUseCase {
	return &UpdatePasswordUseCase{
		userRepository: userRepository,
	}
}

type UpdatePasswordInput struct {
	UserID          string
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

type UpdatePasswordOutput struct {
	Success bool
}

var (
	ErrInvalidOldPassword    = errors.New("old password is incorrect")
	ErrPasswordMismatch      = errors.New("new password and confirm password do not match")
	ErrWeakPassword          = errors.New("password must be at least 8 characters long")
)

func (uc *UpdatePasswordUseCase) Execute(ctx context.Context, input UpdatePasswordInput) (*UpdatePasswordOutput, error) {
	if input.NewPassword != input.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	if len(input.NewPassword) < 8 {
		return nil, ErrWeakPassword
	}

	user, err := uc.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !uc.verifyPassword(input.OldPassword, user.Password) {
		return nil, ErrInvalidOldPassword
	}

	hashedPassword := uc.hashPassword(input.NewPassword)
	if err := user.UpdatePassword(hashedPassword); err != nil {
		return nil, err
	}

	err = uc.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &UpdatePasswordOutput{
		Success: true,
	}, nil
}

func (uc *UpdatePasswordUseCase) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (uc *UpdatePasswordUseCase) verifyPassword(password, hashedPassword string) bool {
	return uc.hashPassword(password) == hashedPassword
}