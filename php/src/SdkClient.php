<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

final class SdkClient
{
    public readonly AuthClient $auth;
    public readonly NetworkClient $network;
    public readonly AnalyticsClient $analytics;
    public readonly OrchestratorClient $orchestrator;
    public readonly MCPClient $mcp;
    public readonly TradingControlClient $tradingControl;
    public readonly MatchEngineClient $matchEngine;
    public readonly KeyStoreClient $keystore;
    public readonly HummingbotClient $hummingbot;
    public readonly WiseClient $wise;
    public readonly AirwallexClient $airwallex;
    public readonly PlannerClient $planner;

    public function __construct(
        public readonly SdkConfig $config,
        ?callable $transport = null,
    ) {
        $tokenSource = $config->tokenSource;
        $this->auth = new AuthClient(
            new ServiceClient($config->authUrl, $transport, $tokenSource),
            $config->authRealm,
            $config->authClientId,
            $config->authPhoneClientId,
        );
        $this->network = new NetworkClient(new ServiceClient($config->networkUrl, $transport, $tokenSource));
        $this->analytics = new AnalyticsClient(new ServiceClient($config->analyticsUrl, $transport, $tokenSource));
        $this->orchestrator = new OrchestratorClient(new ServiceClient($config->orchestratorUrl, $transport, $tokenSource));
        $this->mcp = new MCPClient(new ServiceClient($config->mcpUrl, $transport, $tokenSource));
        $this->tradingControl = new TradingControlClient(new ServiceClient($config->tradingControlUrl, $transport, $tokenSource));
        $this->matchEngine = new MatchEngineClient($config->matchEngineWsUrl);
        $this->keystore = new KeyStoreClient(new ServiceClient($config->keystoreUrl, $transport, $tokenSource));
        $this->hummingbot = new HummingbotClient(new ServiceClient($config->hummingbotUrl, $transport, $tokenSource));
        $this->wise = new WiseClient(new ProviderClient(new ServiceClient($config->wiseUrl, $transport, $tokenSource)));
        $this->airwallex = new AirwallexClient(new ProviderClient(new ServiceClient($config->airwallexUrl, $transport, $tokenSource)));
        $this->planner = new PlannerClient(
            $this->mcp,
            $this->orchestrator,
            $this->analytics,
            $this->tradingControl,
        );
    }

    public static function fromEnv(?callable $transport = null): self
    {
        return new self(SdkConfig::fromEnv(), $transport);
    }
}
