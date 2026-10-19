package domain

import "errors"

var (
	// ErrNotFound reusable
	ErrNotFound = errors.New("resource not found")
	// ErrAlreadyExists reusable
	ErrAlreadyExists = errors.New("Resource already exists")
	// ErrUnauthorized reusable
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden reusable
	ErrForbidden = errors.New("Forbidden")
	// ErrInvalidInput reusable
	ErrInvalidInput = errors.New("Invalid input data")
	// ErrInternal reusable
	ErrInternal = errors.New("internal error")
)
