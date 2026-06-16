package harness

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestSpecsAreWellFormed(t *testing.T) {
	specs := loadRepositorySpecs(t)
	if len(specs) == 0 {
		t.Fatal("expected at least one spec")
	}

	ids := map[string]bool{}
	for _, spec := range specs {
		spec := spec
		t.Run(spec.ID, func(t *testing.T) {
			if ids[spec.ID] {
				t.Fatalf("duplicate spec id: %s", spec.ID)
			}
			ids[spec.ID] = true
			if err := Validate(spec); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCriticalSpecsExist(t *testing.T) {
	required := []string{
		"protocol.prepared-transaction.signing-contract",
		"protocol.signed-transaction.network-roundtrip",
		"wallet-manager.unlock-and-sign-prepared-transaction",
		"lifecycle.send.client-direct-start",
		"lifecycle.send.mcp-sign-submit",
		"mcp.finance-tools.registration-contract",
	}
	specs := loadRepositorySpecs(t)
	found := map[string]bool{}
	for _, spec := range specs {
		found[spec.ID] = true
	}
	for _, id := range required {
		if !found[id] {
			t.Fatalf("critical spec %s not found", id)
		}
	}
}

func TestSpecCatalogHasCoverageForCoreBoundaries(t *testing.T) {
	specs := loadRepositorySpecs(t)
	owners := map[string]bool{}
	systems := map[string]bool{}
	for _, spec := range specs {
		owners[spec.Owner] = true
		for _, system := range spec.Systems {
			systems[system] = true
		}
	}

	for _, owner := range []string{"protocol", "wallet_manager", "lifecycle", "mcp"} {
		if !owners[owner] {
			t.Fatalf("missing spec owner %s", owner)
		}
	}
	for _, system := range []string{"protocol", "wallet_manager", "client_2finance", "mcp-2finance", "2finance-network"} {
		if !systems[system] {
			t.Fatalf("missing spec coverage for system %s", system)
		}
	}
}

func TestPreparedTransactionBoundarySpec(t *testing.T) {
	spec := findSpec(t, "protocol.prepared-transaction.signing-contract")

	requireAssertion(t, spec, "unsigned-boundary", "prepared_transaction", map[string]string{
		"hash":      "empty",
		"signature": "empty",
	})
	requireAssertion(t, spec, "signed-boundary", "signed_transaction", map[string]string{
		"hash":      "non_empty",
		"signature": "non_empty",
	})
}

func TestMCPSendLifecycleSpecCoversConversationAndSignature(t *testing.T) {
	spec := findSpec(t, "lifecycle.send.mcp-sign-submit")

	requireStep(t, spec, "deploy-prepare", "mcp_http", "finance.contract.deploy.prepare")
	requireStep(t, spec, "deploy-sign", "wallet_manager", "")
	requireStep(t, spec, "deploy-submit", "mcp_http", "finance.transaction.submit_signed")
	requireStep(t, spec, "send-start-prepare", "mcp_http", "finance.send.start.prepare")
	requireStep(t, spec, "send-start-sign", "wallet_manager", "")
	requireStep(t, spec, "send-start-submit", "mcp_http", "finance.transaction.submit_signed")
	requireStep(t, spec, "send-get", "mcp_http", "finance.send.get")

	requireAssertion(t, spec, "mcp-never-signs", "prepared_transaction", map[string]string{
		"hash":      "empty",
		"signature": "empty",
	})
	requireAssertion(t, spec, "client-signs-before-submit", "signed_transaction", map[string]string{
		"hash":      "non_empty",
		"signature": "non_empty",
	})
}

func loadRepositorySpecs(t *testing.T) []Spec {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "specs"))
	specs, err := LoadSpecs(root)
	if err != nil {
		t.Fatalf("load specs: %v", err)
	}
	return specs
}

func findSpec(t *testing.T, id string) Spec {
	t.Helper()

	for _, spec := range loadRepositorySpecs(t) {
		if spec.ID == id {
			return spec
		}
	}
	t.Fatalf("spec %s not found", id)
	return Spec{}
}

func requireStep(t *testing.T, spec Spec, id, via, tool string) {
	t.Helper()

	for _, step := range spec.Steps {
		if step.ID != id {
			continue
		}
		if step.Via != via {
			t.Fatalf("step %s via = %s, want %s", id, step.Via, via)
		}
		if step.Tool != tool {
			t.Fatalf("step %s tool = %s, want %s", id, step.Tool, tool)
		}
		return
	}
	t.Fatalf("step %s not found", id)
}

func requireAssertion(t *testing.T, spec Spec, id, subject string, must map[string]string) {
	t.Helper()

	for _, assertion := range spec.Assertions {
		if assertion.ID != id {
			continue
		}
		if assertion.Subject != subject {
			t.Fatalf("assertion %s subject = %s, want %s", id, assertion.Subject, subject)
		}
		for key, want := range must {
			if got, _ := assertion.Must[key].(string); got != want {
				t.Fatalf("assertion %s must[%s] = %v, want %s", id, key, assertion.Must[key], want)
			}
		}
		return
	}
	t.Fatalf("assertion %s not found", id)
}
