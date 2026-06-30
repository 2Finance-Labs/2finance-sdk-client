from __future__ import annotations

from dataclasses import dataclass

from .auth import TokenSource
from .config import SDKConfig, config_from_env
from .domains import (
    AirwallexClient,
    AnalyticsClient,
    AuthClient,
    HummingbotClient,
    KeyStoreClient,
    MatchEngineClient,
    MCPClient,
    NetworkClient,
    OrchestratorClient,
    PlannerClient,
    TradingControlClient,
    WiseClient,
)


@dataclass
class ProvidersClient:
    wise: WiseClient
    airwallex: AirwallexClient


class TwoFinanceClient:
    def __init__(self, config: SDKConfig | None = None, token_source: TokenSource | None = None):
        self.config = config or SDKConfig()
        self.auth = AuthClient(
            self.config.auth_url,
            realm=self.config.auth_realm,
            client_id=self.config.auth_client_id,
            phone_client_id=self.config.auth_phone_client_id,
            token_source=token_source,
        )
        self.network = NetworkClient(self.config.network_url, token_source=token_source)
        self.analytics = AnalyticsClient(self.config.analytics_url, token_source=token_source)
        self.orchestrator = OrchestratorClient(self.config.orchestrator_url, token_source=token_source)
        self.mcp = MCPClient(self.config.mcp_url, token_source=token_source)
        self.trading_control = TradingControlClient(self.config.trading_control_url, token_source=token_source)
        self.matchengine = MatchEngineClient(self.config.matchengine_ws_url)
        self.keystore = KeyStoreClient(self.config.keystore_url, token_source=token_source)
        self.hummingbot = HummingbotClient(self.config.hummingbot_url, token_source=token_source)
        self.providers = ProvidersClient(
            wise=WiseClient(self.config.wise_url, token_source=token_source),
            airwallex=AirwallexClient(self.config.airwallex_url, token_source=token_source),
        )
        self.planner = PlannerClient(
            mcp=self.mcp,
            orchestrator=self.orchestrator,
            analytics=self.analytics,
            trading_control=self.trading_control,
        )

    @classmethod
    def from_env(cls, token_source: TokenSource | None = None) -> "TwoFinanceClient":
        return cls(config_from_env(), token_source=token_source)
