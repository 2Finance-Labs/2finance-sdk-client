from __future__ import annotations

import json
import time
from dataclasses import dataclass, field
from typing import Protocol
from urllib.parse import urlencode
from urllib.request import Request, urlopen


class TokenSource(Protocol):
    def token(self) -> str:
        """Return a bearer access token."""


def bearer_authorization(access_token: str) -> str:
    trimmed = (access_token or "").strip()
    if not trimmed:
        return ""
    if trimmed.lower().startswith("bearer "):
        return trimmed
    return f"Bearer {trimmed}"


@dataclass
class StaticTokenSource:
    access_token: str

    def token(self) -> str:
        return self.access_token


@dataclass
class ClientCredentialsTokenSource:
    token_url: str
    client_id: str
    client_secret: str
    scopes: list[str] = field(default_factory=list)
    expiry_skew_seconds: int = 30
    _access_token: str = ""
    _expires_at: float = 0.0

    def token(self) -> str:
        now = time.time()
        if self._access_token and now < self._expires_at - self.expiry_skew_seconds:
            return self._access_token
        if not self.token_url or not self.client_id or not self.client_secret:
            raise ValueError("token_url, client_id and client_secret are required")
        payload = {
            "grant_type": "client_credentials",
            "client_id": self.client_id,
            "client_secret": self.client_secret,
        }
        if self.scopes:
            payload["scope"] = " ".join(self.scopes)
        body = urlencode(payload).encode("utf-8")
        request = Request(
            self.token_url,
            data=body,
            headers={
                "Accept": "application/json",
                "Content-Type": "application/x-www-form-urlencoded",
            },
            method="POST",
        )
        with urlopen(request) as response:
            data = json.loads(response.read().decode("utf-8"))
        access_token = data.get("access_token")
        if not access_token:
            raise ValueError("token response missing access_token")
        self._access_token = access_token
        self._expires_at = now + int(data.get("expires_in", 300))
        return self._access_token
