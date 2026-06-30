# Spec authoring guide

Specs live under `go/tests/specs/**/*.yaml`.

## Required fields

```yaml
id: namespace.short-description
title: Human readable title
kind: unit | contract | integration | e2e
owner: team-or-module
systems:
  - protocol
goal: >
  Why this behavior matters.
actors:
  - name: actor_name
    role: what the actor does
steps:
  - id: step-id
    via: protocol | wallet_manager | client_2finance | mcp_http | network_vm
    tool: only required when via is mcp_http
    action: stable_domain_action
    expect:
      observable.path: expected_value
assertions:
  - id: assertion-id
    subject: domain_object
    must:
      field: invariant
```

## Optional fields

```yaml
requires:
  env:
    - MCP_E2E
fixtures:
  amount: "25.50"
  currency: USD
```

## Naming

- Use lowercase dot-separated spec ids.
- Use lowercase kebab-case step and assertion ids.
- Use stable domain action names, not Go helper names.
- Use `non_empty`, `empty`, `present`, and `ok` for abstract expectations when
  exact values are generated at runtime.

## Choosing `kind`

- `unit`: pure function or local state only.
- `contract`: payload shape, signer boundary, or MCP/client compatibility.
- `integration`: local adapter or multiple in-process modules.
- `e2e`: live MCP, EMQX, or 2finance-network VM.

## Review checklist

- Does the spec state a product or protocol goal?
- Can a reader understand the behavior without reading Go code?
- Are all external dependencies declared?
- Is signing authority explicit?
- Are assertions observable by a test or harness driver?
- Is broad e2e coverage justified, or would a contract spec be enough?
