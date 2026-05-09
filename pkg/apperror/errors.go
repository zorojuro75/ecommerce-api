package apperror

import "errors"

var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrBadRequest   = errors.New("invalid request")
    ErrConflict     = errors.New("resource already exists")
    ErrOutOfStock   = errors.New("out of stock")
    ErrInvalidInput = errors.New("invalid input")
)