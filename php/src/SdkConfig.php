<?php

declare(strict_types=1);

namespace TwoFinance\Sdk;

final class SdkConfig
{
    public string $authUrl = '';
    public string $networkUrl = '';
    public string $analyticsUrl = '';
    public string $orchestratorUrl = '';
    public string $mcpUrl = '';
    public string $tradingControlUrl = '';
    public string $matchEngineWsUrl = '';
    public string $keystoreUrl = '';
    public string $hummingbotUrl = '';
    public string $wiseUrl = '';
    public string $airwallexUrl = '';
    public string $authRealm = '2finance';
    public string $authClientId = '2finance-network';
    public string $authPhoneClientId = '2finance-network-phone';
    public ?TokenSource $tokenSource = null;

    /**
     * @param array<string,string> $env
     */
    public static function fromEnv(array $env = []): self
    {
        $source = $env ?: $_ENV + $_SERVER;
        $config = new self();
        $config->authUrl = $source['TWO_FINANCE_AUTH_URL'] ?? '';
        $config->networkUrl = $source['TWO_FINANCE_NETWORK_URL'] ?? '';
        $config->analyticsUrl = $source['TWO_FINANCE_ANALYTICS_URL'] ?? '';
        $config->orchestratorUrl = $source['TWO_FINANCE_ORCHESTRATOR_URL'] ?? '';
        $config->mcpUrl = $source['TWO_FINANCE_MCP_URL'] ?? '';
        $config->tradingControlUrl = $source['TWO_FINANCE_TRADING_CONTROL_URL'] ?? '';
        $config->matchEngineWsUrl = $source['TWO_FINANCE_MATCHENGINE_WS_URL'] ?? '';
        $config->keystoreUrl = $source['TWO_FINANCE_KEYSTORE_URL'] ?? '';
        $config->hummingbotUrl = $source['TWO_FINANCE_HUMMINGBOT_URL'] ?? '';
        $config->wiseUrl = $source['TWO_FINANCE_WISE_URL'] ?? '';
        $config->airwallexUrl = $source['TWO_FINANCE_AIRWALLEX_URL'] ?? '';
        $config->authRealm = $source['TWO_FINANCE_AUTH_REALM'] ?? '2finance';
        $config->authClientId = $source['TWO_FINANCE_AUTH_CLIENT_ID'] ?? '2finance-network';
        $config->authPhoneClientId = $source['TWO_FINANCE_AUTH_PHONE_CLIENT_ID'] ?? '2finance-network-phone';
        return $config;
    }

    public function serviceUrl(string $domain): string
    {
        return match (self::serviceKey($domain)) {
            'auth' => $this->authUrl,
            'network' => $this->networkUrl,
            'analytics' => $this->analyticsUrl,
            'orchestrator' => $this->orchestratorUrl,
            'mcp', 'planner' => $this->mcpUrl,
            'tradingcontrol' => $this->tradingControlUrl,
            'matchengine' => $this->matchEngineWsUrl,
            'keystore' => $this->keystoreUrl,
            'hummingbot' => $this->hummingbotUrl,
            'wise' => $this->wiseUrl,
            'airwallex' => $this->airwallexUrl,
            default => '',
        };
    }

    /** @return array<string,string> */
    public function serviceUrls(): array
    {
        $urls = [];
        foreach (Metadata::SERVICE_CATALOG as $service) {
            $url = $this->serviceUrl((string) $service['name']);
            if ($url !== '') {
                $urls[(string) $service['name']] = $url;
            }
        }
        return $urls;
    }

    /** @return array<int,ConfiguredServiceEntry> */
    public function configuredServices(): array
    {
        $services = [];
        foreach (Metadata::SERVICE_CATALOG as $service) {
            $url = $this->serviceUrl((string) $service['name']);
            if ($url !== '') {
                $services[] = new ConfiguredServiceEntry((string) $service['name'], (string) $service['env'], $url);
            }
        }
        return $services;
    }

    /** @return array<int,ServiceCatalogEntry> */
    public function missingServiceUrls(): array
    {
        $services = [];
        foreach (Metadata::SERVICE_CATALOG as $service) {
            if ($this->serviceUrl((string) $service['name']) === '') {
                $services[] = new ServiceCatalogEntry((string) $service['name'], (string) $service['env']);
            }
        }
        return $services;
    }

    private static function serviceKey(string $domain): string
    {
        return str_replace(['-', '_', ' '], '', strtolower(trim($domain)));
    }
}
