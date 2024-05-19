package dcontext

import (
	"context"
)

type Request interface {
	SetRequest(p *SetRequestParam)
	GetIP() string
}

type SetRequestParam struct {
	IP string
}

type request struct {
	// ip address
	ip string
}

func NewRequest() Request {
	return &request{}
}

func SetRequestInContext(ctx context.Context, rctx Request) context.Context {
	return withValue[*request](ctx, rctx)
}

func ExtractRequest(ctx context.Context) Request {
	rctx, _ := value[*request](ctx)
	return rctx
}

func (c *request) SetRequest(p *SetRequestParam) {
	if c == nil {
		return
	}
	if p == nil {
		return
	}
	c.ip = p.IP
}

func (c *request) GetIP() string {
	if c == nil {
		return ""
	}
	return c.ip
}
