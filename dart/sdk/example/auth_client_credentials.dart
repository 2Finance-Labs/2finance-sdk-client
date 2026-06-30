import 'package:two_finance_sdk_client/two_finance_sdk_client.dart';

Future<void> main() async {
  final tokenSource = ClientCredentialsTokenSource(
    tokenUrl: const String.fromEnvironment('TWO_FINANCE_AUTH_TOKEN_URL'),
    clientId: const String.fromEnvironment('TWO_FINANCE_AUTH_CLIENT_ID'),
    clientSecret: const String.fromEnvironment(
      'TWO_FINANCE_AUTH_CLIENT_SECRET',
    ),
    scopes: const ['2finance.sdk'],
    transport: (request) async {
      return const AuthTokenResponse(
        statusCode: 200,
        body: '{"access_token":"example-token","expires_in":300}',
      );
    },
  );

  final client = TwoFinanceClient.fromEnvironment(tokenSource: tokenSource);
  await client.analytics.indicators();
}
