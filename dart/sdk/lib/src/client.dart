import 'auth.dart';
import 'config.dart';
import 'domains.dart';
import 'service_client.dart';

final class TwoFinanceClient {
  TwoFinanceClient(
    this.config, {
    Transport? transport,
    TokenSource? tokenSource,
  }) : auth = AuthClient(
         ServiceClient(
           config.authUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
         realm: config.authRealm,
         clientId: config.authClientId,
         phoneClientId: config.authPhoneClientId,
       ),
       network = NetworkClient(
         ServiceClient(
           config.networkUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       analytics = AnalyticsClient(
         ServiceClient(
           config.analyticsUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       orchestrator = OrchestratorClient(
         ServiceClient(
           config.orchestratorUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       mcp = MCPClient(
         ServiceClient(
           config.mcpUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       tradingControl = TradingControlClient(
         ServiceClient(
           config.tradingControlUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       matchEngine = MatchEngineClient(config.matchEngineWsUrl),
       keystore = KeyStoreClient(
         ServiceClient(
           config.keystoreUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       hummingbot = HummingbotClient(
         ServiceClient(
           config.hummingbotUrl,
           transport: transport,
           tokenSource: tokenSource,
         ),
       ),
       providers = ProvidersClient(
         wise: WiseClient(
           ServiceClient(
             config.wiseUrl,
             transport: transport,
             tokenSource: tokenSource,
           ),
         ),
         airwallex: AirwallexClient(
           ServiceClient(
             config.airwallexUrl,
             transport: transport,
             tokenSource: tokenSource,
           ),
         ),
       ) {
    planner = PlannerClient(
      mcp: mcp,
      orchestrator: orchestrator,
      analytics: analytics,
      tradingControl: tradingControl,
    );
  }

  factory TwoFinanceClient.fromEnvironment({
    Transport? transport,
    TokenSource? tokenSource,
  }) {
    return TwoFinanceClient(
      SdkConfig.fromEnvironment(),
      transport: transport,
      tokenSource: tokenSource,
    );
  }

  final SdkConfig config;
  final AuthClient auth;
  final NetworkClient network;
  final AnalyticsClient analytics;
  final OrchestratorClient orchestrator;
  final MCPClient mcp;
  final TradingControlClient tradingControl;
  final MatchEngineClient matchEngine;
  final KeyStoreClient keystore;
  final HummingbotClient hummingbot;
  final ProvidersClient providers;
  late final PlannerClient planner;
}
