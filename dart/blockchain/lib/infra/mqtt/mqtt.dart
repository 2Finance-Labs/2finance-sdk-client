// lib/mqtt/mqtt_client_wrapper.dart
import 'dart:io';

import 'dart:async';
import 'dart:convert';

import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';
import 'package:two_finance_blockchain/infra/event/request_response.dart';
import 'package:two_finance_blockchain/infra/transport/transport.dart';
import 'package:two_finance_blockchain/infra/utils/tls_utils.dart';

typedef MessageHandler =
    void Function(MqttClient client, MqttReceivedMessage<MqttMessage> message);
typedef MqttServerClientFactory =
    MqttServerClient Function(String host, String clientId);

abstract class MqttClientInterface implements FinanceNetworkTransport {
  Future<void> connect();
  Future<void> disconnect();
  Future<void> publish(String topic, String payload);
  Future<void> subscribe(String topic, {MessageHandler? handler});
  Future<void> unsubscribe(String topic);
  MqttClient? get client;
}

class MqttClientWrapper implements MqttClientInterface {
  final String host;
  final String port;
  final String clientId;
  final bool useSSL;
  final String? username;
  final String? password;
  final String? caCertPath;
  MqttServerClient? _client;

  MqttClientWrapper({
    required this.host,
    required this.port,
    required this.clientId,
    this.useSSL = false,
    this.username,
    this.password,
    this.caCertPath,
    MqttServerClientFactory? clientFactory,
  }) : clientFactory = clientFactory ?? ((h, id) => MqttServerClient(h, id));

  final MqttServerClientFactory clientFactory;

  @override
  MqttClient? get client => _client;

  @override
  Future<void> connect() async {
    _client = clientFactory(host, clientId);
    _client!
      ..logging(on: false)
      ..keepAlivePeriod = 60
      ..onDisconnected = _onDisconnected
      ..onConnected = _onConnected
      ..onSubscribed = _onSubscribed
      ..pongCallback = _onPong
      ..autoReconnect = true
      ..resubscribeOnAutoReconnect = true;

    final parsedPort = int.tryParse(port);
    if (parsedPort == null) {
      print(
        '⚠️ Warning: Invalid MQTT port "$port", using default ${useSSL ? 8883 : 1883}',
      );
    }
    _client!.port = parsedPort ?? (useSSL ? 8883 : 1883);

    if (useSSL) {
      _client!.secure = true;

      if (caCertPath != null && caCertPath!.isNotEmpty) {
        _client!.securityContext = createSecurityContext(caCertPath!);
      } else {
        _client!.securityContext = SecurityContext.defaultContext;
      }
    }

    final connMess = MqttConnectMessage()
        .authenticateAs(username ?? '', password ?? '')
        .withClientIdentifier(clientId)
        .withWillTopic('will/topic')
        .withWillMessage('client disconnected unexpectedly')
        .startClean()
        .withWillQos(MqttQos.atLeastOnce);

    _client!.connectionMessage = connMess;

    try {
      await _client!.connect();
    } catch (e) {
      print('Connection error: $e');
      _client!.disconnect();
      rethrow;
    }

    if (_client!.connectionStatus?.state != MqttConnectionState.connected) {
      throw Exception('Failed to connect to MQTT broker');
    }

    print('Connected to MQTT broker at $host');
  }

  @override
  Future<void> disconnect() async {
    _client?.disconnect();
    print('Disconnected from MQTT broker');
  }

  @override
  Future<void> publish(String topic, String payload) async {
    final builder = MqttClientPayloadBuilder()..addString(payload);
    _client?.publishMessage(topic, MqttQos.atLeastOnce, builder.payload!);
  }

  @override
  Future<void> subscribe(String topic, {MessageHandler? handler}) async {
    _client?.subscribe(topic, MqttQos.atLeastOnce);

    if (handler != null) {
      _client?.updates?.listen((List<MqttReceivedMessage<MqttMessage>> c) {
        for (var message in c) {
          handler(_client!, message);
        }
      });
    }
  }

  @override
  Future<void> unsubscribe(String topic) async {
    _client?.unsubscribe(topic);
  }

  @override
  Future<dynamic> sendRequest(
    String method,
    dynamic params,
    String replyTo,
  ) async {
    final topicBase = TRANSACTIONS_REQUEST_TOPIC.replaceAll('/+', '');
    final requestTopic = "$topicBase/$replyTo";
    final responseTopic = "$TRANSACTIONS_RESPONSE_TOPIC/$replyTo";

    final responseCompleter = Completer<String>();

    await subscribe(
      responseTopic,
      handler: (client, msg) {
        final publishMessage = msg.payload as MqttPublishMessage;
        final payloadStr = MqttPublishPayload.bytesToStringAsString(
          publishMessage.payload.message,
        );
        if (!responseCompleter.isCompleted) {
          responseCompleter.complete(payloadStr);
        }
      },
    );

    final payload = RequestPayload(method: method, params: params);
    await publish(requestTopic, json.encode(payload.toJson()));

    final responseBytes = await responseCompleter.future.timeout(
      const Duration(seconds: 10),
      onTimeout: () => throw Exception(
        "timeout waiting for response on topic $responseTopic",
      ),
    );
    final decoded = json.decode(responseBytes);
    final resp = ResponsePayload.fromJson(Map<String, dynamic>.from(decoded));

    if (resp.status == RESPONSE_STATUS_ERROR) {
      if (resp.message?.contains("record not found") == true) {
        return 0;
      }
      throw Exception("error in response: ${resp.message}");
    }

    return resp.data;
  }

  void _onConnected() => print('MQTT connected');
  void _onDisconnected() => print('MQTT disconnected');
  void _onSubscribed(String topic) => print('Subscribed to $topic');
  void _onPong() => print('Pong received from broker');
}
