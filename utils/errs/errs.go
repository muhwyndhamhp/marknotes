package errs

import (
	"fmt"
	"runtime"
	"strings"
)

// Wrap is wrapper that includes error trace in errors
func Wrap(any interface{}, a ...interface{}) error {
	if any != nil {
		err := error(nil)

		switch any := any.(type) {
		case string:
			err = fmt.Errorf(any, a...)
		case error:
			err = fmt.Errorf(any.Error(), a...)
		default:
			err = fmt.Errorf("%v", err)
		}

		_, fn, line, _ := runtime.Caller(1)

		return fmt.Errorf("%s:%d %v", fn, line, err)
	}

	return nil
}

// Errorf is wrapper to create new error along with error trace
func Errorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	_, fn, line, _ := runtime.Caller(1)

	return fmt.Errorf("%s:%d %v", fn, line, err)
}

// ParseError is a method to parse error code and error message from a single error trace.
// Please make sure that the error have following format:
// <<ErrorCode>>: <<ErrorMessage>>
// Example:
// PaymentUnauthorized: payment cannot be authorized
func ParseError(err error) (code, message string) {
	splits := strings.Split(err.Error(), ":")
	if len(splits) != 2 {
		return "GeneralError", err.Error()
	}

	return strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
}
