// test/key_manager_test.dart
import 'dart:typed_data';

import 'package:test/test.dart';

// Adjust this import to your project path:
// e.g. import 'package:two_finance_blockchain/blockchain/encryption/key_manager.dart';
import 'package:two_finance_blockchain/blockchain/keys/keys.dart';

void main() {
  group('KeyManager.bytesToHex / hexToBytes', () {
    test('bytesToHex converts bytes to lowercase hex with leading zeros', () {
      final bytes = <int>[0, 1, 15, 16, 255];
      final hex = KeyManager.bytesToHex(bytes);

      expect(hex, '00010f10ff');
    });

    test('hexToBytes converts hex to Uint8List', () {
      final hex = '00010f10ff';
      final bytes = KeyManager.hexToBytes(hex);

      expect(bytes, Uint8List.fromList([0, 1, 15, 16, 255]));
    });

    test('hexToBytes throws FormatException when hex length is odd', () {
      expect(
        () => KeyManager.hexToBytes('abc'),
        throwsA(isA<FormatException>()),
      );
    });

    test('bytesToHex -> hexToBytes roundtrip returns original bytes', () {
      final original = Uint8List.fromList(List<int>.generate(32, (i) => i));
      final hex = KeyManager.bytesToHex(original);
      final back = KeyManager.hexToBytes(hex);

      expect(back, original);
      expect(hex.length, 64); // 32 bytes => 64 hex chars
    });
  });

  group('KeyManager.generateKeyEd25519', () {
    test('generates a KeyPair2Finance with valid hex lengths', () async {
      final km = KeyManager();
      final pair = await km.generateKeyEd25519();

      expect(pair, isA<KeyPair2Finance>());
      expect(pair.publicKey, isNotEmpty);
      expect(pair.privateKey, isNotEmpty);

      // Ed25519 public key is 32 bytes => 64 hex chars
      expect(pair.publicKey.length, 64);

      // cryptography's extractPrivateKeyBytes() returns 64 bytes for Ed25519
      expect(pair.privateKey.length, 64);

      // Ensure both are valid hex and convertible back to bytes
      final pubBytes = KeyManager.hexToBytes(pair.publicKey);
      final privBytes = KeyManager.hexToBytes(pair.privateKey);
      expect(pubBytes.length, 32);
      expect(privBytes.length, 32);
    });

    test('generated public key passes validateEDDSAPublicKeyHex', () async {
      final km = KeyManager();
      final pair = await km.generateKeyEd25519();

      expect(
        () => KeyManager.validateEDDSAPublicKeyHex(pair.publicKey),
        returnsNormally,
      );
    });

    test('two generated keypairs are different (very likely)', () async {
      final km = KeyManager();

      final p1 = await km.generateKeyEd25519();
      final p2 = await km.generateKeyEd25519();

      expect(p1.publicKey, isNot(p2.publicKey));
      expect(p1.privateKey, isNot(p2.privateKey));
    });
  });

  group('KeyManager.validateEDDSAPublicKeyHex', () {
    test('throws if public key is empty', () {
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex(''),
        throwsA(isA<FormatException>()),
      );
    });

    test('throws if public key hex length is not 64 (not 32 bytes)', () {
      // 2 bytes => 4 hex chars
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex('aabb'),
        throwsA(isA<FormatException>()),
      );

      // 31 bytes => 62 hex chars
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex('00' * 31),
        throwsA(isA<FormatException>()),
      );

      // 33 bytes => 66 hex chars
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex('01' * 33),
        throwsA(isA<FormatException>()),
      );
    });

    test('throws if public key is all zeros', () {
      final allZero64Hex = '00' * 32; // 32 bytes
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex(allZero64Hex),
        throwsA(isA<FormatException>()),
      );
    });

    test('throws if public key hex is not parseable', () {
      // hexToBytes will throw FormatException from int.parse
      expect(
        () => KeyManager.validateEDDSAPublicKeyHex('zz' * 32),
        throwsA(isA<FormatException>()),
      );
    });
  });

  group('KeyPair2Finance.toString', () {
    test('includes both keys', () {
      final kp = KeyPair2Finance(publicKey: 'aa', privateKey: 'bb');
      final s = kp.toString();

      expect(s, contains('Public Key: aa'));
      expect(s, contains('Private Key: bb'));
    });
  });
}
