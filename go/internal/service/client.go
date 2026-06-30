package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client is the shared JSON-over-HTTP transport used by service SDKs.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RequestOptions carries per-call HTTP options shared by service SDKs.
type RequestOptions struct {
	Headers        http.Header
	IdempotencyKey string
	Query          url.Values
	Timeout        time.Duration
	MaxRetries     int
	Page           int
	Limit          int
}

// RequestOption mutates RequestOptions.
type RequestOption func(*RequestOptions)

// WithHeader adds a per-call HTTP header.
func WithHeader(key, value string) RequestOption {
	return func(options *RequestOptions) {
		if options.Headers == nil {
			options.Headers = make(http.Header)
		}
		options.Headers.Set(key, value)
	}
}

// WithIdempotencyKey adds the standard Idempotency-Key header.
func WithIdempotencyKey(key string) RequestOption {
	return func(options *RequestOptions) {
		options.IdempotencyKey = key
	}
}

// WithQueryParam adds a per-call query parameter.
func WithQueryParam(key, value string) RequestOption {
	return func(options *RequestOptions) {
		if options.Query == nil {
			options.Query = make(url.Values)
		}
		options.Query.Add(key, value)
	}
}

// WithPagination adds standard page and limit query parameters.
func WithPagination(page, limit int) RequestOption {
	return func(options *RequestOptions) {
		options.Page = page
		options.Limit = limit
	}
}

// WithTimeout bounds a single HTTP request.
func WithTimeout(timeout time.Duration) RequestOption {
	return func(options *RequestOptions) {
		options.Timeout = timeout
	}
}

// WithMaxRetries retries retryable HTTP responses up to maxRetries times.
func WithMaxRetries(maxRetries int) RequestOption {
	return func(options *RequestOptions) {
		options.MaxRetries = maxRetries
	}
}

// New returns a service client rooted at baseURL.
func New(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		HTTPClient: httpClient,
	}
}

// URL resolves a service-relative path against BaseURL.
func (c *Client) URL(path string) (string, error) {
	if c.BaseURL == "" {
		return "", fmt.Errorf("service: base URL is required")
	}
	if strings.TrimSpace(path) == "" {
		path = "/"
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	base, err := url.Parse(c.BaseURL + "/")
	if err != nil {
		return "", fmt.Errorf("service: invalid base URL: %w", err)
	}
	ref, err := url.Parse(strings.TrimLeft(path, "/"))
	if err != nil {
		return "", fmt.Errorf("service: invalid path: %w", err)
	}
	return base.ResolveReference(ref).String(), nil
}

// DoJSON sends a JSON request and decodes a JSON response into out.
func (c *Client) DoJSON(ctx context.Context, method string, path string, body any, out any) error {
	return c.DoJSONWithOptions(ctx, method, path, body, out)
}

// DoJSONWithOptions sends a JSON request with per-call options.
func (c *Client) DoJSONWithOptions(ctx context.Context, method string, path string, body any, out any, opts ...RequestOption) error {
	var options RequestOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	var requestBody []byte
	if body != nil {
		switch value := body.(type) {
		case []byte:
			requestBody = append([]byte(nil), value...)
		case json.RawMessage:
			requestBody = append([]byte(nil), value...)
		case io.Reader:
			payload, err := io.ReadAll(value)
			if err != nil {
				return fmt.Errorf("service: read request body: %w", err)
			}
			requestBody = payload
		default:
			payload, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("service: encode request: %w", err)
			}
			requestBody = payload
		}
	}

	requestURL, err := c.URL(path)
	if err != nil {
		return err
	}
	if options.Page > 0 || options.Limit > 0 {
		if options.Query == nil {
			options.Query = make(url.Values)
		}
		if options.Page > 0 {
			options.Query.Set("page", strconv.Itoa(options.Page))
		}
		if options.Limit > 0 {
			options.Query.Set("limit", strconv.Itoa(options.Limit))
		}
	}
	if len(options.Query) > 0 {
		parsed, err := url.Parse(requestURL)
		if err != nil {
			return fmt.Errorf("service: invalid request URL: %w", err)
		}
		query := parsed.Query()
		for key, values := range options.Query {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		parsed.RawQuery = query.Encode()
		requestURL = parsed.String()
	}

	attempts := options.MaxRetries + 1
	if attempts < 1 {
		attempts = 1
	}
	var lastResponseBody []byte
	var lastStatusCode int
	for attempt := 0; attempt < attempts; attempt++ {
		var reader io.Reader
		if requestBody != nil {
			reader = bytes.NewReader(requestBody)
		}
		req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
		if err != nil {
			return fmt.Errorf("service: build request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		for key, values := range options.Headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		if strings.TrimSpace(options.IdempotencyKey) != "" {
			req.Header.Set("Idempotency-Key", strings.TrimSpace(options.IdempotencyKey))
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("service: %s %s: %w", method, requestURL, err)
		}

		responseBody, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if readErr != nil {
			return fmt.Errorf("service: read response: %w", readErr)
		}
		if closeErr != nil {
			return fmt.Errorf("service: close response: %w", closeErr)
		}
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
			if out == nil || len(responseBody) == 0 {
				return nil
			}
			if raw, ok := out.(*json.RawMessage); ok {
				*raw = append((*raw)[:0], responseBody...)
				return nil
			}
			if bytes.Equal(bytes.TrimSpace(responseBody), []byte("null")) {
				return nil
			}
			if err := json.Unmarshal(responseBody, out); err != nil {
				return fmt.Errorf("service: decode response: %w", err)
			}
			return nil
		}
		lastResponseBody = responseBody
		lastStatusCode = resp.StatusCode
		if attempt+1 < attempts && isRetryableStatus(resp.StatusCode) {
			continue
		}
		break
	}
	return &HTTPError{
		Method:     method,
		URL:        requestURL,
		StatusCode: lastStatusCode,
		Body:       lastResponseBody,
	}
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

// HTTPError preserves the upstream response body for callers that need details.
type HTTPError struct {
	Method     string
	URL        string
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("service: %s %s returned %d: %s", e.Method, e.URL, e.StatusCode, strings.TrimSpace(string(e.Body)))
}
