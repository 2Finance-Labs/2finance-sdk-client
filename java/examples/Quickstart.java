package com.twofinance.sdk.examples;

import com.twofinance.sdk.SdkClient;

public final class Quickstart {
    private Quickstart() {}

    public static void main(String[] args) throws Exception {
        SdkClient client = SdkClient.fromEnv();

        String indicators = client.analytics.indicators();
        System.out.println("analytics indicators: " + indicators);

        String plan = client.planner.tradingPlan(
                "{\"goal\":\"prepare a BTC rebalancing plan\",\"useAnalytics\":true,\"useTrading\":true}",
                true,
                true
        );
        System.out.println("planner response: " + plan);
    }
}
