// Code generated by protoc-gen-days (generator/spanner) . DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package transaction

import (
	"context"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/entity/transaction"
)

type UserRepository interface {
	LoadByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (*transaction.User, error)
	LoadByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (transaction.UserSlice, error)
	SelectByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (*transaction.User, error)
	SelectByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (transaction.UserSlice, error)
	SelectAll(ctx context.Context, tx database.ROTx, limit int, offset int) (transaction.UserSlice, error)
	Insert(ctx context.Context, tx database.RWTx, row *transaction.User) error
	BulkInsert(ctx context.Context, tx database.RWTx, rows transaction.UserSlice) error
	Update(ctx context.Context, tx database.RWTx, row *transaction.User) error
	Save(ctx context.Context, tx database.RWTx, row *transaction.User) error
	Delete(ctx context.Context, tx database.RWTx, pk *transaction.UserPK) error
	BulkDelete(ctx context.Context, tx database.RWTx, pks transaction.UserPKs) error
}