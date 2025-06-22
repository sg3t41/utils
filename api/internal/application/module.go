package application

import (
	"github.com/sg3t41/api/internal/application/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("application",
	fx.Provide(
		usecase.NewCreateUserUseCase,
		usecase.NewGetUserUseCase,
		usecase.NewGetUsersUseCase,
		usecase.NewUpdateUserUseCase,
		usecase.NewUpdatePasswordUseCase,
		usecase.NewDeleteUserUseCase,
		usecase.NewCreateArticleUseCase,
	),
)