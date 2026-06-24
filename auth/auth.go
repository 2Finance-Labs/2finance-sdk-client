package auth

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

const redactedAuthorization = "[REDACTED]"

var sensitiveValuePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9._~+/=-]+`),
	regexp.MustCompile(`(?i)(access_token|refresh_token|id_token|client_secret|password|code)(["']?\s*[:=]\s*["']?)[^"',\s}&]+`),
}

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
	clone.Header.Set("Authorization", BearerAuthorization(token))
	return base.RoundTrip(clone)
}

// BearerAuthorization returns a normalized Authorization header value.
func BearerAuthorization(accessToken string) string {
	trimmed := strings.TrimSpace(accessToken)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(trimmed), "bearer ") {
		return trimmed
	}
	return "Bearer " + trimmed
}

// RedactAuthorization returns a copy of headers with Authorization values masked.
func RedactAuthorization(headers http.Header) http.Header {
	redacted := headers.Clone()
	if _, ok := redacted["Authorization"]; ok {
		redacted["Authorization"] = []string{redactedAuthorization}
	}
	return redacted
}

// RedactSensitive masks bearer tokens and common OAuth credential fields in logs/errors.
func RedactSensitive(value string) string {
	redacted := value
	for _, pattern := range sensitiveValuePatterns {
		redacted = pattern.ReplaceAllStringFunc(redacted, func(match string) string {
			if strings.HasPrefix(strings.ToLower(match), "bearer ") {
				return "Bearer <redacted>"
			}
			separator := regexp.MustCompile(`["']?\s*[:=]\s*["']?`).FindString(match)
			if separator == "" {
				return "<redacted>"
			}
			key := strings.Split(match, separator)[0]
			return key + separator + "<redacted>"
		})
	}
	return redacted
}
