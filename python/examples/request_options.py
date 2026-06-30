from twofinance_sdk_client import RequestOptions, TwoFinanceClient, config_from_env


def main() -> None:
    client = TwoFinanceClient(config_from_env())
    response = client.analytics.post(
        "/analytics/candles:upsert",
        {"symbol": "BTC-USDT"},
        RequestOptions(
            headers={"X-Trace-ID": "trace-1"},
            idempotency_key="candles-upsert-001",
            query={"source": "sdk-example"},
            page=1,
            limit=25,
            timeout=5.0,
            max_retries=1,
        ),
    )
    print("response:", response)


if __name__ == "__main__":
    main()
