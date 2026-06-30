package hummingbot

import (
	"context"
	"encoding/json"
	"net/http"

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

func (c *Client) Assets(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/assets", &response)
	return response, err
}

func (c *Client) Symbols(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/symbols", &response)
	return response, err
}

func (c *Client) Books(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/books", &response)
	return response, err
}

func (c *Client) Wallets(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/wallets", &response)
	return response, err
}

func (c *Client) Balances(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/balances", &response)
	return response, err
}

func (c *Client) ConnectorConfig(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/api/v1/connectors/2finance/config", request, &response)
	return response, err
}

func (c *Client) Latest(ctx context.Context, resource string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/api/v1/"+resource+"/latest", &response)
	return response, err
}
