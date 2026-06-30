import 'package:two_finance_sdk_client/two_finance_sdk_client.dart';

Future<void> main() async {
  final client = TwoFinanceClient.fromEnvironment();

  final response = await client.analytics.service.post(
    '/analytics/candles:upsert',
    {'symbol': 'BTC-USDT'},
    const RequestOptions(
      headers: {'X-Trace-ID': 'trace-1'},
      idempotencyKey: 'candles-upsert-001',
      queryParameters: {'source': 'sdk-example'},
      page: 1,
      limit: 25,
      timeout: Duration(seconds: 5),
      maxRetries: 1,
    ),
  );
  print('response: $response');
}
