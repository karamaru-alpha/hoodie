package spanner

import (
	"context"

	"cloud.google.com/go/spanner"

	"github.com/karamaru-alpha/days/pkg/domain/database"
)

func NewTxManager(client *spanner.Client) database.TxManager {
	return &txManager{
		client: client,
	}
}

type txManager struct {
	client *spanner.Client
}

func (m *txManager) Transaction(ctx context.Context, f func(context.Context, database.RWTx) error) error {
	if _, err := m.client.ReadWriteTransaction(ctx, func(ctx context.Context, spannerRWTx *spanner.ReadWriteTransaction) error {
		if err := f(ctx, &rwTx{spannerRWTx}); err != nil {
			return err
		}
		return bufferWrite(ctx, spannerRWTx)
	}); err != nil {
		return err
	}
	return nil
}

func (m *txManager) ReadOnlyTransaction(ctx context.Context, f func(context.Context, database.ROTx) error) error {
	spannerROTx := m.client.ReadOnlyTransaction()
	defer spannerROTx.Close()
	return f(ctx, &roTx{spannerROTx})
}
