{{ template "autogen_comment" -}}
{{- $goName := .GoName -}}
{{- $camelName := .CamelName -}}
{{- $pkgName := .PkgName -}}

package repository

import (
	"bytes"
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/scylladb/go-set/strset"
	"google.golang.org/grpc/codes"

	"github.com/karamaru-alpha/days/pkg/domain/database"
	"github.com/karamaru-alpha/days/pkg/domain/entity/{{ .PkgName }}"
	"github.com/karamaru-alpha/days/pkg/domain/enum"
	repository "github.com/karamaru-alpha/days/pkg/domain/repository/{{ .PkgName }}"
	qspanner "github.com/karamaru-alpha/days/pkg/infra/spanner"
	"github.com/karamaru-alpha/days/pkg/infra/spanner/repository/base"
	"github.com/karamaru-alpha/days/pkg/dcontext"
	"github.com/karamaru-alpha/days/pkg/derrors"
)

type {{ .CamelName }}Repository struct {}

func New{{ .GoName }}Repository() repository.{{ .GoName }}Repository {
	return &{{ .CamelName }}Repository{}
}

func (r *{{ .CamelName }}Repository) extractQueryCache(ctx context.Context) (base.{{ .GoName }}SearchResultCache, base.{{ .GoName }}MutationWaitBuffer) {
	return base.Extract{{ .GoName }}SearchResultCache(ctx), base.Extract{{ .GoName }}MutationWaitBuffer(ctx)
}

func (r *{{ .CamelName }}Repository) LoadByPK(ctx context.Context, tx database.ROTx, pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, error) {
	row, err := r.SelectByPK(ctx, tx, pk)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, derrors.New(derrors.InvalidArgument, "record not found").SetValues(map[string]any{
			"tableName": base.{{ .GoName }}TableName,
			"pk":        pk,
		})
	}

	return row, nil
}

func (r *{{ .CamelName }}Repository) LoadByPKs(ctx context.Context, tx database.ROTx, pks {{ .PkgName }}.{{ .GoName }}PKs) ({{ .PkgName }}.{{ .GoName }}Slice, error) {
	rows, err := r.SelectByPKs(ctx, tx, pks)
	if err != nil {
		return nil, err
	}

	set := strset.NewWithSize(len(rows))
	for _, row := range rows {
		set.Add(row.GetPK().Key())
	}

	notFoundPKs := make({{ .PkgName }}.{{ .GoName }}PKs, 0, len(pks))
	for _, pk := range pks {
		if !set.Has(pk.Key()) {
			notFoundPKs = append(notFoundPKs, pk)
		}
	}
	if len(notFoundPKs) > 0 {
		return nil, derrors.New(derrors.InvalidArgument, "record not found").SetValues(map[string]any{
			"tableName": base.{{ .GoName }}TableName,
			"pks":       notFoundPKs,
		})
	}

	return rows, nil
}

func (r *{{ .CamelName }}Repository) SelectByPK(ctx context.Context, tx database.ROTx, pk *{{ .PkgName }}.{{ .GoName }}PK) (result *{{ .PkgName }}.{{ .GoName }}, err error) {
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

	row, err := roTx.ReadRow(ctx, base.{{ .GoName }}TableName, spanner.Key(pk.Generate()), base.{{ .GoName }}ColumnNames)
	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			searchResultCache.SetAsNotFound(pk)
			return nil, nil
		}
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
        	"tableName": base.{{ .GoName }}TableName,
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

func (r *{{ .CamelName }}Repository) SelectByPKs(ctx context.Context, tx database.ROTx, pks {{ .PkgName }}.{{ .GoName }}PKs) (rows {{ .PkgName }}.{{ .GoName }}Slice, err error) {
	if len(pks) == 0 {
		return {{ .PkgName }}.{{ .GoName }}Slice{}, nil
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
	ri := roTx.Read(ctx, base.{{ .GoName }}TableName, spanner.KeySets(keySets...), base.{{ .GoName }}ColumnNames)
	rows = make({{ .PkgName }}.{{ .GoName }}Slice, 0)
	keySet := strset.New()
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make({{ .PkgName }}.{{ .GoName }}Slice, 0, ri.RowCount)
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
			"tableName": base.{{ .GoName }}TableName,
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

func (r *{{ .CamelName }}Repository) SelectAll(ctx context.Context, tx database.ROTx, limit, offset int) (rows {{ .PkgName }}.{{ .GoName }}Slice, err error) {
	roTx, err := qspanner.ExtractROTx(tx)
	if err != nil {
		return nil, err
	}

	sql, params := base.New{{ .GoName }}QueryBuilder().
		SelectAll().
		OrderBy(base.OrderPairs{ {{- range $i, $col := .PKColumns }}{{ if $i }}, {{ end }}{Column: base.{{ $goName }}ColumnName{{ $col.GoName }}, OrderType: base.OrderTypeASC}{{ end -}} }).
		Limit(limit).
		Offset(offset).
		GetQuery()
	stmt := spanner.Statement{
		SQL: sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make({{ .PkgName }}.{{ .GoName }}Slice, 0)
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make({{ .PkgName }}.{{ .GoName }}Slice, 0, ri.RowCount)
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
			"tableName": base.{{ .GoName }}TableName,
			"limit": limit,
			"offset": offset,
		})
	}

	return rows, nil
}
{{- range .Methods }}

func (r *{{ $camelName }}Repository) SelectBy{{ .Name }}(ctx context.Context, tx database.ROTx, {{ .Args }}) (rows {{ $pkgName }}.{{ .ReturnName }}Slice, err error) {
	{{ if .SliceArgName }}
	if len({{ .SliceArgName }}) == 0 {
		return {{ $pkgName }}.{{ .ReturnName }}Slice{}, nil
	}

	{{ end -}}

	roTx, err := qspanner.ExtractROTx(tx)
	if err != nil {
		return nil, err
	}

	qb := base.New{{ $goName }}QueryBuilder().
		Select{{ .SelectType }}().
		Where().{{ .Wheres }}

	{{ if .UseCache -}}
	searchResultCache, _ := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	if cachedEntities, resultType := searchResultCache.GetByQueryConditions(qb.GetQueryConditions()); resultType != base.SearchResultTypeNotSearched {
		return cachedEntities, nil
	}
	{{- end }}

	sql, params := qb.GetQuery()
	stmt := spanner.Statement{
		SQL: sql,
		Params: params,
	}
	ri := roTx.Query(ctx, stmt)

	rows = make({{ $pkgName }}.{{ .ReturnName }}Slice, 0)
	{{ if .UseCache }}keySet := strset.New(){{ end }}
	if err := ri.Do(func(row *spanner.Row) error {
		if len(rows) == 0 {
			rows = make({{ $pkgName }}.{{ .ReturnName }}Slice, 0, ri.RowCount)
			{{ if .UseCache }}keySet = strset.NewWithSize(int(ri.RowCount)){{ end }}
		}
		result, err := r.decode{{ .SelectType }}Columns(row)
		if err != nil {
			return err
		}
		rows = append(rows, result)
		{{ if .UseCache }}keySet.Add(result.GetPK().Key()){{ end }}
		return nil
	}); err != nil {
		if err, ok := derrors.As(err); ok {
			return nil, err
		}
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.{{ $goName }}TableName,
		})
	}

	{{ if .UseCache -}}
	rows = searchResultCache.ReplaceWithCache(rows)
	searchResultCache.SetResults(rows)
	searchResultCache.AppendQueryConditions(qb.GetQueryConditions())

	for _, value := range searchResultCache.FilterByQueryConditions(qb.GetQueryConditions()) {
		if !keySet.Has(value.GetPK().Key()) {
			rows = append(rows, value)
		}
	}
	{{- end }}

	return rows, nil
}
{{- end }}

func (r *{{ .CamelName }}Repository) Insert(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) (err error) {
	now := dcontext.Now(ctx)
	row.CreatedTime = now
	row.UpdatedTime = now

	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)
	searchResultCache.Lock()
	defer searchResultCache.Unlock()
	mutationWaitBuffer.Lock()
	defer mutationWaitBuffer.Unlock()

	if err := mutationWaitBuffer.MergeOperation(&base.{{ .GoName }}MutationWaitBufferRecordOperation{
		Type:   		base.OperationTypeInsert,
		PK:		     	row.GetPK(),
		ColumnValueMap: row.ToKeyValue(),
	}); err != nil {
		return err
	}

	{{- if .NeedCommonResponse }}
	if err := dcontext.ExtractDataChange(ctx).AddUpdated{{ .GoName }}(ctx, row); err != nil {
		return err
	}
	{{- end }}

	searchResultCache.SetResult(row)

	return nil
}

func (r *{{ .CamelName }}Repository) BulkInsert(ctx context.Context, tx database.RWTx, rows {{ .PkgName }}.{{ .GoName }}Slice) (err error) {
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
		if err := mutationWaitBuffer.MergeOperation(&base.{{ .GoName }}MutationWaitBufferRecordOperation{
			Type:    		base.OperationTypeInsert,
			PK:      		row.GetPK(),
			ColumnValueMap: row.ToKeyValue(),
		}); err != nil {
			return err
		}

		{{- if .NeedCommonResponse }}
		if err := dcontext.ExtractDataChange(ctx).AddUpdated{{ .GoName }}(ctx, row); err != nil {
			return err
		}
		{{- end }}
	}

	searchResultCache.SetResults(rows)

	return nil
}

func (r *{{ .CamelName }}Repository) Update(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) (err error) {
	now := dcontext.Now(ctx)
	row.UpdatedTime = now
	searchResultCache, mutationWaitBuffer := r.extractQueryCache(ctx)

	pk := row.GetPK()
	originRow, resultType := searchResultCache.GetOriginByPK(pk)
	sourceRow := originRow
	if sourceRow == nil {
		sourceRow = pk.ToEntity()
	}
	updatedRow := mutationWaitBuffer.GetLatestEntity(sourceRow);
	var cachedRow *{{ .PkgName }}.{{ .GoName }}
	switch {
	case updatedRow != nil:
		cachedRow = updatedRow
	case originRow != nil:
		cachedRow = originRow
	default:
		switch resultType {
		case base.SearchResultTypeNotSearched:
			return derrors.New(derrors.InvalidArgument, "cannot update record not selected").SetValues(map[string]any{
				"tableName": base.{{ .GoName }}TableName,
				"pk":        pk,
			})
		case base.SearchResultTypeNotFound:
			return derrors.New(derrors.InvalidArgument, "cannot update record not found").SetValues(map[string]any{
				"tableName": base.{{ .GoName }}TableName,
				"pk":        pk,
			})
		default:
			return derrors.New(derrors.InvalidArgument, "invalid cache").SetValues(map[string]any{
				"tableName": base.{{ .GoName }}TableName,
				"pk":        pk,
			})
		}
	}

	if err := mutationWaitBuffer.MergeOperation(&base.{{ .GoName }}MutationWaitBufferRecordOperation{
		Type:  		 	base.OperationTypeUpdate,
		PK:   		 	pk,
		ColumnValueMap: r.diffEntity(cachedRow, row),
	}); err != nil {
		return err
	}

{{- if .NeedCommonResponse }}
	if err := dcontext.ExtractDataChange(ctx).AddUpdated{{ .GoName }}(ctx, row); err != nil {
		return err
	}
{{- end }}

	searchResultCache.SetResult(row)

	return nil
}

func (r *{{ .CamelName }}Repository) Save(ctx context.Context, tx database.RWTx, row *{{ .PkgName }}.{{ .GoName }}) error {
	var err error
	if row.CreatedTime.IsZero() {
		err = r.Insert(ctx, tx, row)
	} else {
		err = r.Update(ctx, tx, row)
	}

	return err
}

func (r *{{ .CamelName }}Repository) Delete(ctx context.Context, tx database.RWTx, pk *{{ .PkgName }}.{{ .GoName }}PK) (err error) {
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

	if err := mutationWaitBuffer.MergeOperation(&base.{{ .GoName }}MutationWaitBufferRecordOperation{
		Type: base.OperationTypeDelete,
		PK:   pk,
	}); err != nil {
		return err
	}

{{- if and .NeedCommonResponse (ge (len .PKColumns) 2) }}
	if err := dcontext.ExtractDataChange(ctx).AddDeleted{{ .GoName }}(ctx, pk); err != nil {
		return err
	}
{{- end }}

	searchResultCache.SetAsNotFound(pk)

	return nil
}

func (r *{{ .CamelName }}Repository) BulkDelete(ctx context.Context, tx database.RWTx, pks {{ .PkgName }}.{{ .GoName }}PKs) (err error) {
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

		if err := mutationWaitBuffer.MergeOperation(&base.{{ .GoName }}MutationWaitBufferRecordOperation{
			Type: base.OperationTypeDelete,
			PK:   pk,
		}); err != nil {
			return err
		}
		{{- if and .NeedCommonResponse (ge (len .PKColumns) 2) }}
		if err := dcontext.ExtractDataChange(ctx).AddDeleted{{ .GoName }}(ctx, pk); err != nil {
			return err
		}
		{{- end }}

		searchResultCache.SetAsNotFound(pk)
	}

	return nil
}
{{- range .Types }}

func (r *{{ $camelName }}Repository) decode{{ .Key }}Columns(row *spanner.Row) (*{{ $pkgName }}.{{ .GoName }}, error) {
	{{- range .Columns }}
	{{ if .IsEnum -}}
	var {{ .LocalName }} {{ if and .IsList }}[]{{ end }}spanner.NullInt64
	{{- else if .IsList -}}
	var {{ .LocalName }} []{{ .Type }}
	{{- else if eq .Type "time.Time" -}}
	var {{ .LocalName }} spanner.NullTime
	{{- else -}}
	var {{ .LocalName }} spanner.Null{{ title .Type }}
	{{- end -}}
	{{- end }}

	if err := row.Columns(
		{{ range .Columns -}}
		&{{ .LocalName }},
		{{ end -}}
	); err != nil {
		return nil, derrors.Wrap(err, derrors.Internal, err.Error()).SetValues(map[string]any{
			"tableName": base.{{ $goName }}TableName,
			"key":       "{{ .Key }}",
		})
	}

	var result {{ $pkgName }}.{{ .GoName }}
	{{- range .Columns }}
	{{ if and .IsList .IsEnum -}}
	{
		s := make({{ .Type }}Slice, 0, len({{ .LocalName }}))
		for _, v := range {{ .LocalName }} {
			s = append(s, {{ .Type }}(v.Int64))
		}
		result.{{ .GoName }} = s
	}
	{{- else if .IsList -}}
	result.{{ .GoName }} = {{ .LocalName }}
	{{- else -}}
	if !{{ .LocalName }}.IsNull() {
		{{- if eq .Type "string" }}
		result.{{ .GoName }} = {{ .LocalName }}.StringVal
		{{- else if eq .Type "int64" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Int64
		{{- else if eq .Type "bool" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Bool
		{{- else if eq .Type "time.Time" -}}
		result.{{ .GoName }} = {{ .LocalName }}.Time.In(time.Local)
		{{- else if .IsEnum -}}
		result.{{ .GoName }} = {{ .Type }}({{ .LocalName }}.Int64)
		{{- end }}
	}
	{{- end -}}
	{{- end }}
	return &result, nil
}
{{- end }}

func (r *{{ $camelName }}Repository) diffEntity(source, target *{{ $pkgName }}.{{ .GoName }}) map[string]any {
	result := make(map[string]any, 0)

	{{ range (index .Types 0).Columns -}}
	{{- if .PK }}{{ continue }}
	{{ else if eq .Type "time.Time" -}}
	if !source.{{ .GoName }}.Equal(target.{{ .GoName }}) {
		result[transaction.{{ $goName }}ColumnName_{{ .GoName }}] = target.{{ .GoName }}
	}
	{{ else if eq .Type "byte" -}}
	if !bytes.Equal(source.{{ .GoName }}, target.{{ .GoName }}) {
		result[transaction.{{ $goName }}ColumnName_{{ .GoName }}] = target.{{ .GoName }}
	}
	{{ else if .IsList -}}
	if !slices.Equal(source.{{ .GoName }}, target.{{ .GoName }}) {
		result[transaction.{{ $goName }}ColumnName_{{ .GoName }}] = target.{{ .GoName }}
	}
	{{ else -}}
	if source.{{ .GoName }} != target.{{ .GoName }} {
		result[transaction.{{ $goName }}ColumnName_{{ .GoName }}] = target.{{ .GoName }}
	}
	{{ end -}}
	{{ end }}

	return result
}
