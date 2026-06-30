package twofinance

// DefaultServiceCatalog returns the canonical 2Finance SDK service catalog.
func DefaultServiceCatalog() ServiceCatalog {
	return ServiceCatalog{Services: []ServiceCatalogEntry{
		{Name: "auth", Env: "TWO_FINANCE_AUTH_URL"},
		{Name: "network", Env: "TWO_FINANCE_NETWORK_URL"},
		{Name: "analytics", Env: "TWO_FINANCE_ANALYTICS_URL"},
		{Name: "orchestrator", Env: "TWO_FINANCE_ORCHESTRATOR_URL"},
		{Name: "mcp", Env: "TWO_FINANCE_MCP_URL"},
		{Name: "planner", Env: "TWO_FINANCE_MCP_URL"},
		{Name: "tradingcontrol", Env: "TWO_FINANCE_TRADING_CONTROL_URL"},
		{Name: "matchengine", Env: "TWO_FINANCE_MATCHENGINE_WS_URL"},
		{Name: "keystore", Env: "TWO_FINANCE_KEYSTORE_URL"},
		{Name: "hummingbot", Env: "TWO_FINANCE_HUMMINGBOT_URL"},
		{Name: "wise", Env: "TWO_FINANCE_WISE_URL"},
		{Name: "airwallex", Env: "TWO_FINANCE_AIRWALLEX_URL"},
	}}
}
