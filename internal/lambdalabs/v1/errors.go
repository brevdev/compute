package v1

import (
	"errors"
	"fmt"
)

var (
	// Lambda Labs specific errors
	ErrLambdaLabsAPIError                = errors.New("lambda labs api error")
	ErrLambdaLabsAuthError               = errors.New("lambda labs authentication error")
	ErrLambdaLabsQuotaExceeded           = errors.New("lambda labs quota exceeded")
	ErrLambdaLabsInstanceNotFound        = errors.New("lambda labs instance not found")
	ErrLambdaLabsInstanceTypeUnavailable = errors.New("lambda labs instance type unavailable")
)

// LambdaLabsError wraps Lambda Labs specific errors with additional context
type LambdaLabsError struct {
	Err     error
	Message string
	Code    int
}

func (e *LambdaLabsError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("lambda labs error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("lambda labs error: %v", e.Err)
}

func (e *LambdaLabsError) Unwrap() error {
	return e.Err
}
