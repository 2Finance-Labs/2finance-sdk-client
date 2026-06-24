package auth

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

type captureTransport struct {
	header http.Header
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.header = req.Header.Clone()
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

func TestAuthTransportAddsBearerToken(t *testing.T) {
	base := &captureTransport{}
	transport := AuthTransport{
		Source: StaticTokenSource("access-token"),
		Base:   base,
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	defer resp.Body.Close()

	if got := base.header.Get("Authorization"); got != "Bearer access-token" {
		t.Fatalf("Authorization = %q, want bearer token", got)
	}
	if got := req.Header.Get("Authorization"); got != "" {
		t.Fatalf("original request Authorization = %q, want unchanged", got)
	}
}

func TestAuthTransportDoesNotDuplicateBearerPrefix(t *testing.T) {
	base := &captureTransport{}
	transport := AuthTransport{
		Source: StaticTokenSource("Bearer access-token"),
		Base:   base,
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	defer resp.Body.Close()

	if got := base.header.Get("Authorization"); got != "Bearer access-token" {
		t.Fatalf("Authorization = %q, want normalized bearer token", got)
	}
}

func TestRedactAuthorizationMasksAuthorizationHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer secret")
	headers.Set("Content-Type", "application/json")

	redacted := RedactAuthorization(headers)

	if got := redacted.Get("Authorization"); got != redactedAuthorization {
		t.Fatalf("redacted Authorization = %q, want %q", got, redactedAuthorization)
	}
	if got := redacted.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := headers.Get("Authorization"); got != "Bearer secret" {
		t.Fatalf("original Authorization = %q, want unchanged", got)
	}
}

func TestRedactSensitiveMasksOAuthValues(t *testing.T) {
	redacted := RedactSensitive(`Bearer abc.def client_secret="secret" refresh_token=refresh code=abc`)
	for _, forbidden := range []string{"abc.def", `"secret"`, "=refresh", "code=abc"} {
		if strings.Contains(redacted, forbidden) {
			t.Fatalf("redacted value leaked %q: %s", forbidden, redacted)
		}
	}
}
