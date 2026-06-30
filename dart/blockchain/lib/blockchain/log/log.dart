import 'dart:convert';

/// Interface-like abstract class
abstract class ILog {
  String getLogType();
  int getLogIndex();
  String getTransactionHash();
  String getEvent();
  String getContractVersion();
  String getContractAddress();
  void validateLog();
}

/// Concrete implementation
class Log implements ILog {
  final String logType;
  final int logIndex;
  final String transactionHash;
  final String event;
  final String contractVersion;
  final String contractAddress;

  Log({
    required this.logType,
    required this.logIndex,
    required this.transactionHash,
    required this.event,
    required this.contractVersion,
    required this.contractAddress,
  });

  factory Log.fromJson(Map<String, dynamic> json) {
    final dynamic ev = json['event'];
    if (ev == null || ev is! String || ev.isEmpty) {
      throw Exception('log event is required');
    }
    return Log(
      logType: json['log_type'] as String,
      logIndex: (json['log_index'] as num).toInt(),
      transactionHash: json['transaction_hash'] as String,
      event: json['event'] as String, // base64
      contractVersion: json['contract_version'] as String,
      contractAddress: json['contract_address'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'log_type': logType,
      'log_index': logIndex,
      'transaction_hash': transactionHash,
      'event': event,
      'contract_version': contractVersion,
      'contract_address': contractAddress,
    };
  }

  @override
  String getLogType() => logType;

  @override
  int getLogIndex() => logIndex;

  @override
  String getTransactionHash() => transactionHash;

  @override
  String getEvent() => event;

  @override
  String getContractVersion() => contractVersion;

  @override
  String getContractAddress() => contractAddress;

  @override
  void validateLog() {
    if (logType.isEmpty) {
      throw Exception('log_type cannot be empty');
    }
    if (logIndex == 0) {
      throw Exception('log_index cannot be zero');
    }
    if (transactionHash.length != 64) {
      throw Exception('transaction_hash must be 64 characters');
    }
    if (event.isEmpty) {
      throw Exception('event (state_data) cannot be empty');
    }
    if (contractAddress.isEmpty) {
      throw Exception('contract_address cannot be empty');
    }
    if (contractVersion.isEmpty) {
      throw Exception('contract_version cannot be empty');
    }
  }
}

/// Factory function equivalent to `NewLog(...)`
Log newLog({
  required String logType,
  required int logIndex,
  required String transactionHash,
  required String event,
  required String contractVersion,
  required String contractAddress,
}) {
  return Log(
    logType: logType,
    logIndex: logIndex,
    transactionHash: transactionHash,
    event: event,
    contractVersion: contractVersion,
    contractAddress: contractAddress,
  );
}
