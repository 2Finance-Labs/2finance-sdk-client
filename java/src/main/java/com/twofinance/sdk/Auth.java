package com.twofinance.sdk;

public final class Auth {
    private Auth() {}

    public static String bearerAuthorization(String token) {
        String trimmed = token == null ? "" : token.trim();
        if (trimmed.isEmpty()) {
            return "";
        }
        if (trimmed.toLowerCase().startsWith("bearer ")) {
            return trimmed;
        }
        return "Bearer " + trimmed;
    }
}
