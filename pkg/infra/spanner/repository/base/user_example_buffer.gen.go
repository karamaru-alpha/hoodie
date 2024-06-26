// Code generated by protoc-gen-days (generator/spanner) . DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package base

import (
	"context"
	"sync"
	"time"

	"github.com/karamaru-alpha/days/pkg/dcontext"
	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/entity/transaction"
)

type UserExampleMutationWaitBuffer interface {
	Lock()
	Unlock()
	GetLatestEntity(sourceRow *transaction.UserExample) *transaction.UserExample
	MergeOperation(operation *UserExampleMutationWaitBufferRecordOperation) error
	GetState() map[string]*UserExampleMutationWaitBufferRecordOperation
}

type userExampleMutationWaitBufferKey struct{}

func ExtractUserExampleMutationWaitBuffer(ctx context.Context) UserExampleMutationWaitBuffer {
	qctx := dcontext.ExtractQueryCache(ctx)
	if qctx.Nil() {
		return &userExampleMutationWaitBuffer{}
	}

	cacher, ok := qctx.GetCacher(userExampleMutationWaitBufferKey{})
	if !ok {
		cacher = &userExampleMutationWaitBuffer{
			state: make(map[string]*UserExampleMutationWaitBufferRecordOperation),
		}
		qctx.SetCacher(userExampleMutationWaitBufferKey{}, cacher)
	}
	return cacher.(*userExampleMutationWaitBuffer)
}

type UserExampleMutationWaitBufferRecordOperation struct {
	Type           OperationType
	PK             *transaction.UserExamplePK
	ColumnValueMap map[string]any
}

type userExampleMutationWaitBuffer struct {
	mutex sync.RWMutex
	state map[string]*UserExampleMutationWaitBufferRecordOperation
}

func (mwb *userExampleMutationWaitBuffer) Lock() {
	mwb.mutex.Lock()
}

func (mwb *userExampleMutationWaitBuffer) Unlock() {
	mwb.mutex.Unlock()
}

func (mwb *userExampleMutationWaitBuffer) GetLatestEntity(sourceRow *transaction.UserExample) *transaction.UserExample {
	operation, ok := mwb.state[sourceRow.GetPK().Key()]
	if !ok {
		return nil
	}

	result := sourceRow.FullDeepCopy()
	for key, value := range operation.ColumnValueMap {
		switch key {
		case transaction.UserExampleColumnName_UserID:
			result.UserID = value.(string)
		case transaction.UserExampleColumnName_Example:
			result.Example = value.(int64)
		case transaction.UserExampleColumnName_CreatedTime:
			result.CreatedTime = value.(time.Time)
		case transaction.UserExampleColumnName_UpdatedTime:
			result.UpdatedTime = value.(time.Time)
		}
	}

	return result
}

func (mwb *userExampleMutationWaitBuffer) MergeOperation(operation *UserExampleMutationWaitBufferRecordOperation) error {
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
				"tableName": UserExampleTableName,
				"pk":        currentOperation.PK,
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
				"tableName": UserExampleTableName,
				"pk":        currentOperation.PK,
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

func (mwb *userExampleMutationWaitBuffer) GetState() map[string]*UserExampleMutationWaitBufferRecordOperation {
	return mwb.state
}

func (o *UserExampleMutationWaitBufferRecordOperation) GetColumnsAndValuesWithPK() ([]string, []any) {
	pkMap := map[string]any{
		transaction.UserExampleColumnName_UserID:  o.PK.UserID,
		transaction.UserExampleColumnName_Example: o.PK.Example,
	}
	columns := make([]string, 0, len(o.ColumnValueMap)+len(pkMap))
	values := make([]any, 0, len(o.ColumnValueMap)+len(pkMap))

	for _, col := range transaction.UserExampleColumns {
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
