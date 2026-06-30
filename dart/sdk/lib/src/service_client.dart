import 'dart:convert';

import 'auth.dart';
import 'models.dart';

typedef Transport = Future<ServiceResponse> Function(ServiceRequest request);

final class ServiceRequest {
  const ServiceRequest({
    required this.method,
    required this.url,
    required this.headers,
    this.body,
  });

  final String method;
  final Uri url;
  final Map<String, String> headers;
  final String? body;
}

final class ServiceResponse {
  const ServiceResponse({required this.statusCode, this.body = ''});

  final int statusCode;
  final String body;
}

final class RequestOptions {
  const RequestOptions({
    this.headers = const {},
    this.idempotencyKey,
    this.queryParameters = const {},
    this.timeout,
    this.maxRetries = 0,
    this.page,
    this.limit,
  });

  final Map<String, String> headers;
  final String? idempotencyKey;
  final Map<String, String> queryParameters;
  final Duration? timeout;
  final int maxRetries;
  final int? page;
  final int? limit;
}

final class ServiceException implements Exception {
  ServiceException(this.method, this.url, this.statusCode, this.body);

  final String method;
  final Uri url;
  final int statusCode;
  final String body;

  @override
  String toString() {
    return '2finance: $method $url returned $statusCode: $body';
  }
}

final class ServiceClient {
  ServiceClient(this.baseUrl, {this.transport, this.tokenSource});

  final String baseUrl;
  final Transport? transport;
  final TokenSource? tokenSource;

  Uri url(String path) {
    if (path.startsWith('http://') || path.startsWith('https://')) {
      return Uri.parse(path);
    }
    final base = baseUrl.trim().replaceFirst(RegExp(r'/+$'), '');
    if (base.isEmpty) {
      throw ArgumentError('baseUrl is required');
    }
    return Uri.parse('$base/${path.replaceFirst(RegExp(r'^/+'), '')}');
  }

  Future<Object?> get(String path, [RequestOptions? options]) =>
      request('GET', path, null, options);

  Future<Object?> post(String path, [Object? body, RequestOptions? options]) =>
      request('POST', path, body, options);

  Future<Object?> put(String path, [Object? body, RequestOptions? options]) =>
      request('PUT', path, body, options);

  Future<Object?> delete(String path, [RequestOptions? options]) =>
      request('DELETE', path, null, options);

  Future<Object?> request(
    String method,
    String path, [
    Object? body,
    RequestOptions? options,
  ]) async {
    final headers = <String, String>{'Accept': 'application/json'};
    String? payload;
    if (body != null) {
      headers['Content-Type'] = 'application/json';
      payload = body is String ? body : jsonEncode(body);
    }
    final source = tokenSource;
    if (source != null) {
      final authorization = bearerAuthorization(await source.token());
      if (authorization.isNotEmpty) {
        headers['Authorization'] = authorization;
      }
    }
    final requestOptions = options;
    if (requestOptions != null) {
      headers.addAll(requestOptions.headers);
      final idempotencyKey = requestOptions.idempotencyKey?.trim();
      if (idempotencyKey != null && idempotencyKey.isNotEmpty) {
        headers['Idempotency-Key'] = idempotencyKey;
      }
    }
    final clientTransport = transport;
    if (clientTransport == null) {
      throw StateError('transport is required');
    }
    var requestUrl = url(path);
    if (requestOptions != null && requestOptions.queryParameters.isNotEmpty) {
      requestUrl = requestUrl.replace(
        queryParameters: {
          ...requestUrl.queryParameters,
          ...requestOptions.queryParameters,
          if (requestOptions.page != null) 'page': '${requestOptions.page}',
          if (requestOptions.limit != null) 'limit': '${requestOptions.limit}',
        },
      );
    } else if (requestOptions != null &&
        (requestOptions.page != null || requestOptions.limit != null)) {
      requestUrl = requestUrl.replace(
        queryParameters: {
          ...requestUrl.queryParameters,
          if (requestOptions.page != null) 'page': '${requestOptions.page}',
          if (requestOptions.limit != null) 'limit': '${requestOptions.limit}',
        },
      );
    }
    final timeout = requestOptions?.timeout;
    final maxRetries = requestOptions?.maxRetries ?? 0;
    final attempts = 1 + (maxRetries < 0 ? 0 : maxRetries);
    ServiceResponse? response;
    for (var attempt = 0; attempt < attempts; attempt++) {
      final pendingResponse = clientTransport(
        ServiceRequest(
          method: method,
          url: requestUrl,
          headers: headers,
          body: payload,
        ),
      );
      response = timeout == null
          ? await pendingResponse
          : await pendingResponse.timeout(timeout);
      if (!_isRetryableStatus(response.statusCode) || attempt + 1 >= attempts) {
        break;
      }
    }
    response!;
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw ServiceException(
        method,
        requestUrl,
        response.statusCode,
        response.body,
      );
    }
    if (response.body.isEmpty) {
      return null;
    }
    return jsonDecode(response.body) as Object?;
  }

  Future<Object?> requestOperation(
    ResolvedOperation operation, [
    Object? body,
    RequestOptions? options,
  ]) {
    return request(operation.method, operation.path, body, options);
  }

  Future<Object?> requestCatalogOperation(
    DomainOperationsCatalog catalog,
    String domainName,
    String operationName, {
    Map<String, Object?> pathParams = const {},
    Map<String, Object?> queryParameters = const {},
    Object? body,
    RequestOptions? options,
  }) {
    final operation = catalog.resolveOperation(
      domainName,
      operationName,
      pathParams: pathParams,
      queryParameters: queryParameters,
    );
    return requestOperation(operation, body, options);
  }

  static bool _isRetryableStatus(int statusCode) =>
      statusCode == 429 || statusCode >= 500;
}
