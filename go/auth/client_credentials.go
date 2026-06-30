package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// ClientCredentialsConfig configures OAuth2 client credentials token acquisition.
type ClientCredentialsConfig struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scopes       []string
	HTTPClient   *http.Client
	Now          func() time.Time
	ExpirySkew   time.Duration
}

// ClientCredentialsTokenSource fetches and caches OIDC access tokens for backend jobs.
type ClientCredentialsTokenSource struct {
	config ClientCredentialsConfig

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

// NewClientCredentialsTokenSource returns a TokenSource backed by the OAuth2
// client_credentials grant. Secrets must come from the caller's environment or
// secret manager, not source code.
func NewClientCredentialsTokenSource(config ClientCredentialsConfig) (*ClientCredentialsTokenSource, error) {
	if strings.TrimSpace(config.TokenURL) == "" {
		return nil, errors.New("auth: token URL is required")
	}
	if err := rejectInsecureProductionURL(config.TokenURL); err != nil {
		return nil, err
	}
	if strings.TrimSpace(config.ClientID) == "" {
		return nil, errors.New("auth: client ID is required")
	}
	if strings.TrimSpace(config.ClientSecret) == "" {
		return nil, errors.New("auth: client secret is required")
	}
	if config.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.ExpirySkew == 0 {
		config.ExpirySkew = 30 * time.Second
	}
	return &ClientCredentialsTokenSource{config: config}, nil
}

// Token returns a cached access token or fetches a fresh one when it is near expiry.
func (s *ClientCredentialsTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.config.Now()
	if s.accessToken != "" && now.Before(s.expiresAt.Add(-s.config.ExpirySkew)) {
		return s.accessToken, nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", s.config.ClientID)
	form.Set("client_secret", s.config.ClientSecret)
	if len(s.config.Scopes) > 0 {
		form.Set("scope", strings.Join(s.config.Scopes, " "))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("auth: build token request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.config.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth: fetch token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("auth: decode token response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if tokenResponse.Error != "" {
			detail := tokenResponse.Error
			if tokenResponse.Description != "" {
				detail += ": " + tokenResponse.Description
			}
			return "", fmt.Errorf("auth: token endpoint returned %d: %s", resp.StatusCode, RedactSensitive(detail))
		}
		return "", fmt.Errorf("auth: token endpoint returned %d", resp.StatusCode)
	}
	if strings.TrimSpace(tokenResponse.AccessToken) == "" {
		return "", errors.New("auth: token response missing access_token")
	}
	if tokenResponse.ExpiresIn <= 0 {
		tokenResponse.ExpiresIn = 300
	}

	s.accessToken = tokenResponse.AccessToken
	s.expiresAt = now.Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	return s.accessToken, nil
}

func rejectInsecureProductionURL(rawURL string) error {
	env := strings.ToLower(strings.TrimSpace(firstEnv("TWO_FINANCE_ENV", "APP_ENV", "ENV")))
	if env != "prod" && env != "production" && env != "prod_secrets" {
		return nil
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("auth: invalid token URL: %w", err)
	}
	if parsed.Scheme != "http" {
		return nil
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return nil
	}
	return errors.New("auth: production token URL must use HTTPS")
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return ""
}
