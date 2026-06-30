import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:two_finance_blockchain/infra/event/request_response.dart';
import 'package:two_finance_blockchain/infra/transport/transport.dart';

typedef TokenProvider = Future<String?> Function();

class HttpFinanceNetworkTransport implements FinanceNetworkTransport {
  HttpFinanceNetworkTransport({
    String? baseUrl,
    http.Client? httpClient,
    TokenProvider? tokenProvider,
  }) : baseUrl = (baseUrl ?? 'http://127.0.0.1:9095').replaceFirst(
         RegExp(r'/$'),
         '',
       ),
       _httpClient = httpClient ?? http.Client(),
       _tokenProvider = tokenProvider;

  final String baseUrl;
  final http.Client _httpClient;
  final TokenProvider? _tokenProvider;

  @override
  Future<dynamic> sendRequest(
    String method,
    dynamic params,
    String replyTo,
  ) async {
    final headers = <String, String>{'Content-Type': 'application/json'};
    final accessToken = await _tokenProvider?.call();
    if (accessToken != null && accessToken.trim().isNotEmpty) {
      headers['Authorization'] = _bearer(accessToken);
    }

    final response = await _httpClient.post(
      Uri.parse('$baseUrl/v1/2finance-network/virtual-machine'),
      headers: headers,
      body: jsonEncode(RequestPayload(method: method, params: params).toJson()),
    );

    final responseBody = response.body;
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw Exception(
        '2finance-network HTTP ${response.statusCode}: $responseBody',
      );
    }

    final decoded = jsonDecode(responseBody);
    if (decoded is! Map<String, dynamic>) {
      throw const FormatException(
        '2finance-network response is not an object.',
      );
    }

    final code = decoded['code'] ?? decoded['Code'];
    if (code is num && code != 200) {
      final msg = decoded['msg'] ?? decoded['Msg'] ?? 'request failed';
      final data = decoded['data'] ?? decoded['Data'];
      throw Exception('2finance-network error $code: $msg $data');
    }

    return decoded['data'] ?? decoded['Data'];
  }

  String _bearer(String accessToken) {
    final trimmed = accessToken.trim();
    if (trimmed.toLowerCase().startsWith('bearer ')) {
      return trimmed;
    }
    return 'Bearer $trimmed';
  }
}
