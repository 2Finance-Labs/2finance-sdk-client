package com.twofinance.sdk;

import java.io.IOException;

public final class ServiceException extends IOException {
    private final String method;
    private final String url;
    private final int statusCode;
    private final String body;

    public ServiceException(String method, String url, int statusCode, String body) {
        super("2finance: " + method + " " + url + " returned " + statusCode + ": " + body);
        this.method = method;
        this.url = url;
        this.statusCode = statusCode;
        this.body = body;
    }

    public String method() {
        return method;
    }

    public String url() {
        return url;
    }

    public int statusCode() {
        return statusCode;
    }

    public String body() {
        return body;
    }
}
