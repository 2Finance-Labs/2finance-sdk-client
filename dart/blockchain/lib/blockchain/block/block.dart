class Block {
  final int number;
  final DateTime? timestamp;
  final String txRoot;
  final String logRoot;
  final String stateRoot;
  final String timestampRoot;
  final int txCount;
  final int logCount;
  final int stateSnapshotCount;
  final String hash;
  final String parentHash;
  final DateTime? createdAt;

  Block({
    required this.number,
    this.timestamp,
    required this.txRoot,
    required this.logRoot,
    required this.stateRoot,
    required this.timestampRoot,
    required this.txCount,
    required this.logCount,
    required this.stateSnapshotCount,
    required this.hash,
    required this.parentHash,
    this.createdAt,
  });

  factory Block.fromJson(Map<String, dynamic> json) {
    return Block(
      number: _asInt(json['number']),
      timestamp: _asDateTime(json['timestamp']),
      txRoot: json['tx_root'] as String? ?? '',
      logRoot: json['log_root'] as String? ?? '',
      stateRoot: json['state_root'] as String? ?? '',
      timestampRoot: json['timestamp_root'] as String? ?? '',
      txCount: _asInt(json['tx_count']),
      logCount: _asInt(json['log_count']),
      stateSnapshotCount: _asInt(json['state_snapshot_count']),
      hash: json['hash'] as String? ?? '',
      parentHash: json['parent_hash'] as String? ?? '',
      createdAt: _asDateTime(json['created_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'number': number,
      if (timestamp != null) 'timestamp': timestamp!.toUtc().toIso8601String(),
      'tx_root': txRoot,
      'log_root': logRoot,
      'state_root': stateRoot,
      'timestamp_root': timestampRoot,
      'tx_count': txCount,
      'log_count': logCount,
      'state_snapshot_count': stateSnapshotCount,
      'hash': hash,
      'parent_hash': parentHash,
      if (createdAt != null) 'created_at': createdAt!.toUtc().toIso8601String(),
    };
  }
}

int _asInt(dynamic value) {
  if (value == null) return 0;
  if (value is int) return value;
  if (value is num) return value.toInt();
  if (value is String) return int.tryParse(value) ?? 0;
  throw Exception('invalid integer value: $value');
}

DateTime? _asDateTime(dynamic value) {
  if (value == null) return null;
  if (value is DateTime) return value;
  if (value is int) return DateTime.fromMillisecondsSinceEpoch(value);
  if (value is num) return DateTime.fromMillisecondsSinceEpoch(value.toInt());
  if (value is String && value.isNotEmpty) return DateTime.tryParse(value);
  return null;
}
