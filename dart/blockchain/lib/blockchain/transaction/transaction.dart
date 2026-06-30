import 'dart:convert';
import 'dart:typed_data';
import 'package:cryptography/cryptography.dart';
import 'package:web3dart/crypto.dart';
import 'package:two_finance_blockchain/blockchain/keys/keys.dart';
import 'package:two_finance_blockchain/blockchain/types/types.dart';
import 'package:two_finance_blockchain/blockchain/utils/json.dart';
import 'package:two_finance_blockchain/blockchain/utils/uuid.dart';

abstract class ITransaction {
  Future<void> validateUnsignedTransaction();
  Future<void> validateTransaction();
  Future<void> validateHash();
  Future<String> calculateHash();
  Transaction get();
}

class Transaction implements ITransaction {
  int chainID;
  String from;
  String to;
  String method;
  JsonMessage data;
  int version;
  String uuid7;
  String hash;
  String signature;

  Transaction({
    required this.chainID,
    required this.from,
    required this.to,
    required this.method,
    required this.data,
    required this.version,
    required this.uuid7,
    this.hash = '',
    this.signature = '',
  });

  @override
  String toString() {
    return '''
  Transaction(
    chainID: $chainID,
    from: $from,
    to: $to,
    method: $method,
    data: $data,
    version: $version,
    uuid7: $uuid7,
    hash: $hash,
    signature: $signature,
  )
  ''';
  }

  static Transaction create({
    required int chainID,
    required String from,
    required String to,
    required String method,
    required JsonMessage data,
    required int version,
    required String uuid7,
  }) {
    return Transaction(
      chainID: chainID,
      from: from,
      to: to,
      method: method,
      data: data,
      version: version,
      uuid7: uuid7,
    );
  }

  @override
  Future<void> validateUnsignedTransaction() async {
    if (data.isEmpty) {
      throw Exception("data cannot be empty");
    }
  }

  @override
  Future<void> validateTransaction() async {
    if (chainID <= 0) {
      throw Exception("chain ID must be greater than zero");
    }
    if (chainID > 2) {
      throw Exception("unsupported chain ID");
    }
    if (from.isEmpty) {
      throw Exception("sender address is required");
    }
    if (to.isEmpty) {
      throw Exception("recipient address is required");
    }
    if (from == to) {
      throw Exception("sender and recipient cannot be the same");
    }
    if (method.isEmpty) {
      throw Exception("method is required");
    }
    if (data.isEmpty) {
      throw Exception("data cannot be empty");
    }
    if (version == 0) {
      throw Exception("version must be greater than zero");
    }

    // Sender pubkey validation
    try {
      KeyManager.validateEDDSAPublicKeyHex(from);
    } catch (e) {
      throw Exception("invalid sender public key: $e");
    }

    // Recipient pubkey validation (skip deploy address)
    if (to.isNotEmpty && to != DEPLOY_CONTRACT_ADDRESS) {
      try {
        KeyManager.validateEDDSAPublicKeyHex(to);
      } catch (e) {
        throw Exception("invalid recipient public key: $e");
      }
    }

    // UUIDv7 validation
    try {
      validateUUID7(uuid7);
    } catch (e) {
      throw Exception("invalid UUIDv7: $e");
    }

    if (hash.length != 64) {
      throw Exception("hash must be 64 characters long");
    }
    if (signature.isEmpty) {
      throw Exception("signature cannot be empty");
    }
    if (signature.length != 128) {
      throw Exception("signature must be 128 characters long");
    }

    // Hash validation
    try {
      await validateHash();
    } catch (e) {
      throw Exception("transaction hash validation failed: $e");
    }
  }

  @override
  Future<void> validateHash() async {
    final computed = await calculateHash();
    if (computed != hash) {
      throw Exception("invalid hash: expected $computed, got $hash");
    }
  }

  @override
  Future<String> calculateHash() async {
    final temp = toJson();
    temp['hash'] = '';
    temp['signature'] = '';

    // ✅ canonicalize `data` semantically (Map order/whitespace becomes irrelevant)
    final d = temp['data'];
    if (d is Uint8List) {
      final decoded = jsonDecode(utf8.decode(d)); // parse JSON
      temp['data'] = decoded; // store as semantic object
    }

    final canonical = canonicalJsonEncode(temp);
    final encoded = Uint8List.fromList(utf8.encode(canonical));
    final sum = keccak256(encoded);
    return KeyManager.bytesToHex(sum);
  }

  @override
  Transaction get() => this;

  Map<String, dynamic> toJson() => {
    'chain_id': chainID,
    'from': from,
    'to': to,
    'method': method,
    'data': data,
    'version': version,
    'uuid7': uuid7,
    'hash': hash,
    'signature': signature,
  };

  static Transaction fromJson(Map<String, dynamic> json) {
    final rawData = json['data'];
    if (rawData == null) {
      throw Exception("transaction data cannot be null");
    }

    return Transaction(
      chainID: json['chain_id'],
      from: json['from'],
      to: json['to'],
      method: json['method'],
      data: rawData,
      version: json['version'],
      uuid7: json['uuid7'],
      hash: json['hash'],
      signature: json['signature'],
    );
  }
}

/// Signs a transaction using a hex-encoded Ed25519 private key
Future<Transaction> signTransaction(
  String privateKeyHex,
  Transaction tx,
) async {
  await tx.validateUnsignedTransaction();

  final privateKeyBytes = KeyManager.hexToBytes(privateKeyHex);
  if (privateKeyBytes.length < 32) {
    throw Exception("private key must be at least 32 bytes (64 hex chars)");
  }

  final seed = privateKeyBytes.sublist(0, 32);
  final algorithm = Ed25519();
  final keyPair = await algorithm.newKeyPairFromSeed(seed);

  final txHash = await tx.calculateHash();
  final hashBytes = KeyManager.hexToBytes(txHash);
  final signature = await algorithm.sign(hashBytes, keyPair: keyPair);

  tx.hash = txHash;
  tx.signature = KeyManager.bytesToHex(signature.bytes);
  return tx;
}
