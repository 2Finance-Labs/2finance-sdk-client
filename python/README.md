# Python SDK

Dependency-free Python SDK for 2Finance services.

```python
from twofinance_sdk_client import TwoFinanceClient, config_from_env

client = TwoFinanceClient(config_from_env())
client.analytics.indicators()
```

Run `PYTHONPATH=. python3 -m unittest discover -s tests -v` to validate the package.

See `examples/quickstart.py` for a minimal analytics and planner flow.
See `examples/request_options.py` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `examples/auth_client_credentials.py` for client credentials token source
configuration.
See `examples/error_handling.py` for catching `HTTPError`.
