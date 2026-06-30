package com.twofinance.sdk;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;

public final class Models {
    private Models() {}

    public record SdkErrorPayload(
            String error,
            String message,
            String code,
            Map<String, Object> details
    ) {}

    public record PaginationResponse(
            List<Map<String, Object>> items,
            int limit,
            String cursor,
            String nextCursor
    ) {}

    public record IdempotencyRecord(
            String idempotencyKey,
            String operation,
            String scope,
            String requestId
    ) {}

    public record ServiceCatalogEntry(String name, String env) {}

    public record ServiceCatalog(List<ServiceCatalogEntry> services) {}

    public record ConfiguredServiceEntry(String name, String env, String url) {}

    public record DomainOperation(
            String name,
            String method,
            String path,
            List<String> pathParams,
            List<String> query,
            String requestSchema,
            String responseSchema,
            String notes
    ) {
        public ResolvedOperation resolve(Map<String, ?> pathParams, Map<String, ?> query) {
            Map<String, ?> pathValues = pathParams == null ? Map.of() : pathParams;
            Map<String, ?> queryValues = query == null ? Map.of() : query;
            String resolvedPath = path();

            for (String name : pathParams() == null ? List.<String>of() : pathParams()) {
                if (!pathValues.containsKey(name)) {
                    throw new IllegalArgumentException("2finance: missing operation path parameter " + name);
                }
                resolvedPath = resolvedPath.replace("{" + name + "}", encodeComponent(String.valueOf(pathValues.get(name))));
            }

            StringBuilder queryString = new StringBuilder();
            for (String name : query() == null ? List.<String>of() : query()) {
                Object value = queryValues.get(name);
                if (value == null) {
                    continue;
                }
                if (queryString.length() > 0) {
                    queryString.append('&');
                }
                queryString.append(encodeComponent(name));
                queryString.append('=');
                queryString.append(encodeComponent(String.valueOf(value)));
            }
            if (queryString.length() > 0) {
                resolvedPath += resolvedPath.contains("?") ? "&" : "?";
                resolvedPath += queryString;
            }

            String resolvedMethod = method() == null ? "" : method().trim().toUpperCase(Locale.ROOT);
            return new ResolvedOperation(resolvedMethod, resolvedPath);
        }
    }

    public record ResolvedOperation(String method, String path) {}

    public record DomainOperationsDomain(
            String name,
            String env,
            String transport,
            String description,
            List<DomainOperation> operations
    ) {}

    public record DomainOperationsCatalog(
            String schema,
            List<DomainOperationsDomain> domains
    ) {
        public Optional<DomainOperation> operation(String domainName, String operationName) {
            for (DomainOperationsDomain domain : domains) {
                if (!domainKey(domain.name()).equals(domainKey(domainName))) {
                    continue;
                }
                return domain.operations().stream()
                        .filter(operation -> operation.name().equals(operationName))
                        .findFirst();
            }
            return Optional.empty();
        }

        public ResolvedOperation resolveOperation(String domainName, String operationName, Map<String, ?> pathParams, Map<String, ?> query) {
            return operation(domainName, operationName)
                    .orElseThrow(() -> new IllegalArgumentException("2finance: unknown operation " + domainName + "." + operationName))
                    .resolve(pathParams, query);
        }

        private static String domainKey(String value) {
            return value == null ? "" : value.trim().toLowerCase().replace("-", "").replace("_", "").replace(" ", "");
        }
    }

    private static String encodeComponent(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8).replace("+", "%20");
    }
}
