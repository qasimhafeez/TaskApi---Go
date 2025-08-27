package errs

import "errors"

var (
	ErrNotFound = errors.New("Not Found")
	ErrConflict = errors.New("Conflict")
	ErrBadRequest = errors.New("Bad Request")
)