package common

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden: you do not have permission for this resource")
	ErrConflict       = errors.New("resource conflict or already exists")
	ErrValidation     = errors.New("validation failed")
	ErrInternalServer = errors.New("internal server error")
)
