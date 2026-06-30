import { ClientCredentialsTokenSource, TwoFinanceClient, configFromEnv } from "../src/index";

async function main() {
  const tokenSource = new ClientCredentialsTokenSource({
    tokenURL: process.env.TWO_FINANCE_AUTH_TOKEN_URL || "",
    clientID: process.env.TWO_FINANCE_AUTH_CLIENT_ID || "",
    clientSecret: process.env.TWO_FINANCE_AUTH_CLIENT_SECRET || "",
    scopes: ["2finance.sdk"]
  });

  const client = new TwoFinanceClient({
    ...configFromEnv(process.env),
    tokenSource
  });
  console.log(await client.analytics.indicators());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
