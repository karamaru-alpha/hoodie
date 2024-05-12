package herrors

type errorCode int

const (
	errorCodeUnknown errorCode = iota
	errorCodeInvalidArgument
	errorCodeInternal
	// TODO: add another error codes
)
