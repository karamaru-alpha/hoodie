{{ template "autogen_comment" }}

package spanner

import (
	"context"

	"cloud.google.com/go/spanner"

	"github.com/karamaru-alpha/days/pkg/domain/entity/transaction"
	"github.com/karamaru-alpha/days/pkg/infra/spanner/repository/base"
	"github.com/karamaru-alpha/days/pkg/derrors"
)

func bufferWrite(ctx context.Context, spannerRWTx *spanner.ReadWriteTransaction) error {
	mutations := make([]*spanner.Mutation, 0)
	{{- range . }}
	for _, operation := range base.Extract{{ .GoName }}MutationWaitBuffer(ctx).GetState() {
		columns, values := operation.GetColumnsAndValuesWithPK()
		switch operation.Type {
		case base.OperationTypeInsert:
			mutations = append(mutations, spanner.Insert(transaction.{{ .GoName }}TableName, columns, values))
		case base.OperationTypeUpdate:
			mutations = append(mutations, spanner.Update(transaction.{{ .GoName }}TableName, columns, values))
		case base.OperationTypeDelete:
			mutations = append(mutations, spanner.Delete(transaction.{{ .GoName }}TableName, spanner.Key(values)))
		}
	}
	{{- end }}
	if err := spannerRWTx.BufferWrite(mutations); err != nil {
		return derrors.Wrap(err, derrors.Internal, "fail to buffer write to spanner")
	}

	return nil
}
