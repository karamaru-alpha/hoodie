package dcontext

import (
	"context"
)

type User interface {
	SetUser(p *SetUserParam)
	GetUserID() string
}

type user struct {
	userID string
}

func NewUser() User {
	return &user{}
}

func SetUserInContext(ctx context.Context, u User) context.Context {
	return withValue[*user](ctx, u)
}

func ExtractUser(ctx context.Context) User {
	u, _ := value[*user](ctx)
	return u
}

type SetUserParam struct {
	UserID string
}

func (c *user) SetUser(p *SetUserParam) {
	if c == nil {
		return
	}
	if p == nil {
		return
	}
	c.userID = p.UserID
}

func (c *user) GetUserID() string {
	if c == nil {
		return ""
	}
	return c.userID
}
