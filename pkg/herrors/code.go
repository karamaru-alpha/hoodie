package herrors

type errorCode int

const (
	errorCodeUnknown errorCode = iota
	errorCodeInvalidArgument
	errorCodeInternal
	// TODO: add another error codes
)

func (c errorCode) String() string {
	switch c {
	case errorCodeUnknown:
		return "Unknown"
	case errorCodeInvalidArgument:
		return "InvalidArgument"
	case errorCodeInternal:
		return "Internal"
	}
	return "Unknown"
}
