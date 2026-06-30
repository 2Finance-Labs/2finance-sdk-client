import 'models.dart';
import 'transport.dart';

class TwoFinanceAuthClient {
  TwoFinanceAuthClient({
    required Uri baseUrl,
    required String realm,
    required String clientId,
    required String phoneClientId,
    AuthTransport? transport,
  }) : _baseUrl = baseUrl,
       _realm = realm,
       _clientId = clientId,
       _phoneClientId = phoneClientId,
       _transport = transport ?? HttpAuthTransport();

  final Uri _baseUrl;
  final String _realm;
  final String _clientId;
  final String _phoneClientId;
  final AuthTransport _transport;

  @Deprecated(
    'Password login is kept for compatibility with controlled environments. '
    'Use Authorization Code + PKCE with loginPKCE and callbackPKCE instead.',
  )
  Future<JWT> login(LoginInput input) async {
    final response = await _transport.post(
      _url('/login'),
      body: input.toJson(),
    );
    ensureSuccess(response, 'Login failed');
    return _jwtFromDefaultResponse(response.body);
  }

  Future<DefaultResponse> signUp(CreateUserInput input) async {
    final response = await _transport.post(
      _url('/signup'),
      body: input.toJson(),
    );
    ensureSuccess(response, 'Signup failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<JWT> refreshToken(String refreshToken) async {
    final response = await _transport.post(
      _url('/refresh'),
      body: {'refresh_token': refreshToken},
    );
    ensureSuccess(response, 'Refresh token failed');
    return _jwtFromDefaultResponse(response.body);
  }

  Future<DefaultResponse> logout(String refreshToken) async {
    final response = await _transport.post(
      _url('/logout'),
      body: {'refresh_token': refreshToken},
    );
    ensureSuccess(response, 'Logout failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<DefaultResponse> requestAuthenticationCode(String phoneNumber) async {
    final response = await _transport.post(
      _phoneUrl('/phone/sms/request-code'),
      body: {'phone_number': phoneNumber},
    );
    ensureSuccess(response, 'Request authentication code failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<JWT> phoneLogin({
    required String phoneNumber,
    required String code,
  }) async {
    final response = await _transport.post(
      _phoneUrl('/phone/sms/login'),
      body: {'phone_number': phoneNumber, 'code': code},
    );
    ensureSuccess(response, 'Phone login failed');
    return _jwtFromDefaultResponse(response.body);
  }

  Future<PKCELoginResponse> loginPKCE() async {
    final response = await _transport.get(
      _url('/pkce/login-redirect'),
      followRedirects: false,
    );
    if (response.statusCode < 300 || response.statusCode >= 400) {
      throw AuthSdkException(
        'PKCE login redirect failed',
        statusCode: response.statusCode,
        body: response.body,
      );
    }

    final location = response.headers['location'];
    if (location == null || location.isEmpty) {
      throw AuthSdkException('PKCE login redirect missing Location header');
    }

    final authUrl = Uri.parse(location);
    return PKCELoginResponse(
      authUrl: authUrl,
      state: authUrl.queryParameters['state'] ?? '',
      cookies: response.cookies,
      statusCode: response.statusCode,
    );
  }

  Future<PKCECallbackResponse> callbackPKCE({
    required String code,
    required String state,
    List<AuthCookie> cookies = const [],
  }) async {
    if (code.isEmpty) {
      throw AuthSdkException('PKCE callback code is required');
    }
    if (state.isEmpty) {
      throw AuthSdkException('PKCE callback state is required');
    }

    final response = await _transport.get(
      _url(
        '/pkce/callback',
      ).replace(queryParameters: {'code': code, 'state': state}),
      cookies: cookies,
    );
    ensureSuccess(response, 'PKCE callback failed');
    return PKCECallbackResponse.fromJson(decodeObject(response.body));
  }

  Future<DefaultResponse> getUserInfo(String accessToken) async {
    final response = await _transport.get(
      _url('/user-info'),
      accessToken: accessToken,
    );
    ensureSuccess(response, 'Get user info failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<DefaultResponse> createUser(
    String accessToken,
    CreateUserInput input,
  ) async {
    final response = await _transport.post(
      _url('/create-user'),
      body: input.toJson(),
      accessToken: accessToken,
    );
    ensureSuccess(response, 'Create user failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<DefaultResponse> requestSMSCode(
    String accessToken,
    SMSRequest request,
  ) async {
    final response = await _transport.post(
      _url('/request-sms'),
      body: request.toJson(),
      accessToken: accessToken,
    );
    ensureSuccess(response, 'Request SMS code failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Future<DefaultResponse> verifySMSCode(
    String accessToken,
    VerifySMSRequest request,
  ) async {
    final response = await _transport.post(
      _url('/verify-sms'),
      body: request.toJson(),
      accessToken: accessToken,
    );
    ensureSuccess(response, 'Verify SMS code failed');
    return DefaultResponse.fromJson(decodeObject(response.body));
  }

  Uri _url(String endpoint) {
    return _baseUrl.replace(
      path: _joinPath([
        _baseUrl.path,
        'v1',
        '2finance-authenticator',
        _realm,
        _clientId,
        endpoint,
      ]),
    );
  }

  Uri _phoneUrl(String endpoint) {
    return _baseUrl.replace(
      path: _joinPath([
        _baseUrl.path,
        'v1',
        '2finance-authenticator',
        _realm,
        _phoneClientId,
        endpoint,
      ]),
    );
  }

  JWT _jwtFromDefaultResponse(String body) {
    final response = DefaultResponse.fromJson(decodeObject(body));
    final data = response.data;
    if (data is Map<String, Object?>) {
      return JWT.fromJson(data);
    }
    if (data is Map) {
      return JWT.fromJson(Map<String, Object?>.from(data));
    }
    throw AuthSdkException('Expected JWT in response data', body: body);
  }
}

String _joinPath(List<String> parts) {
  return parts
      .expand((part) => part.split('/'))
      .where((part) => part.isNotEmpty)
      .join('/')
      .replaceFirst(RegExp('^'), '/');
}
