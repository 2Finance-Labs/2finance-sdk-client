# Behavior specs

Specs in this directory describe expected behavior before test code decides how
to run it.

Each spec should answer:

- what user or system goal is being protected;
- which systems participate;
- which steps execute the behavior;
- which invariants must stay true.

Run the spec validator with:

```bash
cd go && go test ./tests/harness
```

Live execution belongs in focused Go tests under `go/tests/e2e` or future harness
drivers.

See also:

- `docs/spec-driven-development.md`
- `docs/spec-authoring.md`
- `docs/spec-index.md`
