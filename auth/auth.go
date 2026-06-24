package auth

import (
	"context"
	"net/http"
	"strings"
)

const redactedAuthorization = "[REDACTED]"

// TokenSource provides bearer access tokens for HTTP clients.
type TokenSource interface {
	Token(ctx context.Context) (string, error)
}

// StaticTokenSource is useful when the caller already owns token acquisition.
type StaticTokenSource string

// Token returns the configured static token.
func (s StaticTokenSource) Token(context.Context) (string, error) {
	return string(s), nil
}

// AuthTransport injects bearer authentication into outbound HTTP requests.
type AuthTransport struct {
	Source TokenSource
	Base   http.RoundTripper
}

// RoundTrip adds Authorization: Bearer <token> when Source returns a non-empty token.
func (t AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	if t.Source == nil {
		return base.RoundTrip(req)
	}

	token, err := t.Source.Token(req.Context())
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(token) == "" {
		return base.RoundTrip(req)
	}

	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+token)
	return base.RoundTrip(clone)
}

// RedactAuthorization returns a copy of headers with Authorization values masked.
func RedactAuthorization(headers http.Header) http.Header {
	redacted := headers.Clone()
	if _, ok := redacted["Authorization"]; ok {
		redacted["Authorization"] = []string{redactedAuthorization}
	}
	return redacted
}
