package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	v1 "github.com/brevdev/compute/pkg/v1"
)

func handleLLAPIError(ctx context.Context, resp *http.Response, err error) error {
	if err == nil {
		return nil
	}

	return handleLLErrToCloudErr(err)
}

func handleLLErrToCloudErr(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()
	
	if strings.Contains(errStr, "insufficient capacity") ||
		strings.Contains(errStr, "no capacity") ||
		strings.Contains(errStr, "capacity not available") {
		return v1.ErrInsufficientResources
	}
	
	if strings.Contains(errStr, "quota") ||
		strings.Contains(errStr, "limit exceeded") ||
		strings.Contains(errStr, "too many") {
		return v1.ErrOutOfQuota
	}
	
	if strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "does not exist") {
		return v1.ErrInstanceNotFound
	}
	
	if strings.Contains(errStr, "service unavailable") ||
		strings.Contains(errStr, "temporarily unavailable") {
		return v1.ErrServiceUnavailable
	}
	
	return fmt.Errorf("lambda labs error: %w", err)
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	
	if strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "forbidden") ||
		strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "bad request") {
		return false
	}
	
	return strings.Contains(errStr, "capacity") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary") ||
		strings.Contains(errStr, "service unavailable") ||
		strings.Contains(errStr, "internal server error")
}
