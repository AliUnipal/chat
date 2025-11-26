package server

import "errors"

var (
	ErrInvalidUUID        = errors.New("invalid uuid")
	ErrInvalidBase64Value = errors.New("invalid base64 value")
	ErrInvalidRequestBody = errors.New("invalid request body")
)
