part of 'two_finance_blockchain.dart';

Map<String, dynamic> normalizeLifecyclePayload(Map<String, dynamic> input) {
  return Map<String, dynamic>.from(
    _normalizeLifecycleValue(input) as Map<dynamic, dynamic>,
  );
}

dynamic _normalizeLifecycleValue(dynamic value) {
  if (value is Map) {
    final normalized = <String, dynamic>{};
    for (final entry in value.entries) {
      normalized[entry.key.toString()] = _normalizeLifecycleValue(entry.value);
    }
    return normalized;
  }

  if (value is List) {
    return value.map(_normalizeLifecycleValue).toList();
  }

  return value;
}

extension Lifecycle on TwoFinanceBlockchain {
  Future<ContractOutput> startLifecycle({
    required String address,
    required String method,
    required Map<String, dynamic> data,
  }) async {
    if (address.isEmpty) throw ArgumentError('lifecycle address not set');
    if (method.isEmpty) throw ArgumentError('lifecycle method not set');

    final from = publicKeyHex ?? '';
    if (from.isEmpty) throw ArgumentError('from address not set');

    KeyManager.validateEDDSAPublicKeyHex(from);
    KeyManager.validateEDDSAPublicKeyHex(address);

    final uuid7 = newUUID7();
    const version = 1;

    final payload = normalizeLifecyclePayload({'address': address, ...data});

    return signAndSendTransaction(
      chainID: _chainID,
      from: from,
      to: address,
      method: method,
      data: payload,
      version: version,
      uuid7: uuid7,
    );
  }

  Future<ContractOutput> advanceLifecycle({
    required String address,
    required String method,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(address: address, method: method, data: data);
  }

  Future<ContractOutput> failLifecycle({
    required String address,
    required String method,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(address: address, method: method, data: data);
  }

  Future<ContractOutput> getLifecycleState({
    required String address,
    required String requestId,
    required String method,
  }) async {
    if (address.isEmpty) throw ArgumentError('lifecycle address not set');
    if (requestId.isEmpty) throw ArgumentError('request_id not set');
    if (method.isEmpty) throw ArgumentError('lifecycle method not set');

    KeyManager.validateEDDSAPublicKeyHex(address);

    return getState(
      to: address,
      method: method,
      data: {'request_id': requestId},
    );
  }

  Future<ContractOutput> startFX({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(address: address, method: 'start_FX', data: data);
  }

  Future<ContractOutput> advanceFX({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return advanceLifecycle(address: address, method: 'advance_FX', data: data);
  }

  Future<ContractOutput> failFX({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return failLifecycle(address: address, method: 'fail_FX', data: data);
  }

  Future<ContractOutput> getFX({
    required String address,
    required String requestId,
  }) async {
    return getLifecycleState(
      address: address,
      requestId: requestId,
      method: 'get_FX',
    );
  }

  Future<ContractOutput> startOnboarding({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(
      address: address,
      method: 'start_Onboarding',
      data: data,
    );
  }

  Future<ContractOutput> advanceOnboarding({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return advanceLifecycle(
      address: address,
      method: 'advance_Onboarding',
      data: data,
    );
  }

  Future<ContractOutput> failOnboarding({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return failLifecycle(
      address: address,
      method: 'fail_Onboarding',
      data: data,
    );
  }

  Future<ContractOutput> getOnboarding({
    required String address,
    required String requestId,
  }) async {
    return getLifecycleState(
      address: address,
      requestId: requestId,
      method: 'get_Onboarding',
    );
  }

  Future<ContractOutput> startReceiving({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(
      address: address,
      method: 'start_Receiving',
      data: data,
    );
  }

  Future<ContractOutput> advanceReceiving({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return advanceLifecycle(
      address: address,
      method: 'advance_Receiving',
      data: data,
    );
  }

  Future<ContractOutput> failReceiving({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return failLifecycle(
      address: address,
      method: 'fail_Receiving',
      data: data,
    );
  }

  Future<ContractOutput> getReceiving({
    required String address,
    required String requestId,
  }) async {
    return getLifecycleState(
      address: address,
      requestId: requestId,
      method: 'get_Receiving',
    );
  }

  Future<ContractOutput> startSending({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(
      address: address,
      method: 'start_Sending',
      data: data,
    );
  }

  Future<ContractOutput> advanceSending({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return advanceLifecycle(
      address: address,
      method: 'advance_Sending',
      data: data,
    );
  }

  Future<ContractOutput> failSending({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return failLifecycle(address: address, method: 'fail_Sending', data: data);
  }

  Future<ContractOutput> getSending({
    required String address,
    required String requestId,
  }) async {
    return getLifecycleState(
      address: address,
      requestId: requestId,
      method: 'get_Sending',
    );
  }

  Future<ContractOutput> startMultiCurrency({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return startLifecycle(
      address: address,
      method: 'start_MultiCurrency',
      data: data,
    );
  }

  Future<ContractOutput> advanceMultiCurrency({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return advanceLifecycle(
      address: address,
      method: 'advance_MultiCurrency',
      data: data,
    );
  }

  Future<ContractOutput> failMultiCurrency({
    required String address,
    required Map<String, dynamic> data,
  }) async {
    return failLifecycle(
      address: address,
      method: 'fail_MultiCurrency',
      data: data,
    );
  }

  Future<ContractOutput> getMultiCurrency({
    required String address,
    required String requestId,
  }) async {
    return getLifecycleState(
      address: address,
      requestId: requestId,
      method: 'get_MultiCurrency',
    );
  }
}
