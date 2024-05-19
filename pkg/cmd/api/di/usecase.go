package di

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/cmd/api/usecase/user"
)

var usecaseSet = fx.Provide(
	user.NewUsecase,
)
