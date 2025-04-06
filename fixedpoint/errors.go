package fixedpoint

import "fmt"

type internalError struct {
	data any
	msg  string
}

func (e *internalError) Error() string {
	return fmt.Sprintf("internal error: %s: %v", e.msg, e.data)
}

func newInternalError(data any, msg string) error {
	return &internalError{
		data: data,
		msg:  msg,
	}
}
