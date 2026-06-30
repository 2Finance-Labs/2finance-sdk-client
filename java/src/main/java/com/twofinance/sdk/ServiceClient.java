package com.twofinance.sdk;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;

public class ServiceClient {
    protected final String baseUrl;
    protected final HttpClient httpClient;
    protected final TokenSource tokenSource;

    public ServiceClient(String baseUrl, HttpClient httpClient, TokenSource tokenSource) {
        this.baseUrl = trimTrailingSlash(baseUrl);
        this.httpClient = httpClient == null ? HttpClient.newHttpClient() : httpClient;
        this.tokenSource = tokenSource;
    }

    public String url(String path) {
        if (path.startsWith("http://") || path.startsWith("https://")) {
            return path;
        }
        if (baseUrl.isEmpty()) {
            throw new IllegalStateException("baseUrl is required");
        }
        return baseUrl + "/" + path.replaceFirst("^/+", "");
    }

    public String get(String path) throws IOException, InterruptedException {
        return request("GET", path, null);
    }

    public String get(String path, RequestOptions options) throws IOException, InterruptedException {
        return request("GET", path, null, options);
    }

    public String post(String path, String jsonBody) throws IOException, InterruptedException {
        return request("POST", path, jsonBody);
    }

    public String post(String path, String jsonBody, RequestOptions options) throws IOException, InterruptedException {
        return request("POST", path, jsonBody, options);
    }

    public String put(String path, String jsonBody) throws IOException, InterruptedException {
        return request("PUT", path, jsonBody);
    }

    public String put(String path, String jsonBody, RequestOptions options) throws IOException, InterruptedException {
        return request("PUT", path, jsonBody, options);
    }

    public String delete(String path) throws IOException, InterruptedException {
        return request("DELETE", path, null);
    }

    public String delete(String path, RequestOptions options) throws IOException, InterruptedException {
        return request("DELETE", path, null, options);
    }

    public String request(String method, String path, String jsonBody) throws IOException, InterruptedException {
        return request(method, path, jsonBody, null);
    }

    public String request(String method, String path, String jsonBody, RequestOptions options) throws IOException, InterruptedException {
        String requestUrl = url(path);
        if (options != null && (!options.query().isEmpty() || options.page() != null || options.limit() != null)) {
            requestUrl = urlWithQuery(requestUrl, options);
        }
        HttpRequest.Builder builder = HttpRequest.newBuilder(URI.create(requestUrl))
                .header("Accept", "application/json");
        if (tokenSource != null) {
            try {
                String authorization = Auth.bearerAuthorization(tokenSource.token());
                if (!authorization.isEmpty()) {
                    builder.header("Authorization", authorization);
                }
            } catch (Exception exc) {
                throw new IOException("failed to load bearer token", exc);
            }
        }
        if (options != null) {
            for (var header : options.headers().entrySet()) {
                builder.header(header.getKey(), header.getValue());
            }
            String idempotencyKey = options.idempotencyKey() == null ? "" : options.idempotencyKey().trim();
            if (!idempotencyKey.isEmpty()) {
                builder.header("Idempotency-Key", idempotencyKey);
            }
            if (options.timeout() != null) {
                builder.timeout(options.timeout());
            }
        }
        if (jsonBody == null) {
            builder.method(method, HttpRequest.BodyPublishers.noBody());
        } else {
            builder.header("Content-Type", "application/json");
            builder.method(method, HttpRequest.BodyPublishers.ofString(jsonBody));
        }
        int maxRetries = options == null ? 0 : Math.max(0, options.maxRetries());
        HttpRequest request = builder.build();
        HttpResponse<String> response = null;
        for (int attempt = 0; attempt <= maxRetries; attempt++) {
            response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            if (response.statusCode() >= 200 && response.statusCode() < 300) {
                return response.body();
            }
            if (attempt >= maxRetries || !isRetryableStatus(response.statusCode())) {
                throw new ServiceException(method, requestUrl, response.statusCode(), response.body());
            }
        }
        throw new ServiceException(method, requestUrl, response == null ? 0 : response.statusCode(), response == null ? "" : response.body());
    }

    public String requestOperation(Models.ResolvedOperation operation, String jsonBody, RequestOptions options) throws IOException, InterruptedException {
        return request(operation.method(), operation.path(), jsonBody, options);
    }

    public String requestOperation(Models.ResolvedOperation operation, String jsonBody) throws IOException, InterruptedException {
        return requestOperation(operation, jsonBody, null);
    }

    public String requestCatalogOperation(
            Models.DomainOperationsCatalog catalog,
            String domainName,
            String operationName,
            java.util.Map<String, ?> pathParams,
            java.util.Map<String, ?> query,
            String jsonBody,
            RequestOptions options
    ) throws IOException, InterruptedException {
        return requestOperation(
                catalog.resolveOperation(domainName, operationName, pathParams, query),
                jsonBody,
                options
        );
    }

    private static String trimTrailingSlash(String value) {
        if (value == null) {
            return "";
        }
        return value.trim().replaceAll("/+$", "");
    }

    private static String urlWithQuery(String requestUrl, RequestOptions options) {
        StringBuilder builder = new StringBuilder(requestUrl);
        builder.append(requestUrl.contains("?") ? "&" : "?");
        boolean first = true;
        for (var entry : options.query().entrySet()) {
            if (entry.getValue() == null) {
                continue;
            }
            if (!first) {
                builder.append("&");
            }
            builder.append(URLEncoder.encode(entry.getKey(), StandardCharsets.UTF_8));
            builder.append("=");
            builder.append(URLEncoder.encode(entry.getValue(), StandardCharsets.UTF_8));
            first = false;
        }
        if (options.page() != null) {
            if (!first) {
                builder.append("&");
            }
            builder.append("page=");
            builder.append(URLEncoder.encode(options.page().toString(), StandardCharsets.UTF_8));
            first = false;
        }
        if (options.limit() != null) {
            if (!first) {
                builder.append("&");
            }
            builder.append("limit=");
            builder.append(URLEncoder.encode(options.limit().toString(), StandardCharsets.UTF_8));
            first = false;
        }
        if (first) {
            builder.setLength(builder.length() - 1);
        }
        return builder.toString();
    }

    private static boolean isRetryableStatus(int statusCode) {
        return statusCode == 429 || statusCode >= 500;
    }
}
