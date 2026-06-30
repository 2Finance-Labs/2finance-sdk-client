package twofinance

import (
	"net/http"
	"os"
	"strings"

	"github.com/2Finance-Labs/2finance-sdk-client/auth"
)

// Config defines service endpoints shared by all 2Finance clients.
type Config struct {
	AuthURL           string
	NetworkURL        string
	AnalyticsURL      string
	OrchestratorURL   string
	MCPURL            string
	TradingControlURL string
	MatchEngineWSURL  string
	KeyStoreURL       string
	HummingbotURL     string
	WiseURL           string
	AirwallexURL      string
	HTTPClient        *http.Client
	TokenSource       auth.TokenSource
	AuthRealm         string
	AuthClientID      string
	AuthPhoneClientID string
}

// ConfigFromEnv loads the standard 2Finance SDK environment variables.
func ConfigFromEnv() Config {
	return Config{
		AuthURL:           env("TWO_FINANCE_AUTH_URL"),
		NetworkURL:        env("TWO_FINANCE_NETWORK_URL"),
		AnalyticsURL:      env("TWO_FINANCE_ANALYTICS_URL"),
		OrchestratorURL:   env("TWO_FINANCE_ORCHESTRATOR_URL"),
		MCPURL:            env("TWO_FINANCE_MCP_URL"),
		TradingControlURL: env("TWO_FINANCE_TRADING_CONTROL_URL"),
		MatchEngineWSURL:  env("TWO_FINANCE_MATCHENGINE_WS_URL"),
		KeyStoreURL:       env("TWO_FINANCE_KEYSTORE_URL"),
		HummingbotURL:     env("TWO_FINANCE_HUMMINGBOT_URL"),
		WiseURL:           env("TWO_FINANCE_WISE_URL"),
		AirwallexURL:      env("TWO_FINANCE_AIRWALLEX_URL"),
		AuthRealm:         envDefault("TWO_FINANCE_AUTH_REALM", "2finance"),
		AuthClientID:      envDefault("TWO_FINANCE_AUTH_CLIENT_ID", "2finance-network"),
		AuthPhoneClientID: envDefault("TWO_FINANCE_AUTH_PHONE_CLIENT_ID", "2finance-network-phone"),
	}
}

func (c Config) ServiceURL(domain string) string {
	switch serviceKey(domain) {
	case "auth":
		return c.AuthURL
	case "network":
		return c.NetworkURL
	case "analytics":
		return c.AnalyticsURL
	case "orchestrator":
		return c.OrchestratorURL
	case "mcp", "planner":
		return c.MCPURL
	case "tradingcontrol":
		return c.TradingControlURL
	case "matchengine":
		return c.MatchEngineWSURL
	case "keystore":
		return c.KeyStoreURL
	case "hummingbot":
		return c.HummingbotURL
	case "wise":
		return c.WiseURL
	case "airwallex":
		return c.AirwallexURL
	default:
		return ""
	}
}

func (c Config) ServiceURLs() map[string]string {
	urls := make(map[string]string)
	for _, service := range DefaultServiceCatalog().Services {
		if url := c.ServiceURL(service.Name); url != "" {
			urls[service.Name] = url
		}
	}
	return urls
}

func (c Config) ConfiguredServices() []ConfiguredServiceEntry {
	services := make([]ConfiguredServiceEntry, 0)
	for _, service := range DefaultServiceCatalog().Services {
		if url := c.ServiceURL(service.Name); url != "" {
			services = append(services, ConfiguredServiceEntry{
				Name: service.Name,
				Env:  service.Env,
				URL:  url,
			})
		}
	}
	return services
}

func (c Config) MissingServiceURLs() []ServiceCatalogEntry {
	services := make([]ServiceCatalogEntry, 0)
	for _, service := range DefaultServiceCatalog().Services {
		if c.ServiceURL(service.Name) == "" {
			services = append(services, service)
		}
	}
	return services
}

func env(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func envDefault(key string, fallback string) string {
	if value := env(key); value != "" {
		return value
	}
	return fallback
}

func serviceKey(domain string) string {
	return strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(strings.TrimSpace(domain)))
}
