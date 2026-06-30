package com.twofinance.sdk;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

public final class SDKConfig {
    public String authUrl = "";
    public String networkUrl = "";
    public String analyticsUrl = "";
    public String orchestratorUrl = "";
    public String mcpUrl = "";
    public String tradingControlUrl = "";
    public String matchEngineWsUrl = "";
    public String keyStoreUrl = "";
    public String hummingbotUrl = "";
    public String wiseUrl = "";
    public String airwallexUrl = "";
    public String authRealm = "2finance";
    public String authClientId = "2finance-network";
    public String authPhoneClientId = "2finance-network-phone";
    public TokenSource tokenSource;

    public static SDKConfig fromEnv() {
        return fromEnv(System.getenv());
    }

    public static SDKConfig fromEnv(Map<String, String> env) {
        SDKConfig config = new SDKConfig();
        config.authUrl = env.getOrDefault("TWO_FINANCE_AUTH_URL", "");
        config.networkUrl = env.getOrDefault("TWO_FINANCE_NETWORK_URL", "");
        config.analyticsUrl = env.getOrDefault("TWO_FINANCE_ANALYTICS_URL", "");
        config.orchestratorUrl = env.getOrDefault("TWO_FINANCE_ORCHESTRATOR_URL", "");
        config.mcpUrl = env.getOrDefault("TWO_FINANCE_MCP_URL", "");
        config.tradingControlUrl = env.getOrDefault("TWO_FINANCE_TRADING_CONTROL_URL", "");
        config.matchEngineWsUrl = env.getOrDefault("TWO_FINANCE_MATCHENGINE_WS_URL", "");
        config.keyStoreUrl = env.getOrDefault("TWO_FINANCE_KEYSTORE_URL", "");
        config.hummingbotUrl = env.getOrDefault("TWO_FINANCE_HUMMINGBOT_URL", "");
        config.wiseUrl = env.getOrDefault("TWO_FINANCE_WISE_URL", "");
        config.airwallexUrl = env.getOrDefault("TWO_FINANCE_AIRWALLEX_URL", "");
        config.authRealm = env.getOrDefault("TWO_FINANCE_AUTH_REALM", "2finance");
        config.authClientId = env.getOrDefault("TWO_FINANCE_AUTH_CLIENT_ID", "2finance-network");
        config.authPhoneClientId = env.getOrDefault(
                "TWO_FINANCE_AUTH_PHONE_CLIENT_ID",
                "2finance-network-phone"
        );
        return config;
    }

    public String serviceUrl(String domain) {
        switch (serviceKey(domain)) {
            case "auth":
                return authUrl;
            case "network":
                return networkUrl;
            case "analytics":
                return analyticsUrl;
            case "orchestrator":
                return orchestratorUrl;
            case "mcp":
            case "planner":
                return mcpUrl;
            case "tradingcontrol":
                return tradingControlUrl;
            case "matchengine":
                return matchEngineWsUrl;
            case "keystore":
                return keyStoreUrl;
            case "hummingbot":
                return hummingbotUrl;
            case "wise":
                return wiseUrl;
            case "airwallex":
                return airwallexUrl;
            default:
                return "";
        }
    }

    public Map<String, String> serviceUrls() {
        Map<String, String> urls = new LinkedHashMap<>();
        for (Models.ServiceCatalogEntry service : SDKMetadata.serviceCatalog().services()) {
            String url = serviceUrl(service.name());
            if (!url.isBlank()) {
                urls.put(service.name(), url);
            }
        }
        return urls;
    }

    public List<Models.ConfiguredServiceEntry> configuredServices() {
        List<Models.ConfiguredServiceEntry> services = new ArrayList<>();
        for (Models.ServiceCatalogEntry service : SDKMetadata.serviceCatalog().services()) {
            String url = serviceUrl(service.name());
            if (!url.isBlank()) {
                services.add(new Models.ConfiguredServiceEntry(service.name(), service.env(), url));
            }
        }
        return services;
    }

    public List<Models.ServiceCatalogEntry> missingServiceUrls() {
        List<Models.ServiceCatalogEntry> services = new ArrayList<>();
        for (Models.ServiceCatalogEntry service : SDKMetadata.serviceCatalog().services()) {
            if (serviceUrl(service.name()).isBlank()) {
                services.add(service);
            }
        }
        return services;
    }

    private static String serviceKey(String domain) {
        return domain == null ? "" : domain.trim().toLowerCase().replace("-", "").replace("_", "").replace(" ", "");
    }
}
