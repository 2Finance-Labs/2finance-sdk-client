package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	ID         string         `yaml:"id"`
	Title      string         `yaml:"title"`
	Kind       string         `yaml:"kind"`
	Owner      string         `yaml:"owner"`
	Systems    []string       `yaml:"systems"`
	Goal       string         `yaml:"goal"`
	Requires   Requires       `yaml:"requires"`
	Actors     []Actor        `yaml:"actors"`
	Fixtures   map[string]any `yaml:"fixtures"`
	Steps      []Step         `yaml:"steps"`
	Assertions []Assertion    `yaml:"assertions"`
}

type Requires struct {
	Env []string `yaml:"env"`
}

type Actor struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`
}

type Step struct {
	ID     string         `yaml:"id"`
	Via    string         `yaml:"via"`
	Tool   string         `yaml:"tool"`
	Action string         `yaml:"action"`
	Input  map[string]any `yaml:"input"`
	Expect map[string]any `yaml:"expect"`
}

type Assertion struct {
	ID      string         `yaml:"id"`
	Subject string         `yaml:"subject"`
	Must    map[string]any `yaml:"must"`
}

func LoadSpec(path string) (Spec, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, err
	}
	var spec Spec
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		return Spec{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return spec, nil
}

func LoadSpecs(root string) ([]Spec, error) {
	paths := []string{}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".yaml", ".yml":
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Strings(paths)

	specs := make([]Spec, 0, len(paths))
	for _, path := range paths {
		spec, err := LoadSpec(path)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

func Validate(spec Spec) error {
	var problems []string
	required := map[string]string{
		"id":    spec.ID,
		"title": spec.Title,
		"kind":  spec.Kind,
		"owner": spec.Owner,
		"goal":  strings.TrimSpace(spec.Goal),
	}
	for field, value := range required {
		if strings.TrimSpace(value) == "" {
			problems = append(problems, field+" is required")
		}
	}
	if !validKind(spec.Kind) {
		problems = append(problems, "kind must be one of unit, contract, integration, e2e")
	}
	if len(spec.Systems) == 0 {
		problems = append(problems, "systems must not be empty")
	}
	if len(spec.Actors) == 0 {
		problems = append(problems, "actors must not be empty")
	}
	if len(spec.Steps) == 0 {
		problems = append(problems, "steps must not be empty")
	}
	if len(spec.Assertions) == 0 {
		problems = append(problems, "assertions must not be empty")
	}
	if spec.Kind == "e2e" && len(spec.Requires.Env) == 0 {
		problems = append(problems, "e2e specs must declare required env vars")
	}
	for _, env := range spec.Requires.Env {
		if env != strings.ToUpper(env) {
			problems = append(problems, "required env var must be uppercase: "+env)
		}
	}

	seenSteps := map[string]bool{}
	for i, step := range spec.Steps {
		prefix := fmt.Sprintf("steps[%d]", i)
		if strings.TrimSpace(step.ID) == "" {
			problems = append(problems, prefix+".id is required")
		}
		if seenSteps[step.ID] {
			problems = append(problems, "duplicate step id: "+step.ID)
		}
		seenSteps[step.ID] = true
		if strings.TrimSpace(step.Via) == "" {
			problems = append(problems, prefix+".via is required")
		}
		if !validVia(step.Via) {
			problems = append(problems, prefix+".via must be one of protocol, wallet_manager, client_2finance, mcp_http, network_vm")
		}
		if strings.TrimSpace(step.Action) == "" {
			problems = append(problems, prefix+".action is required")
		}
		if strings.TrimSpace(step.Via) == "mcp_http" && strings.TrimSpace(step.Tool) == "" {
			problems = append(problems, prefix+".tool is required for mcp_http")
		}
		if len(step.Expect) == 0 {
			problems = append(problems, prefix+".expect must not be empty")
		}
	}

	seenAssertions := map[string]bool{}
	for i, assertion := range spec.Assertions {
		prefix := fmt.Sprintf("assertions[%d]", i)
		if strings.TrimSpace(assertion.ID) == "" {
			problems = append(problems, prefix+".id is required")
		}
		if seenAssertions[assertion.ID] {
			problems = append(problems, "duplicate assertion id: "+assertion.ID)
		}
		seenAssertions[assertion.ID] = true
		if strings.TrimSpace(assertion.Subject) == "" {
			problems = append(problems, prefix+".subject is required")
		}
		if len(assertion.Must) == 0 {
			problems = append(problems, prefix+".must must not be empty")
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("%s: %s", spec.ID, strings.Join(problems, "; "))
	}
	return nil
}

func validKind(kind string) bool {
	switch kind {
	case "unit", "contract", "integration", "e2e":
		return true
	default:
		return false
	}
}

func validVia(via string) bool {
	switch via {
	case "protocol", "wallet_manager", "client_2finance", "mcp_http", "network_vm":
		return true
	default:
		return false
	}
}
