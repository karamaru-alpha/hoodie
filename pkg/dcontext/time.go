package dcontext

import (
	"context"
	"time"
)

type Time interface {
	Now() time.Time
}

type timeImpl struct {
	now time.Time
}

func SetTimeContext(ctx context.Context) context.Context {
	return SetTimeContextWithNow(ctx, time.Now())
}

func SetTimeContextWithNow(ctx context.Context, now time.Time) context.Context {
	t := &timeImpl{now: now}
	return withValue[*timeImpl](ctx, t)
}

func ExtractTime(ctx context.Context) Time {
	rctx, _ := value[*timeImpl](ctx)
	return rctx
}

func Now(ctx context.Context) time.Time {
	return ExtractTime(ctx).Now()
}

func (c *timeImpl) Now() time.Time {
	if c == nil {
		return time.Now()
	}
	return c.now
}
