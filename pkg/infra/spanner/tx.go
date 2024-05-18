package spanner

import (
	"context"

	"cloud.google.com/go/spanner"

	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/database"
)

// ReadTransaction common type of spanner.ReadWriteTransaction and spanner.ReadOnlyTransaction
type ReadTransaction interface {
	ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
	Read(ctx context.Context, table string, keys spanner.KeySet, columns []string) *spanner.RowIterator
	Query(ctx context.Context, statement spanner.Statement) *spanner.RowIterator
}

type roTx struct {
	*spanner.ReadOnlyTransaction
}

func (tx *roTx) ROTxImpl() {}

type rwTx struct {
	*spanner.ReadWriteTransaction
}

func (tx *rwTx) ROTxImpl() {}
func (tx *rwTx) RWTxImpl() {}

func ExtractROTx(tx database.ROTx) (ReadTransaction, error) {
	switch txObject := tx.(type) {
	case *roTx:
		return txObject, nil
	case *rwTx:
		return txObject, nil
	default:
		return nil, derrors.New(derrors.Internal, "txObject is neither spanner.ReadOnlyTransaction nor spanner.ReadWriteTransaction")
	}
}

func ExtractRWTx(tx database.RWTx) (*spanner.ReadWriteTransaction, error) {
	txObject, ok := tx.(*rwTx)
	if !ok {
		return nil, derrors.New(derrors.Internal, "txObject isn't spanner.ReadWriteTransaction")
	}
	return txObject.ReadWriteTransaction, nil
}
