package tradingcontrol

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

func (c *Client) Put(ctx context.Context, path string, body any, out any) error {
	return c.service.DoJSON(ctx, http.MethodPut, path, body, out)
}

func (c *Client) Health(ctx context.Context) error {
	return c.Get(ctx, "/health", nil)
}

func (c *Client) CreateRobot(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/robots", request, &response)
	return response, err
}

func (c *Client) Robots(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/robots", &response)
	return response, err
}

func (c *Client) Robot(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/robots/"+url.PathEscape(id), &response)
	return response, err
}

func (c *Client) StartRobot(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/robots/"+url.PathEscape(id)+":start", nil, &response)
	return response, err
}

func (c *Client) PauseRobot(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/robots/"+url.PathEscape(id)+":pause", nil, &response)
	return response, err
}

func (c *Client) ResumeRobot(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/robots/"+url.PathEscape(id)+":resume", nil, &response)
	return response, err
}

func (c *Client) StopRobot(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/robots/"+url.PathEscape(id)+":stop", nil, &response)
	return response, err
}

func (c *Client) RiskPolicy(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/robots/"+url.PathEscape(id)+"/risk-policy", &response)
	return response, err
}

func (c *Client) SetRiskPolicy(ctx context.Context, id string, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Put(ctx, "/robots/"+url.PathEscape(id)+"/risk-policy", request, &response)
	return response, err
}

func (c *Client) RiskView(ctx context.Context, id string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/risk-view/"+url.PathEscape(id), &response)
	return response, err
}

func (c *Client) Strategies(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/strategies", &response)
	return response, err
}

func (c *Client) CreateStrategy(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/strategies", request, &response)
	return response, err
}

func (c *Client) Directives(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/directives", &response)
	return response, err
}

func (c *Client) CreateDirective(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/directives", request, &response)
	return response, err
}

func (c *Client) Audit(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/audit", &response)
	return response, err
}

func (c *Client) Activity(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/activity", &response)
	return response, err
}

func (c *Client) MCPTools(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/mcp/tools", &response)
	return response, err
}
