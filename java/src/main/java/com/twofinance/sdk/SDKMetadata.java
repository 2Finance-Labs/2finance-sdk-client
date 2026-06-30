package com.twofinance.sdk;

import java.util.List;

public final class SDKMetadata {
    public static final String SDK_NAME = "2finance-sdk-client";
    public static final String SDK_VERSION = "0.1.0";

    public static Models.ServiceCatalog serviceCatalog() {
        return new Models.ServiceCatalog(List.of(
                new Models.ServiceCatalogEntry("auth", "TWO_FINANCE_AUTH_URL"),
                new Models.ServiceCatalogEntry("network", "TWO_FINANCE_NETWORK_URL"),
                new Models.ServiceCatalogEntry("analytics", "TWO_FINANCE_ANALYTICS_URL"),
                new Models.ServiceCatalogEntry("orchestrator", "TWO_FINANCE_ORCHESTRATOR_URL"),
                new Models.ServiceCatalogEntry("mcp", "TWO_FINANCE_MCP_URL"),
                new Models.ServiceCatalogEntry("planner", "TWO_FINANCE_MCP_URL"),
                new Models.ServiceCatalogEntry("tradingcontrol", "TWO_FINANCE_TRADING_CONTROL_URL"),
                new Models.ServiceCatalogEntry("matchengine", "TWO_FINANCE_MATCHENGINE_WS_URL"),
                new Models.ServiceCatalogEntry("keystore", "TWO_FINANCE_KEYSTORE_URL"),
                new Models.ServiceCatalogEntry("hummingbot", "TWO_FINANCE_HUMMINGBOT_URL"),
                new Models.ServiceCatalogEntry("wise", "TWO_FINANCE_WISE_URL"),
                new Models.ServiceCatalogEntry("airwallex", "TWO_FINANCE_AIRWALLEX_URL")
        ));
    }

    private SDKMetadata() {}
}
