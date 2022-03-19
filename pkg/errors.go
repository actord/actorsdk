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
