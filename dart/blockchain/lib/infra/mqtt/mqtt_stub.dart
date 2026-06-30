import 'package:two_finance_blockchain/infra/transport/transport.dart';

typedef MessageHandler = void Function(Object client, Object message);

abstract class MqttClientInterface implements FinanceNetworkTransport {
  Future<void> connect();
  Future<void> disconnect();
  Future<void> publish(String topic, String payload);
  Future<void> subscribe(String topic, {MessageHandler? handler});
  Future<void> unsubscribe(String topic);
  Object? get client;
}

class MqttClientWrapper implements MqttClientInterface {
  MqttClientWrapper({
    required this.host,
    required this.port,
    required this.clientId,
    this.useSSL = false,
    this.username,
    this.password,
    this.caCertPath,
  });

  final String host;
  final String port;
  final String clientId;
  final bool useSSL;
  final String? username;
  final String? password;
  final String? caCertPath;

  @override
  Object? get client => null;

  @override
  Future<void> connect() => _unsupported();

  @override
  Future<void> disconnect() async {}

  @override
  Future<void> publish(String topic, String payload) => _unsupported();

  @override
  Future<void> subscribe(String topic, {MessageHandler? handler}) =>
      _unsupported();

  @override
  Future<void> unsubscribe(String topic) async {}

  @override
  Future<dynamic> sendRequest(String method, dynamic params, String replyTo) =>
      _unsupported();

  Never _unsupported() {
    throw UnsupportedError(
      'MQTT transport is not available on this platform. Use HttpFinanceNetworkTransport.',
    );
  }
}
