package service

import "errors"

// Common service-layer errors.
var (
	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
	ErrBadRequest = errors.New("bad request")
	ErrConflict   = errors.New("conflict")
)
