import 'package:test/test.dart';
import 'package:two_finance_blockchain/blockchain/contract/constants.dart';
import 'package:two_finance_blockchain/two_finance_blockchain.dart';

import '../../../helpers/fake_mqtt.dart';
import '../../../helpers/helpers.dart';

void main() {
  group('lifecycle payload normalization', () {
    test('keeps lifecycle fields and converts nested maps safely', () {
      final payload = normalizeLifecyclePayload({
        'address': '0xabc',
        'owner': '0xowner',
        'request_id': 'req-1',
        'provider': 'wise',
        'metadata': {'currency': 'USD'},
      });

      expect(payload['address'], '0xabc');
      expect(payload['owner'], '0xowner');
      expect(payload['request_id'], 'req-1');
      expect(payload['provider'], 'wise');
      expect(payload['metadata'], isA<Map<String, dynamic>>());
      expect(payload['metadata']['currency'], 'USD');
    });

    test('keeps Codexa quote-first sending fields', () {
      final payload = normalizeLifecyclePayload({
        'request_id': 'send-codexa-001',
        'provider': 'codexa',
        'source_currency': 'BRL',
        'target_currency': 'USD',
        'currency': 'USD',
        'amount': '1000.00',
        'beneficiary_name': 'John Smith',
        'details': {'provider': 'codexa'},
      });

      expect(payload['provider'], 'codexa');
      expect(payload['source_currency'], 'BRL');
      expect(payload['target_currency'], 'USD');
      expect(payload['details'], isA<Map<String, dynamic>>());
    });
  });

  group('lifecycle client payloads', () {
    test(
      '[spec:lifecycle.start_fx] [spec:lifecycle.advance_sending] start/advance lifecycle methods use Go method names and include address',
      () async {
        final signer = await validKeyPair();
        final lifecycleAddress = await validPublicKeyHex();
        final mqtt = FakeMqttClient();
        final client = TwoFinanceBlockchain(
          keyManager: KeyManager(),
          mqttClient: mqtt,
          chainID: 1,
        );
        await client.setPrivateKey(signer.privateKey);

        final calls = <Future<ContractOutput> Function()>[
          () => client.startFX(
            address: lifecycleAddress,
            data: const {'request_id': 'fx-001'},
          ),
          () => client.startSending(
            address: lifecycleAddress,
            data: const {
              'request_id': 'sending-codexa-001',
              'provider': 'codexa',
              'source_currency': 'BRL',
              'target_currency': 'USD',
              'currency': 'USD',
              'amount': '1000.00',
              'beneficiary_name': 'John Smith',
            },
          ),
          () => client.startReceiving(
            address: lifecycleAddress,
            data: const {
              'request_id': 'receiving-codexa-brl-001',
              'provider': 'codexa',
              'currency': 'BRL',
              'amount': '25.50',
              'details': {'payment_method': 'pix'},
            },
          ),
          () => client.advanceSending(
            address: lifecycleAddress,
            data: const {'request_id': 'sending-001', 'status': 'funded'},
          ),
        ];
        final expectedMethods = <String>[
          'start_FX',
          'start_Sending',
          'start_Receiving',
          'advance_Sending',
        ];
        final expectedData = <Map<String, dynamic>>[
          {'address': lifecycleAddress, 'request_id': 'fx-001'},
          {
            'address': lifecycleAddress,
            'request_id': 'sending-codexa-001',
            'provider': 'codexa',
            'source_currency': 'BRL',
            'target_currency': 'USD',
            'currency': 'USD',
            'amount': '1000.00',
            'beneficiary_name': 'John Smith',
          },
          {
            'address': lifecycleAddress,
            'request_id': 'receiving-codexa-brl-001',
            'provider': 'codexa',
            'currency': 'BRL',
            'amount': '25.50',
            'details': {'payment_method': 'pix'},
          },
          {
            'address': lifecycleAddress,
            'request_id': 'sending-001',
            'status': 'funded',
          },
        ];

        for (var i = 0; i < calls.length; i++) {
          await calls[i]();
          final request = mqtt.lastRequest;
          final params = Map<String, dynamic>.from(request['params'] as Map);
          final data = Map<String, dynamic>.from(params['data'] as Map);
          expect(request['method'], REQUEST_METHOD_SEND_TRANSACTION);
          expect(params['from'], signer.publicKey);
          expect(params['to'], lifecycleAddress);
          expect(params['method'], expectedMethods[i]);
          expect(data, expectedData[i]);
        }
      },
    );

    test(
      '[spec:lifecycle.get_sending] get lifecycle state uses request_id only in get_state payload',
      () async {
        final lifecycleAddress = await validPublicKeyHex();
        final mqtt = FakeMqttClient();
        final client = TwoFinanceBlockchain(
          keyManager: KeyManager(),
          mqttClient: mqtt,
          chainID: 1,
        );

        await client.getSending(
          address: lifecycleAddress,
          requestId: 'sending-001',
        );

        expect(mqtt.lastRequest['method'], REQUEST_METHOD_GET_STATE);
        expect(mqtt.lastRequest['params'], {
          'to': lifecycleAddress,
          'method': 'get_Sending',
          'data': {'request_id': 'sending-001'},
        });
      },
    );
  });
}
