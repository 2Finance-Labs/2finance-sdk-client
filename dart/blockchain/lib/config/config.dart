import 'package:dotenv/dotenv.dart';

class Config {
  static late final String emqxScheme;
  static late final String emqxHost;
  static late final String emqxPort;
  static late final String emqxUsername;
  static late final String emqxPassword;
  static late final bool emqxSSL;
  static late final String emqxClientId;
  static late final String emqxCaCertPath;

  static late final String keycloakClientId;
  static late final String keycloakClientSecret;
  static late final String keycloakState;
  static late final String keycloakHostname;
  static late final String keycloakRealm;
  static late final String keycloakRedirectUrl;

  static bool _isInitialized = false;

  static bool get isInitialized => _isInitialized;

  static Future<void> loadConfig({
    required String env,
    String path = '.env',
  }) async {
    if (_isInitialized) return;
    print('[Config] Initializing config for ENV: $env');
    var environment = DotEnv(includePlatformEnvironment: true);
    if (env == 'prod') {
      try {
        //await dotenv.load(fileName: path);
        environment.load([path]);
      } catch (e) {
        print('[Config] Failed to load .env file at $path: $e');
        throw Exception('Missing or invalid .env file');
      }

      emqxScheme = environment['EMQX_SCHEME'] ?? 'tcp';
      emqxHost = environment['EMQX_HOST'] ?? 'localhost';
      emqxPort = environment['EMQX_PORT'] ?? '1883';
      emqxSSL = (environment['EMQX_SSL']?.toLowerCase() ?? 'false') == 'true';
      emqxUsername = environment['EMQX_USERNAME'] ?? '';
      emqxPassword = environment['EMQX_PASSWORD'] ?? '';
      emqxClientId = environment['EMQX_CLIENT_ID'] ?? 'flutter_plugin_client';
      emqxCaCertPath = environment['EMQX_CA_CERT_PATH'] ?? '';

      keycloakClientId = environment['KEYCLOAK_CLIENT_ID'] ?? '';
      keycloakClientSecret = environment['KEYCLOAK_CLIENT_SECRET'] ?? '';
      keycloakState = environment['KEYCLOAK_STATE'] ?? '';
      keycloakRealm = environment['KEYCLOAK_REALM'] ?? '';
      keycloakHostname = environment['KEYCLOAK_HOSTNAME'] ?? '';
      keycloakRedirectUrl = environment['KEYCLOAK_REDIRECT_URL'] ?? '';
      _isInitialized = true;

      return;
    }
    // EMQX_SCHEME="tcp"
    // EMQX_HOST="k612abc6.ala.us-east-1.emqxsl.com"
    // EMQX_PORT="1883" # Use uma porta padrão para teste, ou 000 como placeholder
    // EMQX_USERNAME="lmenniti"
    // EMQX_PASSWORD="123123123"
    // EMQX_CA_CERT_PATH="../certificates/emqxsl-ca.crt"
    // EMQX_CLIENT_ID="2finance-network"
    //EMQX_SSL=false
    emqxScheme =
        environment['MQTT_SCHEME'] ?? environment['EMQX_SCHEME'] ?? 'tcp';
    emqxHost =
        environment['MQTT_HOST'] ?? environment['EMQX_HOST'] ?? '127.0.0.1';
    emqxPort = environment['MQTT_PORT'] ?? environment['EMQX_PORT'] ?? '1883';
    emqxSSL =
        ((environment['MQTT_SSL'] ?? environment['EMQX_SSL'])?.toLowerCase() ??
            'false') ==
        'true';
    emqxUsername =
        environment['MQTT_USERNAME'] ?? environment['EMQX_USERNAME'] ?? '';
    emqxPassword =
        environment['MQTT_PASSWORD'] ?? environment['EMQX_PASSWORD'] ?? '';
    emqxClientId =
        environment['MQTT_CLIENT_ID'] ??
        environment['EMQX_CLIENT_ID'] ??
        'flutter_plugin_client';
    emqxCaCertPath =
        environment['MQTT_CA_CERT_PATH'] ??
        environment['EMQX_CA_CERT_PATH'] ??
        '';

    keycloakClientId = '';
    keycloakClientSecret = '';
    keycloakState = '';
    keycloakRealm = '';
    keycloakHostname = '';
    keycloakRedirectUrl = '';
    _isInitialized = true;
  }
}
