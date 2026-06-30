package orchestrator

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

func (c *Client) Delete(ctx context.Context, path string, out any) error {
	return c.service.DoJSON(ctx, http.MethodDelete, path, nil, out)
}

func (c *Client) Catalog(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/catalog/packages", &response)
	return response, err
}

func (c *Client) CreateSession(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/v1/mcphost/sessions", request, &response)
	return response, err
}

func (c *Client) SendMessage(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/v1/mcphost/messages", request, &response)
	return response, err
}

func (c *Client) Tools(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/tools", &response)
	return response, err
}

func (c *Client) Prompts(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/prompts", &response)
	return response, err
}

func (c *Client) Resources(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/resources", &response)
	return response, err
}

func (c *Client) Providers(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/providers", &response)
	return response, err
}

func (c *Client) Approvals(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/approvals", &response)
	return response, err
}

func (c *Client) Observability(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/v1/mcphost/observability", &response)
	return response, err
}

func (c *Client) DeleteSession(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Delete(ctx, "/v1/mcphost/sessions/"+url.PathEscape(id), &response)
	return response, err
}
