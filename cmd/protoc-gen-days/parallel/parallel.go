package parallel

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/karamaru-alpha/days/cmd/protoc-gen-days/perrors"
)

const DefaultSize int = 100

type Group struct {
	eg *errgroup.Group
}

func NewGroupWithContext(ctx context.Context, size int) (*Group, context.Context) {
	eg, egctx := errgroup.WithContext(ctx)
	eg.SetLimit(size)
	return &Group{
		eg: eg,
	}, egctx
}

func (g *Group) Go(ctx context.Context, f func(ctx context.Context) error) {
	if g.eg == nil {
		return
	}

	g.eg.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = perrors.New("panic occurred").SetValues(map[string]any{"recovered": r})
			}
		}()
		return f(ctx)
	})
}

func (g *Group) Wait() error {
	if g.eg == nil {
		return nil
	}

	return g.eg.Wait()
}
