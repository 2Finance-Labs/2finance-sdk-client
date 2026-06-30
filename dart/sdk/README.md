# Dart Unified SDK

Dependency-light Dart SDK facade for 2Finance services.

```dart
final client = TwoFinanceClient.fromEnvironment();
await client.analytics.indicators();
```

Run `dart test`.

See `example/quickstart.dart` for a minimal analytics and planner flow.
See `example/request_options.dart` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `example/auth_client_credentials.dart` for client credentials token source
configuration.
See `example/error_handling.dart` for catching `ServiceException`.
