package com.twofinance.sdk;

import java.time.Duration;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;

public final class RequestOptions {
    private final Map<String, String> headers;
    private final String idempotencyKey;
    private final Map<String, String> query;
    private final Duration timeout;
    private final int maxRetries;
    private final Integer page;
    private final Integer limit;

    public RequestOptions() {
        this(Collections.emptyMap(), null, Collections.emptyMap());
    }

    public RequestOptions(Map<String, String> headers, String idempotencyKey) {
        this(headers, idempotencyKey, Collections.emptyMap());
    }

    public RequestOptions(Map<String, String> headers, String idempotencyKey, Map<String, String> query) {
        this(headers, idempotencyKey, query, null);
    }

    public RequestOptions(Map<String, String> headers, String idempotencyKey, Map<String, String> query, Duration timeout) {
        this(headers, idempotencyKey, query, timeout, 0);
    }

    public RequestOptions(Map<String, String> headers, String idempotencyKey, Map<String, String> query, Duration timeout, int maxRetries) {
        this(headers, idempotencyKey, query, timeout, maxRetries, null, null);
    }

    public RequestOptions(
            Map<String, String> headers,
            String idempotencyKey,
            Map<String, String> query,
            Duration timeout,
            int maxRetries,
            Integer page,
            Integer limit
    ) {
        this.headers = headers == null ? Collections.emptyMap() : Collections.unmodifiableMap(new LinkedHashMap<>(headers));
        this.idempotencyKey = idempotencyKey;
        this.query = query == null ? Collections.emptyMap() : Collections.unmodifiableMap(new LinkedHashMap<>(query));
        this.timeout = timeout;
        this.maxRetries = maxRetries;
        this.page = page;
        this.limit = limit;
    }

    public Map<String, String> headers() {
        return headers;
    }

    public String idempotencyKey() {
        return idempotencyKey;
    }

    public Map<String, String> query() {
        return query;
    }

    public Duration timeout() {
        return timeout;
    }

    public int maxRetries() {
        return maxRetries;
    }

    public Integer page() {
        return page;
    }

    public Integer limit() {
        return limit;
    }
}
