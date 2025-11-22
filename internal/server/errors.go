package server

import "errors"

var (
	ErrInvalidUUID        = errors.New("invalid uuid")
	ErrInvalidRequestBody = errors.New("invalid request body")
)
