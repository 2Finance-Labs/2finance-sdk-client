import 'dart:convert';
import 'dart:io';

Future<void> main() async {
  final baseUrl = Platform.environment['ORCHESTRATOR_URL'] ?? 'http://127.0.0.1:8000';
  final client = HttpClient();

  try {
    final uri = Uri.parse('$baseUrl/v1/mcphost/sessions');
    final request = await client.postUrl(uri);
    request.headers.contentType = ContentType.json;
    request.headers.set('X-Tenant-ID', 'local');
    request.headers.set('X-User-ID', 'dart-smoke');
    request.write(jsonEncode({
      'model': 'openai:gpt-4o-mini',
      'system_prompt': 'dart orchestrator smoke test',
    }));

    final response = await request.close();
    final body = await utf8.decodeStream(response);
    if (response.statusCode != HttpStatus.ok) {
      stderr.writeln('orchestrator status=${response.statusCode} body=$body');
      exitCode = 1;
      return;
    }

    final decoded = jsonDecode(body) as Map<String, dynamic>;
    final data = decoded['data'] as Map<String, dynamic>?;
    final sessionID = data?['session_id'];
    if (sessionID is! String || sessionID.isEmpty) {
      stderr.writeln('missing session_id body=$body');
      exitCode = 1;
      return;
    }

    stdout.writeln('connected session_id=$sessionID');
  } finally {
    client.close(force: true);
  }
}
