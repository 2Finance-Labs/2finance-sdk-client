# Feedback loop

Use several feedback layers instead of making every change wait for live e2e
tests.

## Fast local checks

Run before committing most changes:

```bash
make test-fast
```

This runs:

- `gofmt` check on fast-feedback packages;
- `go vet` on non-live packages;
- spec/harness validation;
- protocol and wallet tests.

There are existing files outside the new harness that are not `gofmt` clean.
Use the stricter check when the team is ready for a formatting-only cleanup:

```bash
make fmt-check-all
```

## Spec checks

Run when changing behavior docs or YAML specs:

```bash
make test-harness
```

The harness catches missing fields, duplicate ids, unknown drivers, missing MCP
tool names, and missing e2e environment declarations.

## Live lifecycle checks

Run when EMQX/network dependencies are available:

```bash
make test-e2e-lifecycle
```

Run the MCP-backed signing flow when the MCP server is available:

```bash
MCP_URL=http://127.0.0.1:8089/mcp make mcp-e2e-check
```

## CI policy

The GitHub workflow intentionally runs only fast feedback. Live e2e tests should
be triggered in an environment that owns MCP, EMQX, and the network VM. This
keeps pull-request feedback short and makes infrastructure failures easier to
separate from code failures.

## Future checkers worth adding

- `golangci-lint` once the current codebase is ready for a stricter baseline.
- A spec-to-driver coverage report that shows which specs are only validated and
  which are executed against real drivers.
- A contract test that compares MCP `tools/list` output with checked-in spec
  expectations.
- A cleanup checker for generated `tests/e2e/wallets/*.wallet` artifacts.
