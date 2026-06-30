package providers

import (
	"context"
	"net/http"
	"net/url"

	"github.com/2Finance-Labs/2finance-sdk-client/internal/service"
)

type Config struct {
	WiseURL      string
	AirwallexURL string
	HTTPClient   *http.Client
}

type Client struct {
	Wise      *WiseClient
	Airwallex *AirwallexClient
}

type ProviderClient struct {
	service *service.Client
}

type WiseClient struct {
	*ProviderClient
}

type AirwallexClient struct {
	*ProviderClient
}

func New(config Config) *Client {
	return &Client{
		Wise:      &WiseClient{ProviderClient: &ProviderClient{service: service.New(config.WiseURL, config.HTTPClient)}},
		Airwallex: &AirwallexClient{ProviderClient: &ProviderClient{service: service.New(config.AirwallexURL, config.HTTPClient)}},
	}
}

func (c *ProviderClient) Get(ctx context.Context, path string, out any) error {
	return c.service.DoJSON(ctx, http.MethodGet, path, nil, out)
}

func (c *ProviderClient) Post(ctx context.Context, path string, body any, out any) error {
	return c.service.DoJSON(ctx, http.MethodPost, path, body, out)
}

func (c *WiseClient) Profiles(ctx context.Context) (any, error) {
	var out any
	err := c.Get(ctx, "/v1/profiles", &out)
	return out, err
}

func (c *WiseClient) Profile(ctx context.Context, profileID string) (any, error) {
	var out any
	err := c.Get(ctx, "/v1/profiles/"+url.PathEscape(profileID), &out)
	return out, err
}

func (c *WiseClient) CreateQuote(ctx context.Context, profileID string, request any) (any, error) {
	var out any
	err := c.Post(ctx, "/v3/profiles/"+url.PathEscape(profileID)+"/quotes", request, &out)
	return out, err
}

func (c *WiseClient) CreateTransfer(ctx context.Context, request any) (any, error) {
	var out any
	err := c.Post(ctx, "/v1/transfers", request, &out)
	return out, err
}

func (c *AirwallexClient) Accounts(ctx context.Context) (any, error) {
	var out any
	err := c.Get(ctx, "/api/v1/accounts", &out)
	return out, err
}

func (c *AirwallexClient) Payments(ctx context.Context) (any, error) {
	var out any
	err := c.Get(ctx, "/api/v1/payments", &out)
	return out, err
}

func (c *AirwallexClient) CreatePayment(ctx context.Context, request any) (any, error) {
	var out any
	err := c.Post(ctx, "/api/v1/payments", request, &out)
	return out, err
}

func (c *AirwallexClient) Beneficiaries(ctx context.Context) (any, error) {
	var out any
	err := c.Get(ctx, "/api/v1/beneficiaries", &out)
	return out, err
}

func (c *AirwallexClient) CreateBeneficiary(ctx context.Context, request any) (any, error) {
	var out any
	err := c.Post(ctx, "/api/v1/beneficiaries", request, &out)
	return out, err
}
