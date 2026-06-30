package analytics

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

func (c *Client) GetWithOptions(ctx context.Context, path string, out any, opts ...service.RequestOption) error {
	return c.service.DoJSONWithOptions(ctx, http.MethodGet, path, nil, out, opts...)
}

func (c *Client) Post(ctx context.Context, path string, body any, out any) error {
	return c.service.DoJSON(ctx, http.MethodPost, path, body, out)
}

func (c *Client) PostWithOptions(ctx context.Context, path string, body any, out any, opts ...service.RequestOption) error {
	return c.service.DoJSONWithOptions(ctx, http.MethodPost, path, body, out, opts...)
}

func (c *Client) Health(ctx context.Context) error {
	return c.Get(ctx, "/healthz", nil)
}

func (c *Client) CalculateTechnicalAnalysis(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/analytics/technical-analysis:calculate", request, &response)
	return response, err
}

func (c *Client) ListIndicators(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/analytics/indicators", &response)
	return response, err
}

func (c *Client) UpsertCandles(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/analytics/candles:upsert", request, &response)
	return response, err
}

func (c *Client) OptimizePortfolio(ctx context.Context, request any) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Post(ctx, "/portfolio-manager/optimizer", request, &response)
	return response, err
}

func (c *Client) Rankings(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/portfolio-manager/rankings", &response)
	return response, err
}

func (c *Client) Balances(ctx context.Context, accountID string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/portfolio-manager/balances/"+url.PathEscape(accountID), &response)
	return response, err
}

func (c *Client) BlackScholes(ctx context.Context, query string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/risk-manager/blackscholes?"+query, &response)
	return response, err
}

func (c *Client) Staking(ctx context.Context) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.Get(ctx, "/staking", &response)
	return response, err
}
