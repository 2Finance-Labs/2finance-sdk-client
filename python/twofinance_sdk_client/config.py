from __future__ import annotations

import os
from dataclasses import dataclass
from typing import Mapping

from .models import DEFAULT_SERVICE_CATALOG, ConfiguredServiceEntry


@dataclass
class SDKConfig:
    auth_url: str = ""
    network_url: str = ""
    analytics_url: str = ""
    orchestrator_url: str = ""
    mcp_url: str = ""
    trading_control_url: str = ""
    matchengine_ws_url: str = ""
    keystore_url: str = ""
    hummingbot_url: str = ""
    wise_url: str = ""
    airwallex_url: str = ""
    auth_realm: str = "2finance"
    auth_client_id: str = "2finance-network"
    auth_phone_client_id: str = "2finance-network-phone"

    def service_url(self, domain: str) -> str:
        return {
            "auth": self.auth_url,
            "network": self.network_url,
            "analytics": self.analytics_url,
            "orchestrator": self.orchestrator_url,
            "mcp": self.mcp_url,
            "planner": self.mcp_url,
            "tradingcontrol": self.trading_control_url,
            "matchengine": self.matchengine_ws_url,
            "keystore": self.keystore_url,
            "hummingbot": self.hummingbot_url,
            "wise": self.wise_url,
            "airwallex": self.airwallex_url,
        }.get(_service_key(domain), "")

    def service_urls(self) -> dict[str, str]:
        urls: dict[str, str] = {}
        for service in DEFAULT_SERVICE_CATALOG.services:
            url = self.service_url(service.name)
            if url:
                urls[service.name] = url
        return urls

    def configured_services(self) -> list[ConfiguredServiceEntry]:
        services: list[ConfiguredServiceEntry] = []
        for service in DEFAULT_SERVICE_CATALOG.services:
            url = self.service_url(service.name)
            if url:
                services.append(ConfiguredServiceEntry(service.name, service.env, url))
        return services

    def missing_service_urls(self):
        return [
            service
            for service in DEFAULT_SERVICE_CATALOG.services
            if not self.service_url(service.name)
        ]


def config_from_env(env: Mapping[str, str] | None = None) -> SDKConfig:
    source = env if env is not None else os.environ
    return SDKConfig(
        auth_url=source.get("TWO_FINANCE_AUTH_URL", ""),
        network_url=source.get("TWO_FINANCE_NETWORK_URL", ""),
        analytics_url=source.get("TWO_FINANCE_ANALYTICS_URL", ""),
        orchestrator_url=source.get("TWO_FINANCE_ORCHESTRATOR_URL", ""),
        mcp_url=source.get("TWO_FINANCE_MCP_URL", ""),
        trading_control_url=source.get("TWO_FINANCE_TRADING_CONTROL_URL", ""),
        matchengine_ws_url=source.get("TWO_FINANCE_MATCHENGINE_WS_URL", ""),
        keystore_url=source.get("TWO_FINANCE_KEYSTORE_URL", ""),
        hummingbot_url=source.get("TWO_FINANCE_HUMMINGBOT_URL", ""),
        wise_url=source.get("TWO_FINANCE_WISE_URL", ""),
        airwallex_url=source.get("TWO_FINANCE_AIRWALLEX_URL", ""),
        auth_realm=source.get("TWO_FINANCE_AUTH_REALM", "2finance"),
        auth_client_id=source.get("TWO_FINANCE_AUTH_CLIENT_ID", "2finance-network"),
        auth_phone_client_id=source.get(
            "TWO_FINANCE_AUTH_PHONE_CLIENT_ID", "2finance-network-phone"
        ),
    )


def _service_key(domain: str) -> str:
    return "".join(ch for ch in domain.strip().lower() if ch not in "-_ ")
