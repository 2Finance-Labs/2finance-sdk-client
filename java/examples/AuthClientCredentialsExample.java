package com.twofinance.sdk.examples;

import com.twofinance.sdk.ClientCredentialsTokenSource;
import com.twofinance.sdk.SDKConfig;
import com.twofinance.sdk.SdkClient;
import java.util.List;

public final class AuthClientCredentialsExample {
    private AuthClientCredentialsExample() {}

    public static void main(String[] args) throws Exception {
        SDKConfig config = SDKConfig.fromEnv();
        config.tokenSource = new ClientCredentialsTokenSource(
                System.getenv().getOrDefault("TWO_FINANCE_AUTH_TOKEN_URL", ""),
                System.getenv().getOrDefault("TWO_FINANCE_AUTH_CLIENT_ID", ""),
                System.getenv().getOrDefault("TWO_FINANCE_AUTH_CLIENT_SECRET", ""),
                List.of("2finance.sdk")
        );

        SdkClient client = new SdkClient(config);
        System.out.println(client.analytics.indicators());
    }
}
