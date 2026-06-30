package com.twofinance.sdk;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Clock;
import java.time.Instant;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public final class ClientCredentialsTokenSource implements TokenSource {
    private static final Pattern ACCESS_TOKEN_PATTERN = Pattern.compile("\"access_token\"\\s*:\\s*\"([^\"]+)\"");
    private static final Pattern EXPIRES_IN_PATTERN = Pattern.compile("\"expires_in\"\\s*:\\s*(\\d+)");

    private final String tokenUrl;
    private final String clientId;
    private final String clientSecret;
    private final List<String> scopes;
    private final HttpClient httpClient;
    private final Clock clock;
    private final long expirySkewSeconds;
    private String accessToken = "";
    private Instant expiresAt = Instant.EPOCH;

    public ClientCredentialsTokenSource(String tokenUrl, String clientId, String clientSecret, List<String> scopes) {
        this(tokenUrl, clientId, clientSecret, scopes, HttpClient.newHttpClient(), Clock.systemUTC(), 30);
    }

    public ClientCredentialsTokenSource(
            String tokenUrl,
            String clientId,
            String clientSecret,
            List<String> scopes,
            HttpClient httpClient,
            Clock clock,
            long expirySkewSeconds
    ) {
        this.tokenUrl = tokenUrl == null ? "" : tokenUrl;
        this.clientId = clientId == null ? "" : clientId;
        this.clientSecret = clientSecret == null ? "" : clientSecret;
        this.scopes = scopes == null ? List.of() : List.copyOf(scopes);
        this.httpClient = httpClient == null ? HttpClient.newHttpClient() : httpClient;
        this.clock = clock == null ? Clock.systemUTC() : clock;
        this.expirySkewSeconds = expirySkewSeconds;
    }

    @Override
    public synchronized String token() throws IOException, InterruptedException {
        Instant now = clock.instant();
        if (!accessToken.isEmpty() && now.isBefore(expiresAt.minusSeconds(expirySkewSeconds))) {
            return accessToken;
        }
        if (tokenUrl.isBlank() || clientId.isBlank() || clientSecret.isBlank()) {
            throw new IOException("2finance auth: tokenUrl, clientId and clientSecret are required");
        }
        String form = "grant_type=client_credentials"
                + "&client_id=" + encode(clientId)
                + "&client_secret=" + encode(clientSecret)
                + (scopes.isEmpty() ? "" : "&scope=" + encode(String.join(" ", scopes)));
        HttpRequest request = HttpRequest.newBuilder(URI.create(tokenUrl))
                .header("Accept", "application/json")
                .header("Content-Type", "application/x-www-form-urlencoded")
                .POST(HttpRequest.BodyPublishers.ofString(form))
                .build();
        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        if (response.statusCode() < 200 || response.statusCode() >= 300) {
            throw new IOException("2finance auth: token endpoint returned " + response.statusCode());
        }
        String token = find(response.body(), ACCESS_TOKEN_PATTERN);
        if (token.isEmpty()) {
            throw new IOException("2finance auth: token response missing access_token");
        }
        String expiresIn = find(response.body(), EXPIRES_IN_PATTERN);
        long seconds = expiresIn.isEmpty() ? 300 : Long.parseLong(expiresIn);
        accessToken = token;
        expiresAt = now.plusSeconds(seconds);
        return accessToken;
    }

    private static String encode(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8);
    }

    private static String find(String value, Pattern pattern) {
        Matcher matcher = pattern.matcher(value == null ? "" : value);
        return matcher.find() ? matcher.group(1) : "";
    }
}
