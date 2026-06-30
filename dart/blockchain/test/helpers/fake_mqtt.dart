import 'dart:convert';

import 'package:mqtt_client/mqtt_client.dart';
import 'package:two_finance_blockchain/infra/event/request_response.dart';
import 'package:two_finance_blockchain/infra/mqtt/mqtt.dart';

typedef FakeResponseBuilder =
    Map<String, dynamic> Function(Map<String, dynamic> request);

class FakeMqttClient implements MqttClientInterface {
  FakeMqttClient({FakeResponseBuilder? responseBuilder})
    : _responseBuilder =
          responseBuilder ??
          ((_) => {
            'status': RESPONSE_STATUS_SUCCESS,
            'message': null,
            'data': {'states': <dynamic>[], 'logs': <dynamic>[]},
          });

  final FakeResponseBuilder _responseBuilder;
  final MqttClient _client = MqttClient('localhost', 'fake-client');
  final List<Map<String, dynamic>> publishedRequests = <Map<String, dynamic>>[];
  final Map<String, MessageHandler?> _handlers = <String, MessageHandler?>{};

  @override
  MqttClient? get client => _client;

  Map<String, dynamic> get lastRequest => publishedRequests.last;

  @override
  Future<void> connect() async {}

  @override
  Future<void> disconnect() async {}

  @override
  Future<void> publish(String topic, String payload) async {
    final request = json.decode(payload) as Map<String, dynamic>;
    publishedRequests.add(request);

    final replyTo = topic.split('/').last;
    final responseTopic = '$TRANSACTIONS_RESPONSE_TOPIC/$replyTo';
    final handler = _handlers[responseTopic];
    if (handler == null) return;

    final response = _responseBuilder(request);
    final message = MqttPublishMessage()
      ..payload.message.addAll(utf8.encode(json.encode(response)));
    final received = MqttReceivedMessage<MqttMessage>(responseTopic, message);
    handler(_client, received);
  }

  @override
  Future<void> subscribe(String topic, {MessageHandler? handler}) async {
    _handlers[topic] = handler;
  }

  @override
  Future<void> unsubscribe(String topic) async {
    _handlers.remove(topic);
  }

  @override
  Future<dynamic> sendRequest(
    String method,
    dynamic params,
    String replyTo,
  ) async {
    final request =
        json.decode(
              json.encode(
                RequestPayload(method: method, params: params).toJson(),
              ),
            )
            as Map<String, dynamic>;
    publishedRequests.add(request);
    final response = _responseBuilder(request);
    if (response['status'] == RESPONSE_STATUS_ERROR) {
      if (response['message']?.toString().contains('record not found') ==
          true) {
        return 0;
      }
      throw Exception('error in response: ${response['message']}');
    }
    return response['data'];
  }
}
