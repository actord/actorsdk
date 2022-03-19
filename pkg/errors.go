package actorsdk

import "errors"

var ErrorNotFound = errors.New("not found")

func newErrorFromSDKResponse(error string) error {
	switch error {
	case "not found":
		return ErrorNotFound
	default:
		return errors.New(error)
	}
}

type StatusError struct {
	error      string
	statusCode int
}

func (err StatusError) Error() string {
	return err.error
}

func (err StatusError) StatusCode() int {
	return err.statusCode
}

func NewStatusError(statusCode int, error string) StatusError {
	return StatusError{
		statusCode: statusCode,
		error:      error,
	}
}
