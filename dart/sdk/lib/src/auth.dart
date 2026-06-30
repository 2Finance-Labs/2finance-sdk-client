import 'dart:convert';

abstract interface class TokenSource {
  Future<String> token();
}

typedef AuthTokenTransport =
    Future<AuthTokenResponse> Function(AuthTokenRequest request);

final class AuthTokenRequest {
  const AuthTokenRequest({
    required this.url,
    required this.headers,
    required this.body,
  });

  final Uri url;
  final Map<String, String> headers;
  final String body;
}

final class AuthTokenResponse {
  const AuthTokenResponse({required this.statusCode, required this.body});

  final int statusCode;
  final String body;
}

final class StaticTokenSource implements TokenSource {
  StaticTokenSource(this.accessToken);

  final String accessToken;

  @override
  Future<String> token() async => accessToken;
}

final class ClientCredentialsTokenSource implements TokenSource {
  ClientCredentialsTokenSource({
    required this.tokenUrl,
    required this.clientId,
    required this.clientSecret,
    this.scopes = const [],
    this.transport,
    this.expirySkew = const Duration(seconds: 30),
  });

  final String tokenUrl;
  final String clientId;
  final String clientSecret;
  final List<String> scopes;
  final AuthTokenTransport? transport;
  final Duration expirySkew;

  String? _accessToken;
  DateTime? _expiresAt;

  @override
  Future<String> token() async {
    final cached = _accessToken;
    final expiresAt = _expiresAt;
    if (cached != null &&
        expiresAt != null &&
        DateTime.now().isBefore(expiresAt.subtract(expirySkew))) {
      return cached;
    }
    if (tokenUrl.trim().isEmpty ||
        clientId.trim().isEmpty ||
        clientSecret.trim().isEmpty) {
      throw StateError(
        '2finance auth: tokenUrl, clientId and clientSecret are required',
      );
    }
    final authTransport = transport;
    if (authTransport == null) {
      throw StateError('2finance auth: transport is required');
    }
    final form = <String, String>{
      'grant_type': 'client_credentials',
      'client_id': clientId,
      'client_secret': clientSecret,
      if (scopes.isNotEmpty) 'scope': scopes.join(' '),
    };
    final response = await authTransport(
      AuthTokenRequest(
        url: Uri.parse(tokenUrl),
        headers: const {
          'Accept': 'application/json',
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: form.entries
            .map(
              (entry) =>
                  '${Uri.encodeQueryComponent(entry.key)}=${Uri.encodeQueryComponent(entry.value)}',
            )
            .join('&'),
      ),
    );
    if (response.statusCode < 200 || response.statusCode >= 300) {
      throw StateError(
        '2finance auth: token endpoint returned ${response.statusCode}',
      );
    }
    final payload = jsonDecode(response.body) as Map<String, Object?>;
    final accessToken = payload['access_token'] as String?;
    if (accessToken == null || accessToken.isEmpty) {
      throw StateError('2finance auth: token response missing access_token');
    }
    final expiresIn = payload['expires_in'];
    final seconds = expiresIn is num ? expiresIn.toInt() : 300;
    _accessToken = accessToken;
    _expiresAt = DateTime.now().add(Duration(seconds: seconds));
    return accessToken;
  }
}

String bearerAuthorization(String accessToken) {
  final trimmed = accessToken.trim();
  if (trimmed.isEmpty) {
    return '';
  }
  if (trimmed.toLowerCase().startsWith('bearer ')) {
    return trimmed;
  }
  return 'Bearer $trimmed';
}
