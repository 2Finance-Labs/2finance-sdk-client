package com.twofinance.sdk;

public final class StaticTokenSource implements TokenSource {
    private final String token;

    public StaticTokenSource(String token) {
        this.token = token == null ? "" : token;
    }

    @Override
    public String token() {
        return token;
    }
}
