package di

import (
	"go.uber.org/fx"

	"github.com/karamaru-alpha/days/pkg/domain/service/user"
)

var serviceSet = fx.Provide(
	user.NewService,
)
