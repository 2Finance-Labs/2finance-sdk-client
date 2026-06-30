package com.twofinance.sdk.examples;

import com.twofinance.sdk.SdkClient;
import com.twofinance.sdk.ServiceException;

public final class ErrorHandlingExample {
    private ErrorHandlingExample() {}

    public static void main(String[] args) throws Exception {
        SdkClient client = SdkClient.fromEnv();

        try {
            client.analytics.indicators();
        } catch (ServiceException exception) {
            System.out.println(
                    "request failed with status " + exception.statusCode() + ": " + exception.body()
            );
        }
    }
}
