import { TwoFinanceClient, configFromEnv } from "../src/index";

async function main() {
  const client = new TwoFinanceClient(configFromEnv(process.env));

  const indicators = await client.analytics.indicators();
  console.log("analytics indicators:", indicators);

  const plan = await client.planner.tradingPlan({
    goal: "prepare a BTC rebalancing plan",
    useAnalytics: true,
    useTrading: true
  });
  console.log("planner response:", plan);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
