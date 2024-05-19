{{ template "autogen_comment" }}
{{ $goName := .GoName -}}
package base

import (
	"context"
	"sync"
	"time"

	"github.com/scylladb/go-set/f64set"
	"github.com/scylladb/go-set/i32set"
	"github.com/scylladb/go-set/i64set"
	"github.com/scylladb/go-set/strset"

	"github.com/karamaru-alpha/days/pkg/domain/entity/{{ .PkgName }}"
	"github.com/karamaru-alpha/days/pkg/domain/enum"
	"github.com/karamaru-alpha/days/pkg/dcontext"
)

type {{ .GoName }}SearchResultCache interface {
	Lock()
	Unlock()
	GetByPK(pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, SearchResultType)
	GetByPKs(pks {{ .PkgName }}.{{ .GoName }}PKs) ({{ .PkgName }}.{{ .GoName }}Slice, SearchResultType)
	GetOriginByPK(pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, SearchResultType)
	GetByQueryConditions(queryConditions []*{{ .GoName }}QueryCondition) ({{ .PkgName }}.{{ .GoName }}Slice, SearchResultType)
	SetResult(row *{{ .PkgName }}.{{ .GoName }})
	SetResults(rows {{ .PkgName }}.{{ .GoName }}Slice)
	SetAsNotFound(pk *{{ .PkgName }}.{{ .GoName }}PK)
	AppendQueryConditions(queryConditions []*{{ .GoName }}QueryCondition)
	ReplaceWithCache(rows {{ .PkgName }}.{{ .GoName }}Slice) {{ .PkgName }}.{{ .GoName }}Slice
	FilterByQueryConditions(queryConditions []*{{ .GoName }}QueryCondition) {{ .PkgName }}.{{ .GoName }}Slice
}

type {{ .CamelName }}SearchResultRow struct {
	searchResultType SearchResultType
	searchResult     *{{ .PkgName }}.{{ .GoName }}
}

type {{ .CamelName }}SearchResultCacheKey struct{}

type {{ .CamelName }}SearchResultCache struct {
	mutex sync.RWMutex
	origin map[string]*{{ .CamelName }}SearchResultRow
	updated map[string]*{{ .CamelName }}SearchResultRow
	queryConditionsList [][]*{{ .GoName }}QueryCondition
}

func Extract{{ .GoName }}SearchResultCache(ctx context.Context) {{ .GoName }}SearchResultCache {
	qctx := dcontext.ExtractQueryCache(ctx)
	if qctx.Nil() {
		return &{{ .CamelName }}SearchResultCache{}
	}

	cacher, ok := qctx.GetCacher({{ .CamelName }}SearchResultCacheKey{})
	if !ok {
		cacher = &{{ .CamelName }}SearchResultCache{
			origin:              make(map[string]*{{ .CamelName }}SearchResultRow),
			updated:             make(map[string]*{{ .CamelName }}SearchResultRow),
			queryConditionsList: make([][]*{{ .GoName }}QueryCondition, 0),
		}
		qctx.SetCacher({{ .CamelName }}SearchResultCacheKey{}, cacher)
	}
	return cacher.(*{{ .CamelName }}SearchResultCache)
}

func (src *{{ .CamelName }}SearchResultCache) Lock() {
	src.mutex.Lock()
}

func (src *{{ .CamelName }}SearchResultCache) Unlock() {
	src.mutex.Unlock()
}

func (src *{{ .CamelName }}SearchResultCache) GetByPK(pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, SearchResultType) {
	if row, ok := src.updated[pk.Key()]; ok {
		return row.searchResult, row.searchResultType
	}

	resultType := SearchResultTypeNotSearched
	hasCache := src.hasCache([]*{{ .GoName }}QueryCondition{
		{{- range .PKColumns }}
        {
            column: {{ $goName }}ColumnName{{ .GoName }},
            operator: ConditionOperatorEq,
            value: pk.{{ .GoName }},
        },
		{{- end }}
	})

	if hasCache {
		resultType = SearchResultTypeNotFound
	}
	return nil, resultType
}

func (src *{{ .CamelName }}SearchResultCache) GetOriginByPK(pk *{{ .PkgName }}.{{ .GoName }}PK) (*{{ .PkgName }}.{{ .GoName }}, SearchResultType) {
	if row, ok := src.origin[pk.Key()]; ok {
		return row.searchResult, row.searchResultType
	}

	resultType := SearchResultTypeNotSearched
	hasCache := src.hasCache([]*{{ .GoName }}QueryCondition{
		{{- range .PKColumns }}
        {
            column: "{{ .GoName }}",
            operator: ConditionOperatorEq,
            value: pk.{{ .GoName }},
        },
		{{- end }}
	})

	if hasCache {
		resultType = SearchResultTypeNotFound
	}
	return nil, resultType
}

func (src *{{ .CamelName }}SearchResultCache) GetByPKs(pks {{ .PkgName }}.{{ .GoName }}PKs) ({{ .PkgName }}.{{ .GoName }}Slice, SearchResultType) {
	result := make({{ .PkgName }}.{{ .GoName }}Slice, 0, len(pks))

	for _, pk := range pks {
		value, resultType := src.GetByPK(pk)
		if resultType == SearchResultTypeNotSearched {
			return nil, SearchResultTypeNotSearched
		}
		if value != nil {
			result = append(result, value)
		}
	}

	if len(result) == 0 {
		return nil, SearchResultTypeNotFound
	}
	return result, SearchResultTypeFound
}

func (src *{{ .CamelName }}SearchResultCache) GetByQueryConditions(queryConditions []*{{ .GoName }}QueryCondition) ({{ .PkgName }}.{{ .GoName }}Slice, SearchResultType) {
	if !src.hasCache(queryConditions) {
		return {{ .PkgName }}.{{ .GoName }}Slice{}, SearchResultTypeNotSearched
	}

	result := src.FilterByQueryConditions(queryConditions)

	if len(result) == 0 {
		return {{ .PkgName }}.{{ .GoName }}Slice{}, SearchResultTypeNotFound
	}
	return result, SearchResultTypeFound
}

func (src *{{ .CamelName }}SearchResultCache) SetResult(row *{{ .PkgName }}.{{ .GoName }}) {
	if src.origin == nil {
		return
	}

	key := row.GetPK().Key()

	if _, ok := src.origin[key]; !ok {
		src.origin[key] = &{{ .CamelName }}SearchResultRow{
			searchResultType: SearchResultTypeFound,
			searchResult:     row.FullDeepCopy(),
		}
	}
	src.updated[key] = &{{ .CamelName }}SearchResultRow{
		searchResultType: SearchResultTypeFound,
		searchResult:     row,
	}
}

func (src *{{ .CamelName }}SearchResultCache) SetResults(rows {{ .PkgName }}.{{ .GoName }}Slice) {
	for _, row := range rows {
		src.SetResult(row)
	}
}

func (src *{{ .CamelName }}SearchResultCache) SetAsNotFound(pk *{{ .PkgName }}.{{ .GoName }}PK) {
	if src.origin == nil {
		return
	}

	key := pk.Key()

	if _, ok := src.origin[key]; !ok {
		src.origin[key] = &{{ .CamelName }}SearchResultRow{
			searchResultType: SearchResultTypeNotFound,
			searchResult:     nil,
		}
	}
	src.updated[key] = &{{ .CamelName }}SearchResultRow{
		searchResultType: SearchResultTypeNotFound,
		searchResult:     nil,
	}
}

func (src *{{ .CamelName }}SearchResultCache) AppendQueryConditions(queryConditions []*{{ .GoName }}QueryCondition) {
	if src.queryConditionsList == nil {
		return
	}

	src.queryConditionsList = append(src.queryConditionsList, queryConditions)
}

func (src *{{ .CamelName }}SearchResultCache) ReplaceWithCache(rows {{ .PkgName }}.{{ .GoName }}Slice) {{ .PkgName }}.{{ .GoName }}Slice {
	result := make({{ .PkgName }}.{{ .GoName }}Slice, 0, len(rows))

	for _, row := range rows {
		if cache, resultType := src.GetByPK(row.GetPK()); cache != nil && resultType == SearchResultTypeFound {
			result = append(result, cache)
		} else if resultType != SearchResultTypeFound {
			result = append(result, row)
		}
	}

	return result
}

func (src *{{ .CamelName }}SearchResultCache) FilterByQueryConditions(queryConditions []*{{ .GoName }}QueryCondition) {{ .PkgName }}.{{ .GoName }}Slice {
	result := make({{ .PkgName }}.{{ .GoName }}Slice, 0)

	for _, row := range src.updated {
		isSatisfy := true
		for _, cond := range queryConditions {
			if !row.isSatisfy(cond) {
				isSatisfy = false
				break
			}
		}
		if isSatisfy && row.searchResult != nil {
			result = append(result, row.searchResult)
		}
	}

	return result
}

func (src *{{ .CamelName }}SearchResultCache) hasCache(queryConditions []*{{ .GoName }}QueryCondition) bool {
	for _, sourceQueryConditions := range src.queryConditionsList {
		if contains{{ .GoName }}Conditions(sourceQueryConditions, queryConditions) {
			return true
		}
	}

	return false
}

func (row *{{ .CamelName }}SearchResultRow) isSatisfy(queryCondition *{{ .GoName }}QueryCondition) bool {
	searchResult := row.searchResult
	if searchResult == nil {
		return false
	}

	switch queryCondition.column {
	{{ range .PKColumns -}}
	case {{ $goName }}ColumnName{{ .GoName }}:
		switch queryCondition.operator {
		case ConditionOperatorEq:
            conditionValue := queryCondition.value.({{ .Type }})
			return conditionValue == searchResult.{{ .GoName }}
		case ConditionOperatorIn:
            conditionValueSet := queryCondition.value.(*{{ .SetType }}.Set)
			{{ if eq .Type "time.Time" -}}
				return conditionValueSet.Has(searchResult.{{ .GoName }}.UnixNano())
			{{ else if .IsEnum -}}
				return conditionValueSet.Has(searchResult.{{ .GoName }}.Int32())
			{{ else -}}
				return conditionValueSet.Has(searchResult.{{ .GoName }})
			{{ end -}}
		}
	{{ end -}}
	}

	return false
}

func contains{{ .GoName }}Conditions(sourceQueryConditions, targetQueryConditions []*{{ .GoName }}QueryCondition) bool {
	if len(sourceQueryConditions) > len(targetQueryConditions) {
		return false
	}

	for _, sourceQueryCondition := range sourceQueryConditions {
		var contains bool
		for _, targetQueryCondition := range targetQueryConditions {
			if contains{{ .GoName }}Condition(sourceQueryCondition, targetQueryCondition) {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}

	return true
}

func contains{{ .GoName }}Condition(sourceQueryCondition, targetQueryCondition *{{ .GoName }}QueryCondition) bool {
	if sourceQueryCondition.column != targetQueryCondition.column {
		return false
	}

	switch sourceQueryCondition.column {
{{- range .PKColumns }}
	case {{ $goName }}ColumnName{{ .GoName }}:
		switch sourceQueryCondition.operator {
		case ConditionOperatorEq:
			switch targetQueryCondition.operator {
			case ConditionOperatorEq:
                sourceConditionValue := sourceQueryCondition.value.({{ .Type }})
                targetConditionValue := targetQueryCondition.value.({{ .Type }})
				return sourceConditionValue == targetConditionValue
			case ConditionOperatorIn:
				{{- if eq .Type "time.Time" }}
                sourceConditionValue := sourceQueryCondition.value.({{ .Type }}).UnixNano()
				{{- else if .IsEnum }}
                sourceConditionValue := sourceQueryCondition.value.({{ .Type }}).Int32()
				{{- else }}
                sourceConditionValue := sourceQueryCondition.value.({{ .Type }})
				{{- end }}
                targetQueryConditionValueSet := targetQueryCondition.value.(*{{ .SetType }}.Set)
				return targetQueryConditionValueSet.Size() == 1 && targetQueryConditionValueSet.Has(sourceConditionValue)
			}
		case ConditionOperatorIn:
			switch targetQueryCondition.operator {
			case ConditionOperatorEq:
                sourceConditionValues := sourceQueryCondition.value.(*{{ .SetType }}.Set)
				{{- if eq .Type "time.Time" }}
                targetConditionValue := targetQueryCondition.value.({{ .Type }}).UnixNano()
				{{- else if .IsEnum }}
                targetConditionValue := targetQueryCondition.value.({{ .Type }}).Int32()
				{{- else }}
                targetConditionValue := targetQueryCondition.value.({{ .Type }})
				{{- end }}
				return sourceConditionValues.Has(targetConditionValue)
			case ConditionOperatorIn:
                sourceConditionValues := sourceQueryCondition.value.(*{{ .SetType }}.Set)
                targetQueryConditionValueSet := targetQueryCondition.value.(*{{ .SetType }}.Set)
				return targetQueryConditionValueSet.IsSuperset(sourceConditionValues)
			}
		}
	{{- end }}
	}

	return false
}
