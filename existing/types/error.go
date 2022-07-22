package types

import (
	"fmt"
	"weavelab.xyz/monorail/shared/wlib/werror"
)

type Reason string

type Error struct {
	err     *werror.Error
	code    werror.Code
	reason  Reason
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Reason: %s Message: %s Error: %s", e.reason, e.message, e.err.Error())
}

func (e *Error) AddTag(key string, value interface{}) {
	e.err.Add(key, value)
}

type ErrorOption func(e *Error)

func NewError(msg string, options ...ErrorOption) *Error {
	e := &Error{
		err: werror.New(msg),
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func WrapError(err error, msg string, options ...ErrorOption) *Error {
	e := &Error{
		err: werror.Wrap(err, msg),
	}

	for _, option := range options {
		option(e)
	}

	return e
}

func WithFileSystemError() ErrorOption {
	return func(e *Error) {

	}
}

func WithTag(key string, value interface{}) ErrorOption {
	return func(e *Error) {
		e.err = e.err.Add(key, value)
	}
}

var (
	UninitializedError    = werror.CodedTemplate("sync app has not been initialized", werror.CodeUnavailable)
	NotImplementedError   = werror.CodedTemplate("functionality not implemented", werror.CodeUnavailable)
	FunctionalityDisabled = werror.CodedTemplate("functionality has been disabled", werror.CodeConflict)
)
