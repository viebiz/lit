package redis

// Error represents an error in redis package
// suppose to determine expected or unexpected error in another service
type Error string

func newError(msg string) Error {
	return Error(msg)
}

func (e Error) Error() string {
	return string(e)
}

var (
	ErrFailToSetValue = newError("fail to set value")

	ErrUnsupportedInputType = newError("unsupported input type")
)
