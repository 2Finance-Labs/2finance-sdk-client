import { RequestOptions, TwoFinanceClient, configFromEnv } from "../src/index";

async function main() {
  const client = new TwoFinanceClient(configFromEnv(process.env));
  const options: RequestOptions = {
    headers: { "X-Trace-ID": "trace-1" },
    idempotencyKey: "candles-upsert-001",
    query: { source: "sdk-example" },
    pagination: { page: 1, limit: 25 },
    timeoutMs: 5000,
    maxRetries: 1
  };

  const response = await client.analytics.post(
    "/analytics/candles:upsert",
    { symbol: "BTC-USDT" },
    options
  );
  console.log("response:", response);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
