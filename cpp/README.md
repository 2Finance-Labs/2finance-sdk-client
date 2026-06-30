# C++ SDK

Header-only C++17 SDK surface for 2Finance services.

The C++ client accepts an injectable HTTP transport so callers can use libcurl,
Boost.Beast, platform HTTP, or a test fake without this package choosing a
network dependency.

```cpp
twofinance::SdkConfig config = twofinance::config_from_environment();
twofinance::SdkClient client(config, transport);
```

Build the smoke test with CMake:

```bash
cmake -S . -B build
cmake --build build
ctest --test-dir build
```

Consumers can link the header-only target as `twofinance::sdk_client`.

See `examples/quickstart.cpp` for a minimal analytics and planner flow with an
injectable transport.
See `examples/request_options.cpp` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `examples/auth_token_source.cpp` for injectable bearer token source usage.
See `examples/error_handling.cpp` for catching `twofinance::ServiceError`.
