# JavaScript SDK

Dependency-free JavaScript SDK for Node 18+.

```js
const { TwoFinanceClient, configFromEnv } = require("@2finance/sdk-client");

const client = new TwoFinanceClient(configFromEnv(process.env));
await client.analytics.indicators();
```

Run `npm test` to validate the package.

See `examples/quickstart.js` for a minimal analytics and planner flow.
See `examples/request-options.js` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `examples/auth-client-credentials.js` for client credentials token source
configuration.
See `examples/error-handling.js` for catching `ServiceError`.
