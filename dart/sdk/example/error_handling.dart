import 'package:two_finance_sdk_client/two_finance_sdk_client.dart';

Future<void> main() async {
  final client = TwoFinanceClient.fromEnvironment();

  try {
    await client.analytics.indicators();
  } on ServiceException catch (error) {
    print('request failed with status ${error.statusCode}: ${error.body}');
  }
}
