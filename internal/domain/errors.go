package domain

import "errors"

var (
	// ErrNotFound reusable
	ErrNotFound = errors.New("resource not found")
	// ErrAlreadyExists reusable
	ErrAlreadyExists = errors.New("resource already exists")
	// ErrUnauthorized reusable
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden reusable
	ErrForbidden = errors.New("forbidden")
	// ErrInvalidInput reusable
	ErrInvalidInput = errors.New("invalid input data")
	// ErrInternal reusable
	ErrInternal = errors.New("internal error")
)
