import 'dart:convert';
import 'dart:io';

import 'package:test/test.dart';

void main() {
  group('spec files', () {
    final specFiles =
        Directory('specs')
            .listSync(recursive: true)
            .whereType<File>()
            .where((file) => file.path.endsWith('.spec.json'))
            .toList()
          ..sort((a, b) => a.path.compareTo(b.path));

    test('exist', () {
      expect(specFiles, isNotEmpty);
    });

    for (final file in specFiles) {
      test('${file.path} follows the spec schema', () {
        final decoded = json.decode(file.readAsStringSync());
        expect(decoded, isA<Map<String, dynamic>>());

        final spec = Map<String, dynamic>.from(decoded as Map);
        _expectNonEmptyString(spec, 'id');
        _expectNonEmptyString(spec, 'title');
        _expectNonEmptyString(spec, 'owner');
        _expectStatus(spec['status']);
        _expectNonEmptyStringList(spec, 'layers');

        final cases = spec['cases'];
        expect(cases, isA<List<dynamic>>());
        expect(cases as List<dynamic>, isNotEmpty);

        final caseIds = <String>{};
        for (final rawCase in cases) {
          expect(rawCase, isA<Map<dynamic, dynamic>>());
          final specCase = Map<String, dynamic>.from(rawCase as Map);

          _expectNonEmptyString(specCase, 'id');
          _expectNonEmptyString(specCase, 'title');
          _expectNonEmptyStringList(specCase, 'tags');
          _expectMap(specCase, 'given');
          _expectMap(specCase, 'when');
          _expectMap(specCase, 'then');

          final caseId = specCase['id'] as String;
          expect(
            caseIds.add(caseId),
            isTrue,
            reason: 'duplicate case id "$caseId" in ${file.path}',
          );

          final when = Map<String, dynamic>.from(specCase['when'] as Map);
          _expectNonEmptyString(when, 'client_call');

          final then = Map<String, dynamic>.from(specCase['then'] as Map);
          _expectNonEmptyString(then, 'request_method');
          _expectKnownLayerValues(spec['layers'] as List<dynamic>);
          _expectKnownRequestMethod(then['request_method'] as String);
          _expectObservableExpectation(then);
          _expectContractUnitRefs(spec, specCase, file.path);
        }
      });
    }

    test('case ids are globally unique', () {
      final ids = <String, String>{};
      for (final file in specFiles) {
        final spec =
            json.decode(file.readAsStringSync()) as Map<String, dynamic>;
        final cases = spec['cases'] as List<dynamic>;
        for (final rawCase in cases) {
          final specCase = rawCase as Map<String, dynamic>;
          final id = specCase['id'] as String;
          expect(
            ids.containsKey(id),
            isFalse,
            reason: 'case id "$id" appears in ${ids[id]} and ${file.path}',
          );
          ids[id] = file.path;
        }
      }
    });
  });
}

void _expectNonEmptyString(Map<String, dynamic> object, String key) {
  expect(object, contains(key));
  expect(object[key], isA<String>());
  expect((object[key] as String).trim(), isNotEmpty);
}

void _expectNonEmptyStringList(Map<String, dynamic> object, String key) {
  expect(object, contains(key));
  expect(object[key], isA<List<dynamic>>());
  final values = object[key] as List<dynamic>;
  expect(values, isNotEmpty);
  for (final value in values) {
    expect(value, isA<String>());
    expect((value as String).trim(), isNotEmpty);
  }
}

void _expectMap(Map<String, dynamic> object, String key) {
  expect(object, contains(key));
  expect(object[key], isA<Map<dynamic, dynamic>>());
}

void _expectStatus(dynamic status) {
  expect(status, isA<String>());
  expect(status, isIn(<String>['draft', 'active', 'deprecated']));
}

void _expectKnownLayerValues(List<dynamic> layers) {
  for (final layer in layers) {
    expect(
      layer,
      isIn(<String>['harness-small', 'contract-unit', 'e2e', 'manual']),
    );
  }
}

void _expectKnownRequestMethod(String requestMethod) {
  expect(
    requestMethod,
    isIn(<String>[
      'send_transaction',
      'get_state',
      'get_logs',
      'get_transactions',
      'get_blocks',
    ]),
  );
}

void _expectObservableExpectation(Map<String, dynamic> then) {
  final hasExpectation =
      then.containsKey('transaction') ||
      then.containsKey('state_query') ||
      then.containsKey('params') ||
      then.containsKey('output');

  expect(
    hasExpectation,
    isTrue,
    reason:
        'then must describe an observable transaction, state_query, params, '
        'or output expectation',
  );
}

void _expectContractUnitRefs(
  Map<String, dynamic> spec,
  Map<String, dynamic> specCase,
  String specPath,
) {
  final layers = spec['layers'] as List<dynamic>;
  if (!layers.contains('contract-unit')) return;

  _expectNonEmptyStringList(specCase, 'test_refs');

  final caseId = specCase['id'] as String;
  final refs = specCase['test_refs'] as List<dynamic>;
  for (final ref in refs) {
    final path = ref as String;
    final file = File(path);
    expect(file.existsSync(), isTrue, reason: '$specPath references $path');

    final contents = file.readAsStringSync();
    expect(
      contents,
      contains('[spec:$caseId]'),
      reason: '$path must include a [spec:$caseId] marker',
    );
  }
}
