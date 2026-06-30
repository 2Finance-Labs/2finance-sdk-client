class AuthSdkException implements Exception {
  AuthSdkException(String message, {this.statusCode, String? body})
    : message = _redactSensitiveText(message),
      body = body == null ? null : _redactSensitiveText(body);

  final String message;
  final int? statusCode;
  final String? body;

  @override
  String toString() {
    final code = statusCode == null ? '' : ' (status $statusCode)';
    final responseBody = body == null || body!.isEmpty ? '' : ': $body';
    return 'AuthSdkException$code: $message$responseBody';
  }
}

String _redactSensitiveText(String value) {
  var redacted = value.replaceAllMapped(
    RegExp(r'\bBearer\s+[A-Za-z0-9._~+/=-]+', caseSensitive: false),
    (_) => 'Bearer <redacted>',
  );

  redacted = redacted.replaceAllMapped(
    RegExp(
      r'''(["']?(?:access_token|refresh_token|id_token|accessToken|refreshToken|idToken|client_secret|clientSecret|password)["']?\s*[:=]\s*["']?)([^"',\s}]+)''',
      caseSensitive: false,
    ),
    (match) => '${match.group(1)}<redacted>',
  );

  return redacted;
}

class AuthCookie {
  const AuthCookie({required this.name, required this.value});

  final String name;
  final String value;

  factory AuthCookie.fromSetCookieValue(String value) {
    final firstPart = value.split(';').first;
    final separator = firstPart.indexOf('=');
    if (separator <= 0) {
      throw FormatException('Invalid Set-Cookie header: $value');
    }
    return AuthCookie(
      name: firstPart.substring(0, separator).trim(),
      value: firstPart.substring(separator + 1).trim(),
    );
  }
}

class DefaultResponse {
  const DefaultResponse({this.code, this.msg, this.data});

  final int? code;
  final String? msg;
  final Object? data;

  factory DefaultResponse.fromJson(Map<String, Object?> json) {
    return DefaultResponse(
      code: json['code'] as int?,
      msg: json['msg'] as String?,
      data: json['data'],
    );
  }
}

@Deprecated(
  'Password login is kept for compatibility with controlled environments. '
  'Use Authorization Code + PKCE with TwoFinanceAuthClient.loginPKCE and '
  'callbackPKCE instead.',
)
class LoginInput {
  const LoginInput({required this.username, required this.password});

  final String username;
  final String password;

  Map<String, Object?> toJson() {
    return {'username': username, 'password': password};
  }
}

class CreateUserInput {
  const CreateUserInput({
    required this.username,
    required this.email,
    this.firstName,
    this.lastName,
    this.attributes,
    this.password,
  });

  final String username;
  final String email;
  final String? firstName;
  final String? lastName;
  final Map<String, List<String>>? attributes;
  final String? password;

  Map<String, Object?> toJson() {
    return {
      'username': username,
      'email': email,
      if (firstName != null) 'firstName': firstName,
      if (lastName != null) 'lastName': lastName,
      if (attributes != null) 'attributes': attributes,
      if (password != null) 'password': password,
    };
  }
}

class SMSRequest {
  const SMSRequest({required this.phoneNumber, required this.userId});

  final String phoneNumber;
  final String userId;

  Map<String, Object?> toJson() {
    return {'phone_number': phoneNumber, 'user_id': userId};
  }
}

class VerifySMSRequest {
  const VerifySMSRequest({
    required this.phoneNumber,
    required this.code,
    required this.userId,
  });

  final String phoneNumber;
  final String code;
  final String userId;

  Map<String, Object?> toJson() {
    return {'phone_number': phoneNumber, 'code': code, 'user_id': userId};
  }
}

class JWT {
  const JWT({
    required this.accessToken,
    required this.expiresIn,
    required this.refreshToken,
    required this.refreshExpiresIn,
    required this.tokenType,
  });

  final String accessToken;
  final int expiresIn;
  final String refreshToken;
  final int refreshExpiresIn;
  final String tokenType;

  factory JWT.fromJson(Map<String, Object?> json) {
    return JWT(
      accessToken: json['access_token'] as String? ?? '',
      expiresIn: json['expires_in'] as int? ?? 0,
      refreshToken: json['refresh_token'] as String? ?? '',
      refreshExpiresIn: json['refresh_expires_in'] as int? ?? 0,
      tokenType: json['token_type'] as String? ?? '',
    );
  }
}

class PKCELoginResponse {
  const PKCELoginResponse({
    required this.authUrl,
    required this.state,
    required this.cookies,
    required this.statusCode,
  });

  final Uri authUrl;
  final String state;
  final List<AuthCookie> cookies;
  final int statusCode;
}

class PKCECallbackResponse {
  const PKCECallbackResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.idToken,
    required this.expiresIn,
    required this.tokenType,
  });

  final String accessToken;
  final String refreshToken;
  final String idToken;
  final int expiresIn;
  final String tokenType;

  factory PKCECallbackResponse.fromJson(Map<String, Object?> json) {
    return PKCECallbackResponse(
      accessToken: json['access_token'] as String? ?? '',
      refreshToken: json['refresh_token'] as String? ?? '',
      idToken: json['id_token'] as String? ?? '',
      expiresIn: json['expires_in'] as int? ?? 0,
      tokenType: json['token_type'] as String? ?? '',
    );
  }
}
