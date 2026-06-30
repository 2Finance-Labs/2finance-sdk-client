package com.twofinance.sdk;

import java.net.http.HttpClient;

public final class SdkClient {
    public final SDKConfig config;
    public final DomainClients.AuthClient auth;
    public final DomainClients.NetworkClient network;
    public final DomainClients.AnalyticsClient analytics;
    public final DomainClients.OrchestratorClient orchestrator;
    public final DomainClients.MCPClient mcp;
    public final DomainClients.TradingControlClient tradingControl;
    public final DomainClients.MatchEngineClient matchEngine;
    public final DomainClients.KeyStoreClient keystore;
    public final DomainClients.HummingbotClient hummingbot;
    public final DomainClients.WiseClient wise;
    public final DomainClients.AirwallexClient airwallex;
    public final DomainClients.PlannerClient planner;

    public SdkClient(SDKConfig config) {
        this(config, HttpClient.newHttpClient());
    }

    public SdkClient(SDKConfig config, HttpClient httpClient) {
        this.config = config;
        TokenSource tokenSource = config.tokenSource;
        this.auth = new DomainClients.AuthClient(config.authUrl, httpClient, tokenSource, config);
        this.network = new DomainClients.NetworkClient(config.networkUrl, httpClient, tokenSource);
        this.analytics = new DomainClients.AnalyticsClient(config.analyticsUrl, httpClient, tokenSource);
        this.orchestrator = new DomainClients.OrchestratorClient(config.orchestratorUrl, httpClient, tokenSource);
        this.mcp = new DomainClients.MCPClient(config.mcpUrl, httpClient, tokenSource);
        this.tradingControl = new DomainClients.TradingControlClient(config.tradingControlUrl, httpClient, tokenSource);
        this.matchEngine = new DomainClients.MatchEngineClient(config.matchEngineWsUrl);
        this.keystore = new DomainClients.KeyStoreClient(config.keyStoreUrl, httpClient, tokenSource);
        this.hummingbot = new DomainClients.HummingbotClient(config.hummingbotUrl, httpClient, tokenSource);
        this.wise = new DomainClients.WiseClient(config.wiseUrl, httpClient, tokenSource);
        this.airwallex = new DomainClients.AirwallexClient(config.airwallexUrl, httpClient, tokenSource);
        this.planner = new DomainClients.PlannerClient(this.mcp, this.orchestrator, this.analytics, this.tradingControl);
    }

    public static SdkClient fromEnv() {
        return new SdkClient(SDKConfig.fromEnv());
    }
}
