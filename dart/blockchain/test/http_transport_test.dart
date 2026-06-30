import 'dart:convert';

import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:test/test.dart';
import 'package:two_finance_blockchain/infra/http/http_transport.dart';

void main() {
  test(
    'HttpFinanceNetworkTransport posts to virtual-machine endpoint',
    () async {
      late Uri receivedUri;
      late Map<String, String> receivedHeaders;
      late Map<String, dynamic> receivedBody;

      final transport = HttpFinanceNetworkTransport(
        baseUrl: 'http://2finance-network:9095/',
        tokenProvider: () async => 'test-access-token',
        httpClient: MockClient((request) async {
          receivedUri = request.url;
          receivedHeaders = request.headers;
          receivedBody = jsonDecode(request.body) as Map<String, dynamic>;
          return http.Response(
            jsonEncode({
              'code': 200,
              'msg': 'Successfully',
              'data': {'states': <dynamic>[], 'logs': <dynamic>[]},
            }),
            200,
            headers: {'content-type': 'application/json'},
          );
        }),
      );

      final data = await transport.sendRequest('get_state', {
        'to': 'wallet',
        'method': 'get',
      }, 'reply-id');

      expect(
        receivedUri.toString(),
        'http://2finance-network:9095/v1/2finance-network/virtual-machine',
      );
      expect(receivedBody['method'], 'get_state');
      expect(receivedBody['params'], {'to': 'wallet', 'method': 'get'});
      expect(receivedHeaders['Authorization'], 'Bearer test-access-token');
      expect(data, {'states': <dynamic>[], 'logs': <dynamic>[]});
    },
  );

  test('HttpFinanceNetworkTransport works without a token provider', () async {
    late Map<String, String> receivedHeaders;
    final transport = HttpFinanceNetworkTransport(
      baseUrl: 'http://2finance-network:9095/',
      httpClient: MockClient((request) async {
        receivedHeaders = request.headers;
        return http.Response(
          jsonEncode({
            'code': 200,
            'msg': 'Successfully',
            'data': {'ok': true},
          }),
          200,
        );
      }),
    );

    final data = await transport.sendRequest(
      'get_state',
      <String, dynamic>{},
      'reply-id',
    );

    expect(receivedHeaders.containsKey('Authorization'), isFalse);
    expect(data, {'ok': true});
  });

  test('HttpFinanceNetworkTransport ignores empty bearer tokens', () async {
    late Map<String, String> receivedHeaders;
    final transport = HttpFinanceNetworkTransport(
      baseUrl: 'http://2finance-network:9095/',
      tokenProvider: () async => '   ',
      httpClient: MockClient((request) async {
        receivedHeaders = request.headers;
        return http.Response(
          jsonEncode({
            'code': 200,
            'msg': 'Successfully',
            'data': {'ok': true},
          }),
          200,
        );
      }),
    );

    await transport.sendRequest('get_state', <String, dynamic>{}, 'reply-id');

    expect(receivedHeaders.containsKey('Authorization'), isFalse);
  });

  test(
    'HttpFinanceNetworkTransport keeps preformatted bearer tokens',
    () async {
      late Map<String, String> receivedHeaders;
      final transport = HttpFinanceNetworkTransport(
        baseUrl: 'http://2finance-network:9095/',
        tokenProvider: () async => 'Bearer already-prefixed',
        httpClient: MockClient((request) async {
          receivedHeaders = request.headers;
          return http.Response(
            jsonEncode({
              'code': 200,
              'msg': 'Successfully',
              'data': {'ok': true},
            }),
            200,
          );
        }),
      );

      await transport.sendRequest('get_state', <String, dynamic>{}, 'reply-id');

      expect(receivedHeaders['Authorization'], 'Bearer already-prefixed');
    },
  );

  test(
    'HttpFinanceNetworkTransport does not include bearer token in errors',
    () {
      final transport = HttpFinanceNetworkTransport(
        baseUrl: 'http://2finance-network:9095/',
        tokenProvider: () async => 'sensitive-token-value',
        httpClient: MockClient((request) async {
          return http.Response('upstream denied', 503);
        }),
      );

      expect(
        () =>
            transport.sendRequest('get_state', <String, dynamic>{}, 'reply-id'),
        throwsA(
          isA<Exception>().having(
            (error) => error.toString(),
            'message',
            allOf(
              contains('2finance-network HTTP 503'),
              isNot(contains('sensitive-token-value')),
            ),
          ),
        ),
      );
    },
  );
}
