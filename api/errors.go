package api

import "fmt"

// NewError creates an Error with the given code and message.
func NewError(code ErrorCode, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

// NewErrorf creates an Error with a formatted message.
func NewErrorf(code ErrorCode, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// IsNotFound returns true if the error has code ErrNotFound.
func IsNotFound(err *Error) bool {
	return err != nil && err.Code == ErrNotFound
}

// IsUnauthorized returns true if the error has code ErrUnauthorized.
func IsUnauthorized(err *Error) bool {
	return err != nil && err.Code == ErrUnauthorized
}

// IsVersionMismatch returns true if the error has code ErrVersionMismatch.
func IsVersionMismatch(err *Error) bool {
	return err != nil && err.Code == ErrVersionMismatch
}
