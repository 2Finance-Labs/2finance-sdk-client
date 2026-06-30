package twofinance

import (
	"time"

	"github.com/2Finance-Labs/2finance-sdk-client/internal/service"
)

// RequestOption configures per-call HTTP options for SDK service requests.
type RequestOption = service.RequestOption

// WithHeader adds a per-call HTTP header.
func WithHeader(key, value string) RequestOption {
	return service.WithHeader(key, value)
}

// WithIdempotencyKey adds the standard Idempotency-Key header.
func WithIdempotencyKey(key string) RequestOption {
	return service.WithIdempotencyKey(key)
}

// WithQueryParam adds a per-call query parameter.
func WithQueryParam(key, value string) RequestOption {
	return service.WithQueryParam(key, value)
}

// WithPagination adds standard page and limit query parameters.
func WithPagination(page, limit int) RequestOption {
	return service.WithPagination(page, limit)
}

// WithTimeout bounds a single HTTP request.
func WithTimeout(timeout time.Duration) RequestOption {
	return service.WithTimeout(timeout)
}

// WithMaxRetries retries retryable HTTP responses up to maxRetries times.
func WithMaxRetries(maxRetries int) RequestOption {
	return service.WithMaxRetries(maxRetries)
}
