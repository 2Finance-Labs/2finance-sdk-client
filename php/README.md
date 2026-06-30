# PHP SDK

Dependency-free PHP SDK for 2Finance services.

```php
use TwoFinance\Sdk\SdkClient;
use TwoFinance\Sdk\SdkConfig;

$client = new SdkClient(SdkConfig::fromEnv());
$client->analytics->indicators();
```

Run `php tests/SdkClientTest.php` to validate the package.

See `examples/quickstart.php` for a minimal analytics and planner flow.
See `examples/request_options.php` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `examples/auth_client_credentials.php` for client credentials token source
configuration.
See `examples/error_handling.php` for catching `ServiceException`.
