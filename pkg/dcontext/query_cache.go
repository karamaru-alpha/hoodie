package dcontext

import (
	"context"
	"sync"
)

type QueryCache interface {
	Reset()
	GetCacher(key any) (any, bool)
	SetCacher(key any, cacher any)
	Nil() bool
}

type queryCache struct {
	sync.RWMutex
	cacheMap map[any]any
}

func NewQueryCache() QueryCache {
	q := &queryCache{}
	q.init()
	return q
}

func SetQueryCacheInContext(ctx context.Context, qctx QueryCache) context.Context {
	return withValue[*queryCache](ctx, qctx)
}

func ExtractQueryCache(ctx context.Context) QueryCache {
	qctx, _ := value[*queryCache](ctx)
	return qctx
}

func (q *queryCache) init() {
	q.Lock()
	defer q.Unlock()
	q.cacheMap = make(map[any]any)
}

func (q *queryCache) Reset() {
	if q == nil {
		return
	}
	q.init()
}

func (q *queryCache) GetCacher(key any) (any, bool) {
	if q == nil {
		return nil, false
	}
	q.RLock()
	defer q.RUnlock()
	cr, ok := q.cacheMap[key]
	return cr, ok
}

func (q *queryCache) SetCacher(key, cacher any) {
	if q == nil {
		return
	}
	q.Lock()
	defer q.Unlock()
	q.cacheMap[key] = cacher
}

func (q *queryCache) Nil() bool {
	return q == nil
}
