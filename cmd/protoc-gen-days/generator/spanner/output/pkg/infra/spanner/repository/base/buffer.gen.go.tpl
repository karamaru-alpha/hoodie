{{ template "autogen_comment" }}
{{ $goName := .GoName -}}
package base

import (
	"context"
	"sync"
	"time"

	"github.com/karamaru-alpha/days/pkg/domain/entity/{{ .PkgName }}"
	"github.com/karamaru-alpha/days/pkg/domain/enum"
	"github.com/karamaru-alpha/days/pkg/dcontext"
	"github.com/karamaru-alpha/days/pkg/derrors"
)

type {{ .GoName }}MutationWaitBuffer interface {
	Lock()
	Unlock()
	GetLatestEntity(sourceRow *{{ .PkgName }}.{{ .GoName }}) *{{ .PkgName }}.{{ .GoName }}
	MergeOperation(operation *{{ .GoName }}MutationWaitBufferRecordOperation) error
	GetState() map[string]*{{ .GoName }}MutationWaitBufferRecordOperation
}

type {{ .CamelName }}MutationWaitBufferKey struct{}

func Extract{{ .GoName }}MutationWaitBuffer(ctx context.Context) {{ .GoName }}MutationWaitBuffer {
	qctx := dcontext.ExtractQueryCache(ctx)
	if qctx.Nil() {
		return &{{ .CamelName }}MutationWaitBuffer{}
	}

	cacher, ok := qctx.GetCacher({{ .CamelName }}MutationWaitBufferKey{})
	if !ok {
		cacher = &{{ .CamelName }}MutationWaitBuffer{
			state: make(map[string]*{{ .GoName }}MutationWaitBufferRecordOperation),
		}
		qctx.SetCacher({{ .CamelName }}MutationWaitBufferKey{}, cacher)
	}
	return cacher.(*{{ .CamelName }}MutationWaitBuffer)
}

type {{ .GoName }}MutationWaitBufferRecordOperation struct {
	Type   			OperationType
	PK      		*{{ .PkgName }}.{{ .GoName }}PK
	ColumnValueMap  map[string]any
}

type {{ .CamelName }}MutationWaitBuffer struct {
	mutex sync.RWMutex
	state map[string]*{{ .GoName }}MutationWaitBufferRecordOperation
}

func (mwb *{{ .CamelName }}MutationWaitBuffer) Lock() {
	mwb.mutex.Lock()
}

func (mwb *{{ .CamelName }}MutationWaitBuffer) Unlock() {
	mwb.mutex.Unlock()
}

func (mwb *{{ .CamelName }}MutationWaitBuffer) GetLatestEntity(sourceRow *{{ .PkgName }}.{{ .GoName }}) *{{ .PkgName }}.{{ .GoName }} {
	operation, ok := mwb.state[sourceRow.GetPK().Key()]
	if !ok {
		return nil
	}

	result := sourceRow.FullDeepCopy()
	for key, value := range operation.ColumnValueMap {
		switch key {
		{{- range .Columns }}
		case transaction.{{ $goName }}ColumnName_{{ .GoName }}:
			result.{{ .GoName }} = value.({{ .Type }})
		{{- end }}
		}
	}

	return result
}

func (mwb *{{ .CamelName }}MutationWaitBuffer) MergeOperation(operation *{{ .GoName }}MutationWaitBufferRecordOperation) error {
	if operation.Type != OperationTypeDelete && len(operation.ColumnValueMap) == 0 {
		return nil
	}
	if operation.Type == OperationTypeDelete && len(operation.ColumnValueMap) > 0 {
		operation.ColumnValueMap = nil
	}

	key := operation.PK.Key()
	currentOperation, ok := mwb.state[key]
	if !ok {
		mwb.state[key] = operation
		return nil
	}

	switch currentOperation.Type {
	case OperationTypeInsert:
		switch operation.Type {
		case OperationTypeInsert:
			return derrors.New(derrors.InvalidArgument, "cannot insert because key is duplicated").SetValues(map[string]any{
				"tableName": {{ .GoName }}TableName,
				"pk":       currentOperation.PK,
			})
		case OperationTypeUpdate:
			for column, value := range operation.ColumnValueMap {
				currentOperation.ColumnValueMap[column] = value
			}
		case OperationTypeDelete:
			delete(mwb.state, key)
		}
	case OperationTypeUpdate:
		switch operation.Type {
		case OperationTypeInsert:
			return derrors.New(derrors.InvalidArgument, "cannot insert because key is duplicated").SetValues(map[string]any{
				"tableName": {{ .GoName }}TableName,
				"pk":	   currentOperation.PK,
			})
		case OperationTypeUpdate:
			for column, value := range operation.ColumnValueMap {
				currentOperation.ColumnValueMap[column] = value
			}
		case OperationTypeDelete:
			mwb.state[key] = operation
		}
	case OperationTypeDelete:
		switch operation.Type {
		case OperationTypeInsert:
			mwb.state[key] = operation
			mwb.state[key].Type = OperationTypeUpdate
		case OperationTypeUpdate:
			mwb.state[key] = operation
		case OperationTypeDelete:
			return nil
		}
	}

	return nil
}

func (mwb *{{ .CamelName }}MutationWaitBuffer) GetState() map[string]*{{ .GoName }}MutationWaitBufferRecordOperation {
	return mwb.state
}

func (o *{{ .GoName }}MutationWaitBufferRecordOperation) GetColumnsAndValuesWithPK() ([]string, []any) {
	pkMap := map[string]any{
		{{- range .PKColumns }}
		transaction.{{ $goName }}ColumnName_{{ .GoName }}: o.PK.{{ .GoName }},
		{{- end }}
	}
	columns := make([]string, 0, len(o.ColumnValueMap)+len(pkMap))
	values := make([]any, 0, len(o.ColumnValueMap)+len(pkMap))

	for _, col := range transaction.{{ .GoName }}Columns {
		if _, ok := pkMap[col.Name]; ok {
			columns = append(columns, col.Name)
			values = append(values, pkMap[col.Name])
			continue
		}
		if value, ok := o.ColumnValueMap[col.Name]; ok {
			columns = append(columns, col.Name)
			values = append(values, value)
		}
	}

	return columns, values
}
