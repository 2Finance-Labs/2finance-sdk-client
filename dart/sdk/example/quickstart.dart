import 'package:two_finance_sdk_client/two_finance_sdk_client.dart';

Future<void> main() async {
  final client = TwoFinanceClient.fromEnvironment();

  final indicators = await client.analytics.indicators();
  print('analytics indicators: $indicators');

  final plan = await client.planner.tradingPlan({
    'goal': 'prepare a BTC rebalancing plan',
    'useAnalytics': true,
    'useTrading': true,
  });
  print('planner response: $plan');
}
