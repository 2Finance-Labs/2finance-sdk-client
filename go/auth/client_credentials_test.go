package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientCredentialsTokenSourceFetchesAndCachesToken(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.Method; got != http.MethodPost {
			t.Fatalf("method = %s, want POST", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("Content-Type = %q, want form", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		if got := r.Form.Get("grant_type"); got != "client_credentials" {
			t.Fatalf("grant_type = %q", got)
		}
		if got := r.Form.Get("client_id"); got != "2finance-sdk-client" {
			t.Fatalf("client_id = %q", got)
		}
		if got := r.Form.Get("client_secret"); got != "secret" {
			t.Fatalf("client_secret = %q", got)
		}
		if got := r.Form.Get("scope"); got != "network:execute mcp:invoke" {
			t.Fatalf("scope = %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "token-123",
			"expires_in":   300,
			"token_type":   "Bearer",
		})
	}))
	defer server.Close()

	now := time.Unix(1000, 0)
	source, err := NewClientCredentialsTokenSource(ClientCredentialsConfig{
		TokenURL:     server.URL,
		ClientID:     "2finance-sdk-client",
		ClientSecret: "secret",
		Scopes:       []string{"network:execute", "mcp:invoke"},
		Now: func() time.Time {
			return now
		},
	})
	if err != nil {
		t.Fatalf("NewClientCredentialsTokenSource: %v", err)
	}

	token, err := source.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if token != "token-123" {
		t.Fatalf("token = %q", token)
	}
	token, err = source.Token(context.Background())
	if err != nil {
		t.Fatalf("Token cached: %v", err)
	}
	if token != "token-123" {
		t.Fatalf("cached token = %q", token)
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want one cached fetch", requests)
	}
}

func TestClientCredentialsTokenSourceRejectsMissingConfig(t *testing.T) {
	if _, err := NewClientCredentialsTokenSource(ClientCredentialsConfig{}); err == nil {
		t.Fatal("expected missing config error")
	}
}

func TestClientCredentialsTokenSourceRejectsTokenEndpointError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_client",
			"error_description": `client_secret="super-secret" access_token=token-123`,
		})
	}))
	defer server.Close()

	source, err := NewClientCredentialsTokenSource(ClientCredentialsConfig{
		TokenURL:     server.URL,
		ClientID:     "2finance-sdk-client",
		ClientSecret: "bad-secret",
	})
	if err != nil {
		t.Fatalf("NewClientCredentialsTokenSource: %v", err)
	}

	_, err = source.Token(context.Background())
	if err == nil {
		t.Fatal("expected token endpoint error")
	}
	if strings.Contains(err.Error(), "super-secret") || strings.Contains(err.Error(), "token-123") {
		t.Fatalf("error leaked sensitive value: %v", err)
	}
}

func TestClientCredentialsTokenSourceRejectsInsecureProductionEndpoint(t *testing.T) {
	t.Setenv("TWO_FINANCE_ENV", "production")
	if _, err := NewClientCredentialsTokenSource(ClientCredentialsConfig{
		TokenURL:     "http://authenticator.2finance.io/realms/2Finance/protocol/openid-connect/token",
		ClientID:     "2finance-sdk-client",
		ClientSecret: "secret",
	}); err == nil {
		t.Fatal("expected insecure endpoint error")
	}
}

func TestClientCredentialsTokenSourceAllowsLocalHTTPInProduction(t *testing.T) {
	t.Setenv("TWO_FINANCE_ENV", "production")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "token-123",
			"expires_in":   300,
		})
	}))
	defer server.Close()

	source, err := NewClientCredentialsTokenSource(ClientCredentialsConfig{
		TokenURL:     server.URL,
		ClientID:     "2finance-sdk-client",
		ClientSecret: "secret",
	})
	if err != nil {
		t.Fatalf("NewClientCredentialsTokenSource: %v", err)
	}
	token, err := source.Token(context.Background())
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if token != "token-123" {
		t.Fatalf("token = %q", token)
	}
}
