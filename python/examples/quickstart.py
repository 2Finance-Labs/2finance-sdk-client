from twofinance_sdk_client import TwoFinanceClient, config_from_env


def main() -> None:
    client = TwoFinanceClient(config_from_env())

    indicators = client.analytics.indicators()
    print("analytics indicators:", indicators)

    plan = client.planner.trading_plan(
        {
            "goal": "prepare a BTC rebalancing plan",
            "useAnalytics": True,
            "useTrading": True,
        }
    )
    print("planner response:", plan)


if __name__ == "__main__":
    main()
