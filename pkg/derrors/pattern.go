package derrors

import (
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
)

type ErrorPattern struct {
	ErrorCode         errorCode
	HTTPStatusCode    int
	ConnectStatusCode connect.Code
	GRPCStatusCode    codes.Code
}

var (
	Unknown = ErrorPattern{
		ErrorCode:         errorCodeUnknown,
		HTTPStatusCode:    http.StatusInternalServerError,
		ConnectStatusCode: connect.CodeUnknown,
		GRPCStatusCode:    codes.Unknown,
	}
	InvalidArgument = ErrorPattern{
		ErrorCode:         errorCodeInvalidArgument,
		HTTPStatusCode:    http.StatusBadRequest,
		ConnectStatusCode: connect.CodeInvalidArgument,
		GRPCStatusCode:    codes.InvalidArgument,
	}
	Internal = ErrorPattern{
		ErrorCode:         errorCodeInternal,
		HTTPStatusCode:    http.StatusInternalServerError,
		ConnectStatusCode: connect.CodeInternal,
		GRPCStatusCode:    codes.Internal,
	}
)
