package di

import (
	"connectrpc.com/connect"
)

func newHandlerOption() ([]connect.HandlerOption, error) {
	interceptors := []connect.Interceptor{
		// TODO: implement me
	}
	return []connect.HandlerOption{
		connect.WithInterceptors(interceptors...),
	}, nil
}
