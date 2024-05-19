package user

import (
	"context"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/entity/transaction"
	transaction_repository "github.com/karamaru-alpha/days/pkg/domain/repository/transaction"
)

type Service interface {
	UpdateName(ctx context.Context, tx database.RWTx, userID, name string) error
}

type service struct {
	userRepository transaction_repository.UserRepository
}

func NewService(userRepository transaction_repository.UserRepository) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) UpdateName(ctx context.Context, tx database.RWTx, userID, name string) error {
	user, err := s.userRepository.LoadByPK(ctx, tx, &transaction.UserPK{
		UserID: userID,
	})
	if err != nil {
		return err
	}
	if user.Name == name {
		return nil
	}
	user.Name = name
	return s.userRepository.Update(ctx, tx, user)
}
