from twofinance_sdk_client import HTTPError, TwoFinanceClient, config_from_env


def main() -> None:
    client = TwoFinanceClient(config_from_env())

    try:
        client.analytics.indicators()
    except HTTPError as error:
        print(f"request failed with status {error.status_code}: {error.body}")


if __name__ == "__main__":
    main()
