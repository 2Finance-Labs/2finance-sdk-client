import 'models.dart';

final class SdkConfig {
  const SdkConfig({
    this.authUrl = '',
    this.networkUrl = '',
    this.analyticsUrl = '',
    this.orchestratorUrl = '',
    this.mcpUrl = '',
    this.tradingControlUrl = '',
    this.matchEngineWsUrl = '',
    this.keystoreUrl = '',
    this.hummingbotUrl = '',
    this.wiseUrl = '',
    this.airwallexUrl = '',
    this.authRealm = '2finance',
    this.authClientId = '2finance-network',
    this.authPhoneClientId = '2finance-network-phone',
  });

  final String authUrl;
  final String networkUrl;
  final String analyticsUrl;
  final String orchestratorUrl;
  final String mcpUrl;
  final String tradingControlUrl;
  final String matchEngineWsUrl;
  final String keystoreUrl;
  final String hummingbotUrl;
  final String wiseUrl;
  final String airwallexUrl;
  final String authRealm;
  final String authClientId;
  final String authPhoneClientId;

  factory SdkConfig.fromMap(Map<String, String> env) {
    return SdkConfig(
      authUrl: env['TWO_FINANCE_AUTH_URL'] ?? '',
      networkUrl: env['TWO_FINANCE_NETWORK_URL'] ?? '',
      analyticsUrl: env['TWO_FINANCE_ANALYTICS_URL'] ?? '',
      orchestratorUrl: env['TWO_FINANCE_ORCHESTRATOR_URL'] ?? '',
      mcpUrl: env['TWO_FINANCE_MCP_URL'] ?? '',
      tradingControlUrl: env['TWO_FINANCE_TRADING_CONTROL_URL'] ?? '',
      matchEngineWsUrl: env['TWO_FINANCE_MATCHENGINE_WS_URL'] ?? '',
      keystoreUrl: env['TWO_FINANCE_KEYSTORE_URL'] ?? '',
      hummingbotUrl: env['TWO_FINANCE_HUMMINGBOT_URL'] ?? '',
      wiseUrl: env['TWO_FINANCE_WISE_URL'] ?? '',
      airwallexUrl: env['TWO_FINANCE_AIRWALLEX_URL'] ?? '',
      authRealm: _envValue(env, 'TWO_FINANCE_AUTH_REALM', '2finance'),
      authClientId: _envValue(
        env,
        'TWO_FINANCE_AUTH_CLIENT_ID',
        '2finance-network',
      ),
      authPhoneClientId: _envValue(
        env,
        'TWO_FINANCE_AUTH_PHONE_CLIENT_ID',
        '2finance-network-phone',
      ),
    );
  }

  factory SdkConfig.fromEnvironment() {
    return SdkConfig.fromMap(const {
      'TWO_FINANCE_AUTH_URL': String.fromEnvironment('TWO_FINANCE_AUTH_URL'),
      'TWO_FINANCE_NETWORK_URL': String.fromEnvironment(
        'TWO_FINANCE_NETWORK_URL',
      ),
      'TWO_FINANCE_ANALYTICS_URL': String.fromEnvironment(
        'TWO_FINANCE_ANALYTICS_URL',
      ),
      'TWO_FINANCE_ORCHESTRATOR_URL': String.fromEnvironment(
        'TWO_FINANCE_ORCHESTRATOR_URL',
      ),
      'TWO_FINANCE_MCP_URL': String.fromEnvironment('TWO_FINANCE_MCP_URL'),
      'TWO_FINANCE_TRADING_CONTROL_URL': String.fromEnvironment(
        'TWO_FINANCE_TRADING_CONTROL_URL',
      ),
      'TWO_FINANCE_MATCHENGINE_WS_URL': String.fromEnvironment(
        'TWO_FINANCE_MATCHENGINE_WS_URL',
      ),
      'TWO_FINANCE_KEYSTORE_URL': String.fromEnvironment(
        'TWO_FINANCE_KEYSTORE_URL',
      ),
      'TWO_FINANCE_HUMMINGBOT_URL': String.fromEnvironment(
        'TWO_FINANCE_HUMMINGBOT_URL',
      ),
      'TWO_FINANCE_WISE_URL': String.fromEnvironment('TWO_FINANCE_WISE_URL'),
      'TWO_FINANCE_AIRWALLEX_URL': String.fromEnvironment(
        'TWO_FINANCE_AIRWALLEX_URL',
      ),
      'TWO_FINANCE_AUTH_REALM': String.fromEnvironment(
        'TWO_FINANCE_AUTH_REALM',
      ),
      'TWO_FINANCE_AUTH_CLIENT_ID': String.fromEnvironment(
        'TWO_FINANCE_AUTH_CLIENT_ID',
      ),
      'TWO_FINANCE_AUTH_PHONE_CLIENT_ID': String.fromEnvironment(
        'TWO_FINANCE_AUTH_PHONE_CLIENT_ID',
      ),
    });
  }

  String serviceUrl(String domain) {
    switch (_serviceKey(domain)) {
      case 'auth':
        return authUrl;
      case 'network':
        return networkUrl;
      case 'analytics':
        return analyticsUrl;
      case 'orchestrator':
        return orchestratorUrl;
      case 'mcp':
      case 'planner':
        return mcpUrl;
      case 'tradingcontrol':
        return tradingControlUrl;
      case 'matchengine':
        return matchEngineWsUrl;
      case 'keystore':
        return keystoreUrl;
      case 'hummingbot':
        return hummingbotUrl;
      case 'wise':
        return wiseUrl;
      case 'airwallex':
        return airwallexUrl;
      default:
        return '';
    }
  }

  Map<String, String> serviceUrls() {
    final urls = <String, String>{};
    for (final service in defaultServiceCatalog.services) {
      final url = serviceUrl(service.name);
      if (url.isNotEmpty) {
        urls[service.name] = url;
      }
    }
    return urls;
  }

  List<ConfiguredServiceEntry> configuredServices() {
    final services = <ConfiguredServiceEntry>[];
    for (final service in defaultServiceCatalog.services) {
      final url = serviceUrl(service.name);
      if (url.isNotEmpty) {
        services.add(
          ConfiguredServiceEntry(
            name: service.name,
            env: service.env,
            url: url,
          ),
        );
      }
    }
    return services;
  }

  List<ServiceCatalogEntry> missingServiceUrls() {
    final services = <ServiceCatalogEntry>[];
    for (final service in defaultServiceCatalog.services) {
      if (serviceUrl(service.name).isEmpty) {
        services.add(service);
      }
    }
    return services;
  }
}

String _envValue(Map<String, String> env, String key, String fallback) {
  final value = env[key];
  return value == null || value.isEmpty ? fallback : value;
}

String _serviceKey(String domain) {
  return domain.trim().toLowerCase().replaceAll(RegExp(r'[-_\s]+'), '');
}
