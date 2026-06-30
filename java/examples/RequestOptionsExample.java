package com.twofinance.sdk.examples;

import com.twofinance.sdk.RequestOptions;
import com.twofinance.sdk.SdkClient;
import java.time.Duration;
import java.util.Map;

public final class RequestOptionsExample {
    private RequestOptionsExample() {}

    public static void main(String[] args) throws Exception {
        SdkClient client = SdkClient.fromEnv();
        RequestOptions options = new RequestOptions(
                Map.of("X-Trace-ID", "trace-1"),
                "candles-upsert-001",
                Map.of("source", "sdk-example"),
                Duration.ofSeconds(5),
                1,
                1,
                25
        );
        String response = client.analytics.post(
                "/analytics/candles:upsert",
                "{\"symbol\":\"BTC-USDT\"}",
                options
        );
        System.out.println("response: " + response);
    }
}
