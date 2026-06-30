package twofinance

import (
	"net/http"

	"github.com/2Finance-Labs/2finance-sdk-client/analytics"
	"github.com/2Finance-Labs/2finance-sdk-client/auth"
	"github.com/2Finance-Labs/2finance-sdk-client/hummingbot"
	"github.com/2Finance-Labs/2finance-sdk-client/keystore"
	"github.com/2Finance-Labs/2finance-sdk-client/matchengine"
	"github.com/2Finance-Labs/2finance-sdk-client/mcp"
	"github.com/2Finance-Labs/2finance-sdk-client/network"
	"github.com/2Finance-Labs/2finance-sdk-client/orchestrator"
	"github.com/2Finance-Labs/2finance-sdk-client/planner"
	"github.com/2Finance-Labs/2finance-sdk-client/providers"
	"github.com/2Finance-Labs/2finance-sdk-client/tradingcontrol"
)

// Client aggregates the official 2Finance service clients.
type Client struct {
	Config         Config
	Auth           *auth.UserClient
	Network        *network.Client
	Analytics      *analytics.Client
	Orchestrator   *orchestrator.Client
	MCP            *mcp.Client
	Planner        *planner.Client
	TradingControl *tradingcontrol.Client
	MatchEngine    *matchengine.Client
	Hummingbot     *hummingbot.Client
	KeyStore       *keystore.Client
	Providers      *providers.Client
}

// New builds a single SDK client from explicit configuration.
func New(config Config) *Client {
	httpClient := authenticatedHTTPClient(config.HTTPClient, config.TokenSource)
	mcpClient := mcp.New(config.MCPURL, httpClient)
	orchestratorClient := orchestrator.New(config.OrchestratorURL, httpClient)
	tradingClient := tradingcontrol.New(config.TradingControlURL, httpClient)
	analyticsClient := analytics.New(config.AnalyticsURL, httpClient)

	return &Client{
		Config: config,
		Auth: auth.NewUserClient(auth.UserClientConfig{
			BaseURL:       config.AuthURL,
			Realm:         config.AuthRealm,
			ClientID:      config.AuthClientID,
			PhoneClientID: config.AuthPhoneClientID,
			HTTPClient:    httpClient,
		}),
		Network:        network.New(config.NetworkURL, httpClient),
		Analytics:      analyticsClient,
		Orchestrator:   orchestratorClient,
		MCP:            mcpClient,
		TradingControl: tradingClient,
		MatchEngine:    matchengine.New(config.MatchEngineWSURL, httpClient),
		Hummingbot:     hummingbot.New(config.HummingbotURL, httpClient),
		KeyStore:       keystore.New(config.KeyStoreURL, httpClient),
		Providers: providers.New(providers.Config{
			WiseURL:      config.WiseURL,
			AirwallexURL: config.AirwallexURL,
			HTTPClient:   httpClient,
		}),
		Planner: planner.New(planner.Config{
			MCP:            mcpClient,
			Orchestrator:   orchestratorClient,
			TradingControl: tradingClient,
			Analytics:      analyticsClient,
		}),
	}
}

// NewFromEnv builds a SDK client using the standard TWO_FINANCE_* variables.
func NewFromEnv() *Client {
	return New(ConfigFromEnv())
}

func authenticatedHTTPClient(base *http.Client, tokenSource auth.TokenSource) *http.Client {
	if base == nil {
		base = http.DefaultClient
	}
	if tokenSource == nil {
		return base
	}
	clone := *base
	clone.Transport = auth.AuthTransport{
		Source: tokenSource,
		Base:   base.Transport,
	}
	return &clone
}
