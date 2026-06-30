package keystore

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/2Finance-Labs/2finance-sdk-client/internal/service"
)

type Client struct {
	service *service.Client
}

func New(baseURL string, httpClient *http.Client) *Client {
	return &Client{service: service.New(baseURL, httpClient)}
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.service.DoJSON(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, body any, out any) error {
	return c.service.DoJSON(ctx, http.MethodPost, path, body, out)
}

func (c *Client) Health(ctx context.Context) error {
	return c.Get(ctx, "/healthz", nil)
}

func (c *Client) Ready(ctx context.Context) error {
	return c.Get(ctx, "/readyz", nil)
}

func (c *Client) StartKeygen(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/keystore/keygen/start", request, &response)
	return response, err
}

func (c *Client) KeygenSignature(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/keystore/keygen/signature", request, &response)
	return response, err
}

func (c *Client) StartSigning(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/keystore/signing/start", request, &response)
	return response, err
}

func (c *Client) SigningSignature(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/keystore/signing/signature", request, &response)
	return response, err
}

func (c *Client) StartResharing(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/keystore/resharing/start", request, &response)
	return response, err
}

func (c *Client) Keys(ctx context.Context, userPublicKey string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/keystore/keys/"+url.PathEscape(userPublicKey), &response)
	return response, err
}

func (c *Client) Signatures(ctx context.Context, userPublicKey string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/keystore/signatures/"+url.PathEscape(userPublicKey), &response)
	return response, err
}

func (c *Client) Metrics(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/keystore/tss/metrics", &response)
	return response, err
}
