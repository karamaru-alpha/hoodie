package database

import "context"

type TxManager interface {
	Transaction(ctx context.Context, f func(context.Context, RWTx) error) error
	ReadOnlyTransaction(ctx context.Context, f func(context.Context, ROTx) error) error
}

type TransactionTxManager TxManager
