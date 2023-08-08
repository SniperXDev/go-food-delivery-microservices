package customErrors

import (
	"fmt"
	"io"

	"emperror.dev/errors"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/contracts"
)

// https://klotzandrew.com/blog/error-handling-in-golang
// https://banzaicloud.com/blog/error-handling-go/
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
// https://github.com/go-stack/stack
// https://github.com/juju/errors
// https://github.com/emperror/errors
// https://github.com/pkg/errors/issues/75
type customError struct {
	statusCode int
	message    string
	err        error
}

type CustomError interface {
	error
	contracts.Wrapper
	contracts.Causer
	contracts.Formatter
	IsCustomError() bool
	Status() int
	Message() string
}

func NewCustomError(err error, code int, message string) CustomError {
	m := &customError{
		statusCode: code,
		err:        err,
		message:    message,
	}

	return m
}

func (e *customError) IsCustomError() bool {
	return true
}

func (e *customError) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}

	return e.message
}

func (e *customError) Message() string {
	return e.message
}

func (e *customError) Status() int {
	return e.statusCode
}

func (e *customError) Cause() error {
	return e.err
}

func (e *customError) Unwrap() error {
	return e.err
}

func (e *customError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.Cause())
			io.WriteString(s, e.message)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}

func GetCustomError(err error) CustomError {
	if IsCustomError(err) {
		var internalErr CustomError
		errors.As(err, &internalErr)

		return internalErr
	}

	return nil
}

func IsCustomError(err error) bool {
	var customErr CustomError

	_, ok := err.(CustomError)
	if ok {
		return true
	}

	if errors.As(err, &customErr) {
		return customErr.IsCustomError()
	}

	return false
}
