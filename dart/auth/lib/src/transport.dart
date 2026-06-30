import 'dart:convert';

import 'package:http/http.dart' as http;

import 'models.dart';

class AuthHttpResponse {
  const AuthHttpResponse({
    required this.statusCode,
    required this.body,
    required this.headers,
    required this.cookies,
  });

  final int statusCode;
  final String body;
  final Map<String, String> headers;
  final List<AuthCookie> cookies;
}

abstract interface class AuthTransport {
  Future<AuthHttpResponse> get(
    Uri url, {
    String? accessToken,
    List<AuthCookie> cookies,
    bool followRedirects,
  });

  Future<AuthHttpResponse> post(
    Uri url, {
    required Map<String, Object?> body,
    String? accessToken,
  });
}

class HttpAuthTransport implements AuthTransport {
  HttpAuthTransport({http.Client? client}) : _client = client ?? http.Client();

  final http.Client _client;

  @override
  Future<AuthHttpResponse> get(
    Uri url, {
    String? accessToken,
    List<AuthCookie> cookies = const [],
    bool followRedirects = true,
  }) async {
    final request = http.Request('GET', url);
    request.followRedirects = followRedirects;
    _addHeaders(request, accessToken: accessToken, cookies: cookies);

    final streamed = await _client.send(request);
    return _readResponse(streamed);
  }

  @override
  Future<AuthHttpResponse> post(
    Uri url, {
    required Map<String, Object?> body,
    String? accessToken,
  }) async {
    final request = http.Request('POST', url);
    request.body = jsonEncode(body);
    _addHeaders(request, accessToken: accessToken);

    final streamed = await _client.send(request);
    return _readResponse(streamed);
  }

  void _addHeaders(
    http.BaseRequest request, {
    String? accessToken,
    List<AuthCookie> cookies = const [],
  }) {
    request.headers['Accept'] = 'application/json';
    request.headers['Content-Type'] = 'application/json';
    if (accessToken != null && accessToken.isNotEmpty) {
      final lower = accessToken.toLowerCase();
      request.headers['Authorization'] = lower.startsWith('bearer ')
          ? accessToken
          : 'Bearer $accessToken';
    }
    if (cookies.isNotEmpty) {
      request.headers['Cookie'] = cookies
          .map((cookie) => '${cookie.name}=${cookie.value}')
          .join('; ');
    }
  }

  Future<AuthHttpResponse> _readResponse(http.StreamedResponse response) async {
    final body = await response.stream.bytesToString();
    return AuthHttpResponse(
      statusCode: response.statusCode,
      body: body,
      headers: response.headers,
      cookies: _parseSetCookieHeaders(response.headers['set-cookie']),
    );
  }
}

void ensureSuccess(AuthHttpResponse response, String message) {
  if (response.statusCode < 200 || response.statusCode >= 300) {
    throw AuthSdkException(
      message,
      statusCode: response.statusCode,
      body: response.body,
    );
  }
}

Map<String, Object?> decodeObject(String body) {
  final decoded = jsonDecode(body);
  if (decoded is! Map<String, Object?>) {
    throw AuthSdkException('Expected a JSON object response');
  }
  return decoded;
}

List<AuthCookie> _parseSetCookieHeaders(String? header) {
  if (header == null || header.isEmpty) {
    return const [];
  }

  final cookies = <AuthCookie>[];
  for (final value in _splitSetCookieHeader(header)) {
    try {
      cookies.add(AuthCookie.fromSetCookieValue(value));
    } on FormatException {
      // Ignore malformed cookie fragments returned by intermediaries.
    }
  }
  return cookies;
}

List<String> _splitSetCookieHeader(String header) {
  final parts = <String>[];
  final buffer = StringBuffer();
  var inExpires = false;

  for (var i = 0; i < header.length; i++) {
    final char = header[i];
    if (char == ',') {
      if (inExpires) {
        buffer.write(char);
        continue;
      }
      parts.add(buffer.toString().trim());
      buffer.clear();
      continue;
    }

    buffer.write(char);
    final current = buffer.toString().toLowerCase();
    if (current.endsWith('expires=')) {
      inExpires = true;
    } else if (inExpires && char == ';') {
      inExpires = false;
    }
  }

  final last = buffer.toString().trim();
  if (last.isNotEmpty) {
    parts.add(last);
  }
  return parts;
}
