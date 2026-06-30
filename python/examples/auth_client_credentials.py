import os

from twofinance_sdk_client import ClientCredentialsTokenSource, TwoFinanceClient, config_from_env


def main() -> None:
    token_source = ClientCredentialsTokenSource(
        token_url=os.environ.get("TWO_FINANCE_AUTH_TOKEN_URL", ""),
        client_id=os.environ.get("TWO_FINANCE_AUTH_CLIENT_ID", ""),
        client_secret=os.environ.get("TWO_FINANCE_AUTH_CLIENT_SECRET", ""),
        scopes=["2finance.sdk"],
    )
    client = TwoFinanceClient(config_from_env(), token_source=token_source)
    print(client.analytics.indicators())


if __name__ == "__main__":
    main()
