from __future__ import annotations

import json
from dataclasses import dataclass
from typing import Any
from urllib.error import HTTPError as UrlHTTPError
from urllib.parse import parse_qsl, urlencode, urljoin, urlsplit, urlunsplit
from urllib.request import Request, urlopen

from .auth import TokenSource, bearer_authorization
from .models import DomainOperationsCatalog, ResolvedOperation


@dataclass
class HTTPError(Exception):
    method: str
    url: str
    status_code: int
    body: str

    def __str__(self) -> str:
        return f"2finance: {self.method} {self.url} returned {self.status_code}: {self.body}"


@dataclass(frozen=True)
class RequestOptions:
    headers: dict[str, str] | None = None
    idempotency_key: str | None = None
    query: dict[str, Any] | None = None
    timeout: float | None = None
    max_retries: int = 0
    page: int | None = None
    limit: int | None = None


class ServiceClient:
    def __init__(self, base_url: str, token_source: TokenSource | None = None):
        self.base_url = (base_url or "").strip().rstrip("/")
        self.token_source = token_source

    def url(self, path: str) -> str:
        if path.startswith("http://") or path.startswith("https://"):
            return path
        if not self.base_url:
            raise ValueError("base_url is required")
        return urljoin(f"{self.base_url}/", path.lstrip("/"))

    def request(self, method: str, path: str, body: Any = None, options: RequestOptions | None = None) -> Any:
        request_body = None
        headers = {"Accept": "application/json"}
        if body is not None:
            headers["Content-Type"] = "application/json"
            request_body = json.dumps(body).encode("utf-8")
        if self.token_source is not None:
            authorization = bearer_authorization(self.token_source.token())
            if authorization:
                headers["Authorization"] = authorization
        if options is not None:
            headers.update(options.headers or {})
            idempotency_key = (options.idempotency_key or "").strip()
            if idempotency_key:
                headers["Idempotency-Key"] = idempotency_key
        request_url = self.url(path)
        if options is not None and (options.query or options.page is not None or options.limit is not None):
            request_url = self._url_with_query(request_url, options)
        request = Request(request_url, data=request_body, headers=headers, method=method)
        timeout = options.timeout if options is not None else None
        max_retries = max(0, options.max_retries if options is not None else 0)
        for attempt in range(max_retries + 1):
            try:
                response_context = urlopen(request, timeout=timeout) if timeout is not None else urlopen(request)
                with response_context as response:
                    raw = response.read().decode("utf-8")
                break
            except UrlHTTPError as exc:
                raw = exc.read().decode("utf-8")
                if attempt < max_retries and _is_retryable_status(exc.code):
                    continue
                raise HTTPError(method, request_url, exc.code, raw) from exc
        return json.loads(raw) if raw else None

    def get(self, path: str, options: RequestOptions | None = None) -> Any:
        return self.request("GET", path, options=options)

    def post(self, path: str, body: Any = None, options: RequestOptions | None = None) -> Any:
        return self.request("POST", path, body, options)

    def put(self, path: str, body: Any = None, options: RequestOptions | None = None) -> Any:
        return self.request("PUT", path, body, options)

    def delete(self, path: str, options: RequestOptions | None = None) -> Any:
        return self.request("DELETE", path, options=options)

    def request_operation(
        self,
        operation: ResolvedOperation,
        body: Any = None,
        options: RequestOptions | None = None,
    ) -> Any:
        return self.request(operation.method, operation.path, body, options)

    def request_catalog_operation(
        self,
        catalog: DomainOperationsCatalog,
        domain_name: str,
        operation_name: str,
        path_params: dict[str, Any] | None = None,
        query: dict[str, Any] | None = None,
        body: Any = None,
        options: RequestOptions | None = None,
    ) -> Any:
        operation = catalog.resolve_operation(
            domain_name,
            operation_name,
            path_params=path_params,
            query=query,
        )
        return self.request_operation(operation, body=body, options=options)

    @staticmethod
    def _url_with_query(request_url: str, options: RequestOptions) -> str:
        parts = urlsplit(request_url)
        merged = dict(parse_qsl(parts.query, keep_blank_values=True))
        for key, value in (options.query or {}).items():
            if value is not None:
                merged[key] = str(value)
        if options.page is not None:
            merged["page"] = str(options.page)
        if options.limit is not None:
            merged["limit"] = str(options.limit)
        return urlunsplit((parts.scheme, parts.netloc, parts.path, urlencode(merged), parts.fragment))


def _is_retryable_status(status_code: int) -> bool:
    return status_code == 429 or status_code >= 500
