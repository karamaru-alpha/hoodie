// Code generated by protoc-gen-days (generator/spanner) . DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package repository

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/scylladb/go-set/strset"
	"google.golang.org/grpc/codes"

	"github.com/karamaru-alpha/days/pkg/dcontext"
	"github.com/karamaru-alpha/days/pkg/derrors"
	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/entity/transaction"
	repository "github.com/karamaru-alpha/days/pkg/domain/repository/transaction"
	qspanner "github.com/karamaru-alpha/days/pkg/infra/spanner"
	"github.com/karamaru-alpha/days/pkg/infra/spanner/repository/base"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (r *userRepository) extractQueryCache(ctx context.Context) (base.UserSearchResultCache, base.UserMutationWaitBuffer) {
	return base.ExtractUserSearchResultCache(ctx), base.ExtractUserMutationWaitBuffer(ctx)
}

func (r *userRepository) LoadByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (*transaction.User, error) {
	row, err := r.SelectByPK(ctx, tx, pk)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, derrors.New(derrors.InvalidArgument, "record not found").SetValues(map[string]any{
			"tableName": base.UserTableName,
			"pk":        pk,
		})
	}

	return row, nil
}

func (r *userRepository) LoadByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (transaction.UserSlice, error) {
	rows, err := r.SelectByPKs(ctx, tx, pks)
	if err != nil {
		return nil, err
	}

	set := strset.NewWithSize(len(rows))
	for _, row := range rows {
		set.Add(row.GetPK().Key())
	}

	notFoundPKs := make(transaction.UserPKs, 0, len(pks))
	for _, pk := range pks {
		if !set.Has(pk.Key()) {
			notFoundPKs = append(notFoundPKs, pk)
		}
	}
	if len(notFoundPKs) > 0 {
		return nil, derrors.New(derrors.InvalidArgument, "record not found").SetValues(map[string]any{
			"tableName": base.UserTableName,
			"pks":       notFoundPKs,
		})
	}

	return rows, nil
}

func (r *userRepository) SelectByPK(ctx context.Context, tx database.ROTx, pk *transaction.UserPK) (result *transaction.User, err error) {
	searchResultCache, _ := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()

	if cachedEntity, resultType := searchResultCache.GetByPK(pk); resultType != base.SearchResultTypeNotSearched {
		return cachedEntity, nil
	}

	roTx, err := qspanner.ExtractROTx(tx)
	if err != nil {
		return nil, err
	}

	row, err := roTx.ReadRow(ctx, base.UserTableName, spanner.Key(pk.Generate()), base.UserColumnNames)
	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			searchResultCache.SetAsNotFound(pk)
			return nil, nil
		}
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.UserTableName,
			"pk":        pk,
		})
	}

	result, err = r.decodeAllColumns(row)
	if err != nil {
		return nil, err
	}

	searchResultCache.SetResult(result)

	return result, nil
}

func (r *userRepository) SelectByPKs(ctx context.Context, tx database.ROTx, pks transaction.UserPKs) (rows transaction.UserSlice, err error) {
	if len(pks) == 0 {
		return transaction.UserSlice{}, nil
	}

	searchResultCache, _ := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	if cachedEntity, resultType := searchResultCache.GetByPKs(pks); cachedEntity != nil && resultType != base.SearchResultTypeNotSearched {
		return cachedEntity, nil
	}

	roTx, err := qspanner.ExtractROTx(tx)
	if err != nil {
		return nil, err
	}

	keySets := make([]spanner.KeySet, 0, len(pks))
	for _, pk := range pks {
		if _, resultType := searchResultCache.GetByPK(pk); resultType != base.SearchResultTypeNotSearched {
			continue
		}
		keySets = append(keySets, spanner.Key(pk.Generate()))
	}
	ri := roTx.Read(ctx, base.UserTableName, spanner.KeySets(keySets...), base.UserColumnNames)
	rows = make(transaction.UserSlice, 0)
	keySet := strset.New()
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.UserSlice, 0, ri.RowCount)
			keySet = strset.NewWithSize(int(ri.RowCount))
		}
		result, err := r.decodeAllColumns(row)
		if err != nil {
			return err
		}
		rows = append(rows, result)
		keySet.Add(result.GetPK().Key())
		return nil
	}); err != nil {
		if err, ok := derrors.As(err); ok {
			return nil, err
		}
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.UserTableName,
			"pks":       pks,
		})
	}

	rows = searchResultCache.ReplaceWithCache(rows)
	searchResultCache.SetResults(rows)
	for _, pk := range pks {
		if cachedEntity, resultType := searchResultCache.GetByPK(pk); resultType == base.SearchResultTypeNotSearched {
			searchResultCache.SetAsNotFound(pk)
		} else if !keySet.Has(pk.Key()) && cachedEntity != nil {
			rows = append(rows, cachedEntity)
		}
	}

	return rows, nil
}

func (r *userRepository) SelectAll(ctx context.Context, tx database.ROTx, limit, offset int) (rows transaction.UserSlice, err error) {
	roTx, err := qspanner.ExtractROTx(tx)
	if err != nil {
		return nil, err
	}

	sql, params := base.NewUserQueryBuilder().
		SelectAll().
		OrderBy(base.OrderPairs{{Column: base.UserColumnNameUserID, OrderType: base.OrderTypeASC}}).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL:    sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make(transaction.UserSlice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make(transaction.UserSlice, 0, ri.RowCount)
		}
		result, err := r.decodeAllColumns(row)
		if err != nil {
			return err
		}
		rows = append(rows, result)
		return nil
	}); err != nil {
		if err, ok := derrors.As(err); ok {
			return nil, err
		}
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.UserTableName,
			"limit":     limit,
			"offset":    offset,
		})
	}

	return rows, nil
}

func (r *userRepository) Insert(ctx context.Context, tx database.RWTx, row *transaction.User) (err error) {
	now := dcontext.Now(ctx)
	row.CreatedTime = now
	row.UpdatedTime = now

	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	mutationWaitBuffer.Lock()
	defer mutationWaitBuffer.Unlock()

	if err := mutationWaitBuffer.MergeOperation(&base.UserMutationWaitBufferRecordOperation{
		Type:           base.OperationTypeInsert,
		PK:             row.GetPK(),
		ColumnValueMap: row.ToKeyValue(),
	}); err != nil {
		return err
	}

	searchResultCache.SetResult(row)

	return nil
}

func (r *userRepository) BulkInsert(ctx context.Context, tx database.RWTx, rows transaction.UserSlice) (err error) {
	if len(rows) == 0 {
		return nil
	}

	now := dcontext.Now(ctx)
	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	mutationWaitBuffer.Lock()
	defer mutationWaitBuffer.Unlock()

	for _, row := range rows {
		row.CreatedTime = now
		row.UpdatedTime = now
		if err := mutationWaitBuffer.MergeOperation(&base.UserMutationWaitBufferRecordOperation{
			Type:           base.OperationTypeInsert,
			PK:             row.GetPK(),
			ColumnValueMap: row.ToKeyValue(),
		}); err != nil {
			return err
		}
	}

	searchResultCache.SetResults(rows)

	return nil
}

func (r *userRepository) Update(ctx context.Context, tx database.RWTx, row *transaction.User) (err error) {
	now := dcontext.Now(ctx)
	row.UpdatedTime = now
	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)

	pk := row.GetPK()
	originRow, resultType := searchResultCache.GetOriginByPK(pk)
	sourceRow := originRow
	if sourceRow == nil {
		sourceRow = pk.ToEntity()
	}
	updatedRow := mutationWaitBuffer.GetLatestEntity(sourceRow)
	var cachedRow *transaction.User
	switch {
	case updatedRow != nil:
		cachedRow = updatedRow
	case originRow != nil:
		cachedRow = originRow
	default:
		switch resultType {
		case base.SearchResultTypeNotSearched:
			return derrors.New(derrors.InvalidArgument, "cannot update record not selected").SetValues(map[string]any{
				"tableName": base.UserTableName,
				"pk":        pk,
			})
		case base.SearchResultTypeNotFound:
			return derrors.New(derrors.InvalidArgument, "cannot update record not found").SetValues(map[string]any{
				"tableName": base.UserTableName,
				"pk":        pk,
			})
		default:
			return derrors.New(derrors.InvalidArgument, "invalid cache").SetValues(map[string]any{
				"tableName": base.UserTableName,
				"pk":        pk,
			})
		}
	}

	if err := mutationWaitBuffer.MergeOperation(&base.UserMutationWaitBufferRecordOperation{
		Type:           base.OperationTypeUpdate,
		PK:             pk,
		ColumnValueMap: r.diffEntity(cachedRow, row),
	}); err != nil {
		return err
	}

	searchResultCache.SetResult(row)

	return nil
}

func (r *userRepository) Save(ctx context.Context, tx database.RWTx, row *transaction.User) error {
	var err error
	if row.CreatedTime.IsZero() {
		err = r.Insert(ctx, tx, row)
	} else {
		err = r.Update(ctx, tx, row)
	}

	return err
}

func (r *userRepository) Delete(ctx context.Context, tx database.RWTx, pk *transaction.UserPK) (err error) {
	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	mutationWaitBuffer.Lock()
	defer mutationWaitBuffer.Unlock()

	switch _, resultType := searchResultCache.GetByPK(pk); resultType {
	case base.SearchResultTypeNotSearched:
		// allow for the implementation using interleaving and to maintain backward compatibility
	case base.SearchResultTypeNotFound:
		// return early because record not found
		return nil
	}

	if err := mutationWaitBuffer.MergeOperation(&base.UserMutationWaitBufferRecordOperation{
		Type: base.OperationTypeDelete,
		PK:   pk,
	}); err != nil {
		return err
	}

	searchResultCache.SetAsNotFound(pk)

	return nil
}

func (r *userRepository) BulkDelete(ctx context.Context, tx database.RWTx, pks transaction.UserPKs) (err error) {
	if len(pks) == 0 {
		return nil
	}

	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	mutationWaitBuffer.Lock()
	defer mutationWaitBuffer.Unlock()

	for _, pk := range pks {
		switch _, resultType := searchResultCache.GetByPK(pk); resultType {
		case base.SearchResultTypeNotSearched:
			// allow for the implementation using interleaving and to maintain backward compatibility
		case base.SearchResultTypeNotFound:
			// return early because record not found
			continue
		}

		if err := mutationWaitBuffer.MergeOperation(&base.UserMutationWaitBufferRecordOperation{
			Type: base.OperationTypeDelete,
			PK:   pk,
		}); err != nil {
			return err
		}

		searchResultCache.SetAsNotFound(pk)
	}

	return nil
}

func (r *userRepository) decodeAllColumns(row *spanner.Row) (*transaction.User, error) {
	var userID spanner.NullString
	var name spanner.NullString
	var createdTime spanner.NullTime
	var updatedTime spanner.NullTime

	if err := row.Columns(
		&userID,
		&name,
		&createdTime,
		&updatedTime,
	); err != nil {
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.UserTableName,
			"key":       "All",
		})
	}

	var result transaction.User
	if !userID.IsNull() {
		result.UserID = userID.StringVal
	}
	if !name.IsNull() {
		result.Name = name.StringVal
	}
	if !createdTime.IsNull() {
		result.CreatedTime = createdTime.Time.In(time.Local)
	}
	if !updatedTime.IsNull() {
		result.UpdatedTime = updatedTime.Time.In(time.Local)
	}
	return &result, nil
}

func (r *userRepository) diffEntity(source, target *transaction.User) map[string]any {
	result := make(map[string]any, 0)

	if source.Name != target.Name {
		result[transaction.UserColumnName_Name] = target.Name
	}
	if !source.CreatedTime.Equal(target.CreatedTime) {
		result[transaction.UserColumnName_CreatedTime] = target.CreatedTime
	}
	if !source.UpdatedTime.Equal(target.UpdatedTime) {
		result[transaction.UserColumnName_UpdatedTime] = target.UpdatedTime
	}

	return result
}