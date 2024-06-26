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

type UserMutationWaitBuffer interface {
	Lock()
	Unlock()
	GetLatestEntity(sourceRow *transaction.User) *transaction.User
	MergeOperation(operation *UserMutationWaitBufferRecordOperation) error
	GetState() map[string]*UserMutationWaitBufferRecordOperation
}

type userMutationWaitBufferKey struct{}

func ExtractUserMutationWaitBuffer(ctx context.Context) UserMutationWaitBuffer {
	qctx := dcontext.ExtractQueryCache(ctx)
	if qctx.Nil() {
		return &userMutationWaitBuffer{}
	}

	cacher, ok := qctx.GetCacher(userMutationWaitBufferKey{})
	if !ok {
		cacher = &userMutationWaitBuffer{
			state: make(map[string]*UserMutationWaitBufferRecordOperation),
		}
		qctx.SetCacher(userMutationWaitBufferKey{}, cacher)
	}
	return cacher.(*userMutationWaitBuffer)
}

type UserMutationWaitBufferRecordOperation struct {
	Type           OperationType
	PK             *transaction.UserPK
	ColumnValueMap map[string]any
}

type userMutationWaitBuffer struct {
	mutex sync.RWMutex
	state map[string]*UserMutationWaitBufferRecordOperation
}

func (mwb *userMutationWaitBuffer) Lock() {
	mwb.mutex.Lock()
}

func (mwb *userMutationWaitBuffer) Unlock() {
	mwb.mutex.Unlock()
}

func (mwb *userMutationWaitBuffer) GetLatestEntity(sourceRow *transaction.User) *transaction.User {
	operation, ok := mwb.state[sourceRow.GetPK().Key()]
	if !ok {
		return nil
	}

	result := sourceRow.FullDeepCopy()
	for key, value := range operation.ColumnValueMap {
		switch key {
		case transaction.UserColumnName_UserID:
			result.UserID = value.(string)
		case transaction.UserColumnName_Name:
			result.Name = value.(string)
		case transaction.UserColumnName_CreatedTime:
			result.CreatedTime = value.(time.Time)
		case transaction.UserColumnName_UpdatedTime:
			result.UpdatedTime = value.(time.Time)
		}
	}

	return result
}

func (mwb *userMutationWaitBuffer) MergeOperation(operation *UserMutationWaitBufferRecordOperation) error {
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
				"tableName": UserTableName,
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
				"tableName": UserTableName,
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

func (mwb *userMutationWaitBuffer) GetState() map[string]*UserMutationWaitBufferRecordOperation {
	return mwb.state
}

func (o *UserMutationWaitBufferRecordOperation) GetColumnsAndValuesWithPK() ([]string, []any) {
	pkMap := map[string]any{
		transaction.UserColumnName_UserID: o.PK.UserID,
	}
	columns := make([]string, 0, len(o.ColumnValueMap)+len(pkMap))
	values := make([]any, 0, len(o.ColumnValueMap)+len(pkMap))

	for _, col := range transaction.UserColumns {
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
