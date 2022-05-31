package response

import (
	"fmt"

	"github.com/pkg/errors"
)

type ResponseError struct {
	toast string
	err   error
}

// Toast
func (r *ResponseError) Toast() string { return r.toast }

// SetToast
func (r *ResponseError) SetToast(toast string) { r.toast = toast }

// Error
func (r *ResponseError) Error() string { return r.err.Error() }

// SetError
func (r *ResponseError) SetError(err error) { r.err = err }

// Unwrap
func (r *ResponseError) Unwrap() error { return r.err }

// Cause
func (r *ResponseError) Cause() error { return r.err }

// WrapToast return a new ResponseError
func WrapToast(err error, toast string) *ResponseError {
	if err == nil {
		return &ResponseError{
			err:   errors.New(toast),
			toast: toast,
		}
	}

	return &ResponseError{
		err:   errors.Wrap(err, toast),
		toast: toast,
	}
}

// WrapToastf return a new format ResponseError
func WrapToastf(err error, toast string, args ...interface{}) *ResponseError {
	if err == nil {
		return &ResponseError{
			err:   errors.Errorf(toast, args...),
			toast: fmt.Sprintf(toast, args...),
		}
	}

	return &ResponseError{
		err:   errors.Wrapf(err, toast, args...),
		toast: fmt.Sprintf(toast, args...),
	}
}
