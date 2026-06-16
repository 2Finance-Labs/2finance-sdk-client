# Spec-driven test harness

This repository has three different testing concerns that should stay connected
but not mixed in the same file:

- business behavior: what 2Finance promises to support;
- execution harness: how a promise is exercised through client, wallet, MCP, or
  the network VM;
- implementation tests: Go tests that keep the behavior executable.

The structure in `tests/specs` and `tests/harness` makes that separation
explicit. Specs are readable examples. The harness validates that examples stay
well-formed and gives Go tests a stable contract to execute.

## Why this shape

Martin Fowler describes specification by example as useful because examples are
often easier to produce than complete formal pre/post conditions, especially for
business users. He also warns that examples cannot be the only requirements
technique; they work best when backed by collaboration, domain language, and
other checks.

For this codebase, that means:

- keep examples close to the domain: wallet signs prepared transactions,
  MCP prepares unsigned transactions, network accepts signed submissions;
- keep e2e coverage small and valuable: broad tests are slower and harder to
  debug, so they should cover only critical workflows;
- use narrower harness checks for contract shape, required fields, and
  invariants before paying the cost of real network execution.

## Repository layout

```text
docs/spec-driven-harness.md   Research notes and local conventions.
docs/spec-driven-development.md Local SDD workflow for this repository.
docs/spec-authoring.md        YAML authoring guide.
docs/spec-index.md            Current spec catalog.
tests/specs/                  Human-readable behavior specs.
tests/specs/**.yaml           Individual executable examples.
tests/harness/                Go loader/validator for specs.
tests/e2e/                    Existing live integration/e2e tests.
```

## Spec lifecycle

1. Capture the behavior in `tests/specs`.
2. Link the behavior to a local principle or workflow in `docs`.
3. Validate the spec with `go test ./tests/harness`.
4. Add or update a focused Go execution test when the spec needs real IO.
5. Keep broad e2e tests opt-in when they need MCP, EMQX, or the network VM.

## Local conventions

- `kind: unit` means no external process should be needed.
- `kind: contract` checks message shape or compatibility.
- `kind: integration` may use local adapters but should avoid full network
  lifecycle cost.
- `kind: e2e` may call MCP, EMQX, or the network VM and should declare required
  environment variables.
- `steps[].via` names the driver that should execute the action:
  `client_2finance`, `wallet_manager`, `mcp_http`, or `network_vm`.
- Specs should use stable domain words rather than Go helper names where
  possible.

## Current spec families

- `protocol.*`: transaction envelopes and network compatibility.
- `wallet-manager.*`: unlock, local signing, and key-safety boundaries.
- `lifecycle.*`: direct client and MCP-backed lifecycle execution.
- `mcp.*`: MCP tool registry and tool payload compatibility.

See `docs/spec-index.md` for the current list.

## References

- Martin Fowler, Specification By Example:
  https://martinfowler.com/bliki/SpecificationByExample.html
- Martin Fowler, Test Pyramid:
  https://martinfowler.com/bliki/TestPyramid.html
- Thoughtworks/Martin Fowler site, Consumer-Driven Contracts:
  https://martinfowler.com/articles/consumerDrivenContracts.html
- Google Testing Blog, Just Say No to More End-to-End Tests:
  https://testing.googleblog.com/2015/04/just-say-no-to-more-end-to-end-tests.html
- Cucumber BDD docs:
  https://cucumber.io/docs/bdd/
- GitHub Spec Kit:
  https://github.com/github/spec-kit
