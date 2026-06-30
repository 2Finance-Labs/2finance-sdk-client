# TypeScript SDK

Typed TypeScript SDK for 2Finance services.

```ts
import { TwoFinanceClient, configFromEnv } from "@2finance/sdk-client-typescript";

const client = new TwoFinanceClient(configFromEnv(process.env));
await client.analytics.indicators();
```

Run `npm run typecheck` to validate the package.

See `examples/quickstart.ts` for a minimal analytics and planner flow.
See `examples/request-options.ts` for per-call headers, idempotency, query,
pagination, timeout, and retry options.
See `examples/auth-client-credentials.ts` for client credentials token source
configuration.
See `examples/error-handling.ts` for catching `ServiceError`.
