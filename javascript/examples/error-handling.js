"use strict";

const { ServiceError, TwoFinanceClient, configFromEnv } = require("../src");

async function main() {
  const client = new TwoFinanceClient(configFromEnv(process.env));

  try {
    await client.analytics.indicators();
  } catch (error) {
    if (error instanceof ServiceError) {
      console.log(`request failed with status ${error.statusCode}: ${error.body}`);
      return;
    }
    throw error;
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
