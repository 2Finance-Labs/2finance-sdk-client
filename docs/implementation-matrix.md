# SDK Implementation Matrix

This matrix tracks the current implementation surface across the language SDKs.
It is intentionally tied to files, examples, and checks that the repository can
validate.

## Language Coverage

| Language | Client surface | Auth/token source | Request options | Error handling | Package metadata | Examples | Local/CI check |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Go | Complete first implementation with legacy blockchain compatibility and domain clients. | `auth.ClientCredentialsTokenSource`, static token source, bearer injection. | Public `twofinance.With*` options plus domain generic option calls. | Public `twofinance.HTTPError`. | `go/go.mod`, `twofinance.SDKName`, `twofinance.SDKVersion`, `DefaultServiceCatalog`. | quickstart, request options, client credentials, error handling. | `make -C go test-fast`. |
| Dart | Unified SDK plus migrated auth and blockchain packages. | `ClientCredentialsTokenSource`, static token source, injected transport. | `RequestOptions` with headers, idempotency, query, pagination, timeout, retries. | `ServiceException`. | `dart/sdk/pubspec.yaml`, `sdkName`, `sdkVersion`, `defaultServiceCatalog`. | quickstart, request options, client credentials, error handling. | `dart pub get && dart analyze && dart test`. |
| TypeScript | Typed client surface for all public SDK domains. | `ClientCredentialsTokenSource`, `StaticTokenSource`, bearer injection. | `RequestOptions` with headers, idempotency, query, pagination, timeout, retries. | `ServiceError`. | `typescript/package.json`, `SDK_NAME`, `SDK_VERSION`, `SERVICE_CATALOG`. | quickstart, request options, client credentials, error handling. | `npm test`. |
| JavaScript | Dependency-free Node 18+ client mirroring TypeScript. | `ClientCredentialsTokenSource`, `StaticTokenSource`, bearer injection. | Plain object request options with the same shape as TypeScript. | `ServiceError`. | `javascript/package.json`, `SDK_NAME`, `SDK_VERSION`, `SERVICE_CATALOG`. | quickstart, request options, client credentials, error handling. | `npm test`. |
| Python | Dependency-free standard-library client. | `ClientCredentialsTokenSource`, `StaticTokenSource`, bearer injection. | `RequestOptions` dataclass with headers, idempotency, query, pagination, timeout, retries. | `HTTPError`. | `python/pyproject.toml`, `SDK_NAME`, `__version__`, `DEFAULT_SERVICE_CATALOG`. | quickstart, request options, client credentials, error handling. | `PYTHONPATH=. python3 -m unittest discover -s tests -v`. |
| PHP | Dependency-free PHP 8.1+ client with injectable transport. | `ClientCredentialsTokenSource`, `StaticTokenSource`, bearer injection. | `RequestOptions` with headers, idempotency, query, pagination, timeout, retries. | `ServiceException`. | `php/composer.json`, `Metadata::SDK_NAME`, `Metadata::SDK_VERSION`, `Metadata::SERVICE_CATALOG`. | quickstart, request options, client credentials, error handling. | `php tests/SdkClientTest.php` in CI; local runtime may be absent. |
| Java | Java 11+ client using `java.net.http`. | `ClientCredentialsTokenSource`, `StaticTokenSource`, bearer injection. | `RequestOptions` with headers, idempotency, query, pagination, timeout, retries. | `ServiceException`. | `java/pom.xml`, `SDKMetadata.SDK_NAME`, `SDKMetadata.SDK_VERSION`, `SDKMetadata.serviceCatalog`. | quickstart, request options, client credentials, error handling. | `make -C java test` and `mvn -q -DskipTests package` in CI; local runtime may be absent. |
| C++ | Header-only C++17 client with injectable transport. | `TokenSource` callback and `static_token_source`; no built-in OAuth HTTP fetcher yet. | `RequestOptions` struct with headers, idempotency, query, pagination, timeout, retries. | `twofinance::ServiceError`. | `cpp/CMakeLists.txt`, `twofinance::SDK_NAME`, `twofinance::SDK_VERSION`, `default_service_catalog`. | quickstart, request options, token source, error handling. | `cmake -S . -B build && cmake --build build && ctest --test-dir build`. |

## Domain Coverage

The Go SDK is the reference implementation. Dart, TypeScript, JavaScript,
Python, PHP, Java, and C++ expose the same public domain shape at SDK level:

- `auth`
- `network`
- `analytics`
- `orchestrator`
- `mcp`
- `planner`
- `tradingcontrol`
- `matchengine`
- `hummingbot`
- `keystore`
- `providers`

The canonical operation catalog lives in
`contracts/examples/domain-operations.json` and is checked by
`tools/validate-sdk-structure.mjs`.

## Implementation Notes

- Go is the first fully implemented SDK and keeps the migrated blockchain
  compatibility packages.
- Dart is second priority because it includes the app-facing auth and
  blockchain packages.
- TypeScript and JavaScript share the same service-client shape and examples.
- Python, PHP, Java, and C++ currently prioritize a small, dependency-light
  surface with injectable transports and contract-backed examples.
- Java and PHP are validated by CI configuration, but this local environment
  may not have `javac`, `mvn`, or `php` installed.
- C++ intentionally uses an injectable `TokenSource` rather than a built-in
  OAuth HTTP client for now, keeping the header-only SDK dependency-free.
- Every language exposes SDK name/version constants for diagnostics and
  package reporting.
- Every language exposes the canonical service catalog from
  `contracts/examples/service-catalog.json`.
- Every language exposes public `DomainOperationsCatalog` models for the domain
  operation catalog in `contracts/examples/domain-operations.json`.
- Every language exposes a helper to find an operation in the domain operation
  catalog, such as `findDomainOperation`, `operation(...)`, or
  `find_domain_operation`.
- Every language exposes a helper to resolve an operation into a concrete
  method/path pair from contract-declared path parameters and query parameters,
  such as `resolveDomainOperation`, `operation.resolve(...)`, or
  `resolve_domain_operation`.
- Every language can execute a resolved operation through the shared service
  transport, such as `requestOperation`, `request_operation`, or Go's
  `CallOperation`/`CallResolvedOperation`.
- Every language can resolve and execute directly by catalog domain and
  operation name through helpers such as `resolveCatalogOperation`,
  `requestCatalogOperation`, `resolve_operation`, `request_catalog_operation`,
  or Go's `CallCatalogOperation`.
- Every language exposes a config helper to resolve a service URL by SDK domain
  name, such as `serviceUrl("analytics")` or `service_url("match_engine")`.
- Every language exposes an environment bootstrap path, such as `NewFromEnv`,
  `fromEnv`, `from_env`, `fromEnvironment`, `configFromEnv`,
  `config_from_env`, or `config_from_environment`.
- Every language exposes a config helper to list non-empty configured service
  URLs, such as `serviceUrls()` or `service_urls()`.
- Every language exposes typed configured service entries with `name`, `env`,
  and `url`, such as `ConfiguredServiceEntry`, `configuredServices()`, and
  `configured_services()`.
- Every language exposes a helper for missing service URL configuration, such
  as `missingServiceURLs()` or `missing_service_urls()`.

## Required Example Set

Every language folder must keep examples for:

- quickstart;
- request options and idempotency;
- client credentials or token source;
- error handling.

These examples are part of the SDK compatibility surface and are validated by
the structure checker or language-specific build/test commands.
