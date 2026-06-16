# Spec-driven development for go-client-2finance

Spec-driven development treats specifications as the durable product artifact,
not as disposable notes created before coding. For this repository, specs should
describe the financial protocol boundary first, then the implementation should
prove that boundary through focused harness tests and a small number of live
e2e tests.

## Local principles

1. Specs describe `what` and `why` before `how`.
2. Plans describe the chosen implementation path before code is changed.
3. Tasks are small enough to map to a test, checker, or file change.
4. Live e2e tests exist only for behavior that cannot be trusted through a
   smaller contract or integration test.
5. The MCP/client signing boundary is a product invariant: MCP prepares,
   wallet/client signs, network receives signed transactions.

## Recommended workflow

### 1. Constitution

Update this document when the team agrees on a new durable rule. Examples:

- prepared transactions must be unsigned;
- private keys never leave `wallet_manager`;
- MCP tools cannot claim they signed a transaction;
- e2e tests must declare the external services they require.

### 2. Specify

Add or update a YAML file under `tests/specs`.

Good specs contain:

- a domain goal;
- systems involved;
- actors;
- concrete steps;
- observable expectations;
- assertions that should survive implementation changes.

### 3. Plan

Before implementing, decide which driver should prove each step:

- `protocol` for pure data-shape behavior;
- `wallet_manager` for signing and key lifecycle;
- `client_2finance` for client API behavior;
- `mcp_http` for MCP tool boundaries;
- `network_vm` for persisted blockchain/VM state.

### 4. Tasks

Break the plan into a short checklist:

- spec updated;
- harness validation updated;
- unit or contract test added;
- e2e added or intentionally skipped;
- docs updated.

### 5. Implement

Code follows the spec. If code reveals a missing rule, update the spec first or
in the same change.

### 6. Validate

Run the fastest applicable feedback first:

```bash
make test-fast
```

Then run live checks only when dependencies are available:

```bash
make test-e2e-lifecycle
MCP_URL=http://127.0.0.1:8089/mcp make mcp-e2e-check
```

## Traceability rules

- Every critical e2e should point back to a spec id in its test name, comment, or
  future harness metadata.
- Every MCP lifecycle tool used by the client should appear in a spec.
- Every spec with `kind: e2e` must declare required environment variables.
- Every spec assertion should be observable from either structured output,
  wallet state, network state, or explicit logs.

## Relationship to references

GitHub Spec Kit frames SDD as a workflow where specs are first-class and followed
by plans, tasks, and implementation. Fowler's Specification by Example reminds
us to keep examples concrete and collaborative, while Cucumber/BDD emphasizes a
shared vocabulary for behavior. OpenAPI's contract-first approach is a useful
analogy for MCP and protocol payloads: agree on the shape first, then implement
and test against it.
