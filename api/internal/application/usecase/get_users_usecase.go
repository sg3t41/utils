package usecase

import (
	"context"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type GetUsersUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUsersUseCase(userRepository repository.UserRepository) *GetUsersUseCase {
	return &GetUsersUseCase{
		userRepository: userRepository,
	}
}

type GetUsersInput struct {
	Limit       int
	Page        int
	Cursor      string
	Sort        string
	Order       string
	Search      string
	Status      string
	CreatedFrom string
	CreatedTo   string
	UseCursor   bool
}

type GetUsersOutput struct {
	Users []*entity.User
	Total int
}

func (uc *GetUsersUseCase) Execute(ctx context.Context, input GetUsersInput) (*GetUsersOutput, error) {
	if input.Limit == 0 {
		users, err := uc.userRepository.FindAll(ctx)
		if err != nil {
			return nil, err
		}
		return &GetUsersOutput{
			Users: users,
			Total: len(users),
		}, nil
	}

	filter := repository.PaginationFilter{
		Search:      input.Search,
		Status:      input.Status,
		CreatedFrom: input.CreatedFrom,
		CreatedTo:   input.CreatedTo,
	}

	sort := repository.SortOption{
		Field: input.Sort,
		Order: input.Order,
	}

	if input.UseCursor {
		users, err := uc.userRepository.FindWithCursorPagination(ctx, input.Limit, input.Cursor, filter, sort)
		if err != nil {
			return nil, err
		}
		return &GetUsersOutput{
			Users: users,
			Total: 0,
		}, nil
	}

	offset := (input.Page - 1) * input.Limit
	result, err := uc.userRepository.FindWithOffsetPagination(ctx, input.Limit, offset, filter, sort)
	if err != nil {
		return nil, err
	}

	return &GetUsersOutput{
		Users: result.Users,
		Total: result.Total,
	}, nil
}