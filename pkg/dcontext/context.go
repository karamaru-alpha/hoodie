package dcontext

import "context"

type key[K any] struct{}

func withValue[K, V any](ctx context.Context, val V) context.Context {
	return context.WithValue(ctx, key[K]{}, val)
}

func value[T any](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(key[T]{}).(T)
	return val, ok
}
