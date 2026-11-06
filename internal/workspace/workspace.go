// package workspace loads workspaces from directories:
//
//	policy.yaml       the policy (required)
//	constraints.yaml  system constraints (optional)
//	contracts.yaml    requirements and their contracts (optional)
//	worlds/*.yaml     saved worlds, named after their file (optional)
package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const (
	PolicyFile      = "policy.yaml"
	ConstraintsFile = "constraints.yaml"
	ContractsFile   = "contracts.yaml"
	WorldsDir       = "worlds"
)

// Load reads the workspace in dir. Files are parsed but not validated against the policy:
// that happens when the workspace is used.
func Load(dir string) (*types.Workspace, error) {
	policy, err := os.ReadFile(filepath.Join(dir, PolicyFile))
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	model := &types.SystemModel{Policy: string(policy)}

	if data, ok, err := readOptional(filepath.Join(dir, ConstraintsFile)); err != nil {
		return nil, err
	} else if ok {
		if model.Constraints, err = pc.UnmarshalAssumptions(data); err != nil {
			return nil, fmt.Errorf("%v: %w", ConstraintsFile, err)
		}
	}
	if data, ok, err := readOptional(filepath.Join(dir, ContractsFile)); err != nil {
		return nil, err
	} else if ok {
		if model.Requirements, err = pc.UnmarshalRequirements(data); err != nil {
			return nil, fmt.Errorf("%v: %w", ContractsFile, err)
		}
	}

	paths, err := filepath.Glob(filepath.Join(dir, WorldsDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	sort.Strings(paths)
	ws := &types.Workspace{SystemModel: model}
	for _, path := range paths {
		world, err := loadWorld(path)
		if err != nil {
			return nil, err
		}
		ws.Worlds = append(ws.Worlds, world)
	}
	return ws, nil
}

// MarshalWorld renders a world in the worlds/*.yaml format. The name is left out: it's the file name.
func MarshalWorld(world *types.World) ([]byte, error) {
	return yaml.Marshal(worldFile{Relationships: world.Relationships, Assertions: world.Assertions})
}

// worldFile is the format of worlds/*.yaml
type worldFile struct {
	Relationships string `json:"relationships"`
	Assertions    string `json:"assertions,omitempty"`
}

func loadWorld(path string) (*types.World, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	var f worldFile
	if err := yaml.UnmarshalStrict(data, &f); err != nil {
		return nil, fmt.Errorf("%v: %w", filepath.Join(WorldsDir, filepath.Base(path)), err)
	}
	return &types.World{
		Name:          strings.TrimSuffix(filepath.Base(path), ".yaml"),
		Relationships: f.Relationships,
		Assertions:    f.Assertions,
	}, nil
}

func readOptional(path string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("loading workspace: %w", err)
	}
	return string(data), true, nil
}

// Files are the contents of a workspace's files as written, for clients which parse them themselves
// (eg. the playground UI, which reports errors by line). Missing optional files are empty.
type Files struct {
	Dir         string      `json:"dir"`
	Policy      string      `json:"policy"`
	Constraints string      `json:"constraints"`
	Contracts   string      `json:"contracts"`
	Worlds      []WorldFile `json:"worlds"`
	// Version changes whenever any file does
	Version string `json:"version"`
}

// WorldFile is a worlds/<name>.yaml file
type WorldFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ReadFiles reads the files of the workspace in dir without parsing them. Only policy.yaml is required.
func ReadFiles(dir string) (*Files, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	policy, err := os.ReadFile(filepath.Join(abs, PolicyFile))
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	files := &Files{Dir: abs, Policy: string(policy), Worlds: []WorldFile{}}
	for name, dst := range map[string]*string{ConstraintsFile: &files.Constraints, ContractsFile: &files.Contracts} {
		if *dst, _, err = readOptional(filepath.Join(abs, name)); err != nil {
			return nil, err
		}
	}
	paths, err := filepath.Glob(filepath.Join(abs, WorldsDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("loading workspace: %w", err)
	}
	sort.Strings(paths)
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading workspace: %w", err)
		}
		files.Worlds = append(files.Worlds, WorldFile{Name: strings.TrimSuffix(filepath.Base(path), ".yaml"), Content: string(content)})
	}

	h := sha256.New()
	for _, part := range []string{files.Policy, files.Constraints, files.Contracts} {
		fmt.Fprintf(h, "%d:%s", len(part), part)
	}
	for _, w := range files.Worlds {
		fmt.Fprintf(h, "%d:%s%d:%s", len(w.Name), w.Name, len(w.Content), w.Content)
	}
	files.Version = hex.EncodeToString(h.Sum(nil))[:16]
	return files, nil
}
