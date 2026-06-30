import 'dart:convert';

import 'package:test/test.dart';
import 'package:two_finance_auth_sdk/two_finance_auth_sdk.dart';

void main() {
  group('TwoFinanceAuthClient', () {
    test('login posts to the current password login route', () async {
      final transport = FakeTransport();
      transport.nextPost = AuthHttpResponse(
        statusCode: 200,
        body: jsonEncode({
          'code': 200,
          'msg': 'Successfully',
          'data': _jwtJson(),
        }),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      // ignore: deprecated_member_use_from_same_package
      final jwt = await client.login(
        // ignore: deprecated_member_use_from_same_package
        const LoginInput(username: 'luiz', password: 'secret'),
      );

      expect(jwt.accessToken, 'access');
      expect(
        transport.lastUrl.toString(),
        'http://localhost:8080/v1/2finance-authenticator/realm/client/login',
      );
      expect(transport.lastBody, {'username': 'luiz', 'password': 'secret'});
    });

    test('protected routes send Authorization as a bearer token', () async {
      final transport = FakeTransport();
      transport.nextGet = AuthHttpResponse(
        statusCode: 200,
        body: jsonEncode({
          'code': 200,
          'msg': 'Successfully',
          'data': {'sub': 'user-id'},
        }),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      final response = await client.getUserInfo('access-token');

      expect(response.code, 200);
      expect(transport.lastAccessToken, 'access-token');
      expect(
        transport.lastUrl.toString(),
        'http://localhost:8080/v1/2finance-authenticator/realm/client/user-info',
      );
    });

    test('phone login uses the phone client id route', () async {
      final transport = FakeTransport();
      transport.nextPost = AuthHttpResponse(
        statusCode: 200,
        body: jsonEncode({
          'code': 200,
          'msg': 'Token exchanged successfully',
          'data': _jwtJson(),
        }),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      await client.phoneLogin(phoneNumber: '+5511999999999', code: '123456');

      expect(
        transport.lastUrl.toString(),
        'http://localhost:8080/v1/2finance-authenticator/realm/phone-client/phone/sms/login',
      );
    });

    test('logout posts refresh token to the auth client route', () async {
      final transport = FakeTransport();
      transport.nextPost = AuthHttpResponse(
        statusCode: 200,
        body: jsonEncode({'code': 200, 'msg': 'Successfully', 'data': null}),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      final response = await client.logout('refresh-token');

      expect(response.code, 200);
      expect(
        transport.lastUrl.toString(),
        'http://localhost:8080/v1/2finance-authenticator/realm/client/logout',
      );
      expect(transport.lastBody, {'refresh_token': 'refresh-token'});
    });

    test('pkce login exposes redirect location, state and cookies', () async {
      final transport = FakeTransport();
      transport.nextGet = AuthHttpResponse(
        statusCode: 302,
        body: '',
        headers: const {'location': 'https://idp.example/auth?state=state-123'},
        cookies: const [AuthCookie(name: 'oauth_state', value: 'state-123')],
      );
      final client = _client(transport);

      final response = await client.loginPKCE();

      expect(response.authUrl.toString(), contains('state=state-123'));
      expect(response.state, 'state-123');
      expect(response.cookies.single.name, 'oauth_state');
      expect(transport.lastFollowRedirects, isFalse);
    });

    test('pkce callback sends state and auth-flow cookies', () async {
      final transport = FakeTransport();
      transport.nextGet = AuthHttpResponse(
        statusCode: 200,
        body: jsonEncode({
          'access_token': 'access',
          'refresh_token': 'refresh',
          'id_token': 'id',
          'expires_in': 300,
          'token_type': 'Bearer',
        }),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      final response = await client.callbackPKCE(
        code: 'code-123',
        state: 'state-123',
        cookies: const [AuthCookie(name: 'oauth_state', value: 'state-123')],
      );

      expect(response.idToken, 'id');
      expect(transport.lastUrl.queryParameters, {
        'code': 'code-123',
        'state': 'state-123',
      });
      expect(transport.lastCookies.single.name, 'oauth_state');
    });

    test('failed responses redact credentials and tokens in exceptions', () {
      final body = jsonEncode({
        'access_token': 'access-secret',
        'refresh_token': 'refresh-secret',
        'id_token': 'id-secret',
        'password': 'password-secret',
        'error': 'Authorization: Bearer bearer-secret',
      });
      final response = AuthHttpResponse(
        statusCode: 401,
        body: body,
        headers: const {},
        cookies: const [],
      );

      expect(
        () =>
            ensureSuccess(response, 'Login failed with Bearer message-secret'),
        throwsA(
          isA<AuthSdkException>()
              .having((error) => error.body, 'body', isNot(contains('secret')))
              .having(
                (error) => error.body,
                'redacted body',
                allOf([
                  contains('"access_token":"<redacted>"'),
                  contains('"refresh_token":"<redacted>"'),
                  contains('"id_token":"<redacted>"'),
                  contains('"password":"<redacted>"'),
                  contains('Bearer <redacted>'),
                ]),
              )
              .having(
                (error) => error.toString(),
                'toString',
                isNot(contains('secret')),
              ),
        ),
      );
    });

    test('logout errors redact refresh token in exceptions', () async {
      final transport = FakeTransport();
      transport.nextPost = AuthHttpResponse(
        statusCode: 400,
        body: jsonEncode({
          'error': 'client_secret=secret refresh_token=refresh-secret',
        }),
        headers: const {},
        cookies: const [],
      );
      final client = _client(transport);

      await expectLater(
        client.logout('refresh-secret'),
        throwsA(
          isA<AuthSdkException>()
              .having(
                (error) => error.toString(),
                'toString',
                isNot(contains('refresh-secret')),
              )
              .having(
                (error) => error.toString(),
                'client secret',
                isNot(contains('client_secret=secret')),
              ),
        ),
      );
    });
  });
}

TwoFinanceAuthClient _client(FakeTransport transport) {
  return TwoFinanceAuthClient(
    baseUrl: Uri.parse('http://localhost:8080'),
    realm: 'realm',
    clientId: 'client',
    phoneClientId: 'phone-client',
    transport: transport,
  );
}

Map<String, Object?> _jwtJson() {
  return {
    'access_token': 'access',
    'expires_in': 300,
    'refresh_token': 'refresh',
    'refresh_expires_in': 1800,
    'token_type': 'Bearer',
  };
}

class FakeTransport implements AuthTransport {
  AuthHttpResponse? nextGet;
  AuthHttpResponse? nextPost;

  late Uri lastUrl;
  Map<String, Object?>? lastBody;
  String? lastAccessToken;
  List<AuthCookie> lastCookies = const [];
  bool? lastFollowRedirects;

  @override
  Future<AuthHttpResponse> get(
    Uri url, {
    String? accessToken,
    List<AuthCookie> cookies = const [],
    bool followRedirects = true,
  }) async {
    lastUrl = url;
    lastAccessToken = accessToken;
    lastCookies = cookies;
    lastFollowRedirects = followRedirects;
    return nextGet!;
  }

  @override
  Future<AuthHttpResponse> post(
    Uri url, {
    required Map<String, Object?> body,
    String? accessToken,
  }) async {
    lastUrl = url;
    lastBody = body;
    lastAccessToken = accessToken;
    return nextPost!;
  }
}
