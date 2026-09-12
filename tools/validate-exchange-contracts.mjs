import fs from 'node:fs';
import path from 'node:path';

const root = path.resolve(new URL('..', import.meta.url).pathname);
const readJSON = (file) => JSON.parse(fs.readFileSync(path.join(root, file), 'utf8'));
const fail = (message) => { throw new Error(`exchange contracts: ${message}`); };

const eventSchema = readJSON('contracts/schemas/matchengine-ws-event.v1.json');
const responseSchema = readJSON('contracts/schemas/matchengine-order-response.v2.json');
const golden = readJSON('contracts/examples/exchange-websocket-golden.v1.json');

if (eventSchema.$schema !== 'https://json-schema.org/draft/2020-12/schema' || !eventSchema.$id) fail('event schema metadata');
if (responseSchema.$schema !== 'https://json-schema.org/draft/2020-12/schema' || !responseSchema.$id) fail('response schema metadata');
if (golden.schema !== '2finance.exchange_websocket_golden.v1' || golden.frames.length < 5 || golden.command_responses.length !== 2) fail('golden catalog is incomplete');

const eventTypes = new Set(eventSchema.properties.event_type.enum);
const streams = new Set(eventSchema.properties.stream.enum);
const lastByStream = new Map();
for (const frame of golden.frames) {
  for (const field of eventSchema.required) if (!(field in frame)) fail(`frame missing ${field}`);
  if (frame.schema !== 'matchengine.ws_event.v1' || frame.version !== 1) fail('event version mismatch');
  if (frame.engine_id !== golden.engine_id || frame.symbol_id !== golden.symbol_id || frame.route_epoch !== golden.route_epoch) fail('event route identity mismatch');
  if (!streams.has(frame.stream) || !eventTypes.has(frame.event_type) || typeof frame.payload !== 'object' || Array.isArray(frame.payload)) fail('event enum or payload mismatch');
  if (frame.stream === 'PRIVATE' && (!frame.account_id || !Number.isInteger(frame.wallet_id))) fail('private event missing account/wallet identity');
  const previous = lastByStream.get(frame.stream);
  if (previous !== undefined && frame.previous_sequence !== previous) fail(`${frame.stream} sequence is not contiguous`);
  lastByStream.set(frame.stream, frame.sequence);
}

const responseTypes = new Set(responseSchema.properties.message_type.enum);
for (const response of golden.command_responses) {
  for (const field of responseSchema.required) if (!(field in response)) fail(`command response missing ${field}`);
  if (response.schema !== 'matchengine.order_response.v2' || response.version !== 2 || !responseTypes.has(response.message_type)) fail('command response version/type mismatch');
  if (response.idempotency_key.length < 16) fail('short idempotency key');
  if (response.message_type === 'REJECT' && (!response.reason_code || !response.message)) fail('reject lacks actionable reason');
}
if (!golden.invalid_frames.some((item) => item.expected === 'reject-and-resync') || !golden.invalid_frames.some((item) => item.expected === 'recovery-from-44') || !golden.invalid_frames.some((item) => item.expected === 'schema-reject')) fail('negative golden cases are incomplete');

console.log(`exchange contracts passed: ${golden.frames.length} events, ${golden.command_responses.length} responses, ${golden.invalid_frames.length} negative cases`);
