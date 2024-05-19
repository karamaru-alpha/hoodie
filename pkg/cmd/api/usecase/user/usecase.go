package user

import (
	"context"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/service/user"
)

type Usecase interface {
	UpdateName(ctx context.Context, userID, name string) error
}

type usecase struct {
	txManager   database.TransactionTxManager
	userService user.Service
}

func NewUsecase(txManager database.TransactionTxManager, userService user.Service) Usecase {
	return &usecase{
		txManager:   txManager,
		userService: userService,
	}
}

func (u *usecase) UpdateName(ctx context.Context, userID, name string) error {
	return u.txManager.Transaction(ctx, func(ctx context.Context, tx database.RWTx) error {
		return u.userService.UpdateName(ctx, tx, userID, name)
	})
}
