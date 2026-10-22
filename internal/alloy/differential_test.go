package alloy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/internal/simulator"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// workspacePolicy exercises recursive tuple to usersets, nested groups,
// usersets to permissions, management rules, difference and intersection
const workspacePolicy = `
name: workspace
resources:
- name: group
  relations:
  - name: admin
    manages: [member]
    types: [actor]
  - name: member
    types: [actor, group->member]
  permissions:
  - name: view
    expr: admin + member
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: editor
    manages: [viewer]
    types: [actor, group->view]
  - name: viewer
    types: [actor, group->member]
  permissions:
  - name: edit
    expr: editor + parent->edit
  - name: view
    expr: viewer + edit + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: reviewer
    types: [actor]
  - name: banned
    types: [actor]
  permissions:
  - name: view
    expr: folder->view - banned
  - name: review
    expr: reviewer & folder->view
`

const (
	instanceActors  = 3
	instanceObjects = 2
	instanceSeeds   = 5
)

// Test_Differential_EmittedModelAgreesWithEngine generates random relationship sets
// and asserts that the emitted Alloy model computes the same
// relations, permissions and authorities as the acp_core engine.
func Test_Differential_EmittedModelAgreesWithEngine(t *testing.T) {
	if _, err := exec.LookPath("alloy"); err != nil {
		t.Skip("alloy not found in PATH")
	}
	if testing.Short() {
		t.Skip("differential tests are slow")
	}

	policies := map[string]string{
		"filesystem sample": sandbox.Samples[0].Data.PolicyDefinition,
		"workspace":         workspacePolicy,
	}
	for name, policy := range policies {
		t.Run(name, func(t *testing.T) {
			model, err := pc.Compile(policy, types.PolicyMarshalingType_YAML)
			require.NoError(t, err)

			for seed := int64(0); seed < instanceSeeds; seed++ {
				inst := randomInstance(model, rand.New(rand.NewSource(seed)))
				expected := evaluateWithEngine(t, policy, model, inst)
				failed := checkWithAlloy(t, model, inst, expected)
				require.Empty(t, failed, "seed %v: alloy disagrees with engine\nrelationships:\n%v", seed, inst.relationships())
			}
		})
	}
}

type tuple struct {
	field   *pc.Field
	object  string
	subject string
}

type instance struct {
	model   *pc.Model
	actors  []string
	objects map[string][]string
	tuples  []tuple
}

func randomInstance(model *pc.Model, rnd *rand.Rand) *instance {
	inst := &instance{model: model, objects: make(map[string][]string)}
	for i := 0; i < instanceActors; i++ {
		inst.actors = append(inst.actors, fmt.Sprintf("Actor%v", i))
	}
	for _, r := range model.Resources {
		for i := 0; i < instanceObjects; i++ {
			inst.objects[r.Name] = append(inst.objects[r.Name], fmt.Sprintf("%v%v", sigName(r.Name), i))
		}
	}
	for _, r := range model.Resources {
		for _, f := range r.Fields {
			for _, obj := range inst.objects[r.Name] {
				candidates := inst.atoms(f.Target)
				if f.One {
					inst.tuples = append(inst.tuples, tuple{f, obj, candidates[rnd.Intn(len(candidates))]})
					continue
				}
				for _, subject := range candidates {
					if rnd.Float64() < 0.35 {
						inst.tuples = append(inst.tuples, tuple{f, obj, subject})
					}
				}
			}
		}
	}
	return inst
}

func (inst *instance) atoms(sig string) []string {
	if sig == pc.ActorTarget {
		return inst.actors
	}
	return inst.objects[sig]
}

// relationships renders the instance in the PolicyTheoremDSL relationship syntax,
// with owner relationships first so objects are registered before use
func (inst *instance) relationships() string {
	var owners, others []string
	for _, tup := range inst.tuples {
		subject := "did:user:" + tup.subject
		if tup.field.Target != pc.ActorTarget {
			subject = tup.field.Target + ":" + tup.subject
			if tup.field.Userset != nil {
				subject += "#" + tup.field.Userset.Name
			}
		}
		line := fmt.Sprintf("%v:%v#%v@%v", tup.field.Resource, tup.object, tup.field.Relation, subject)
		if tup.field.One {
			owners = append(owners, line)
		} else {
			others = append(others, line)
		}
	}
	return strings.Join(append(owners, others...), "\n")
}

// evaluateWithEngine returns the tuples of every relation, permission and authority computation,
// keyed by the computation's alloy name, as evaluated by the acp_core engine
func evaluateWithEngine(t *testing.T, policy string, model *pc.Model, inst *instance) map[string][]string {
	authorizations := &strings.Builder{}
	delegations := &strings.Builder{}
	for _, r := range model.Resources {
		for _, obj := range inst.objects[r.Name] {
			for _, actor := range inst.actors {
				for _, c := range values(r) {
					fmt.Fprintf(authorizations, "%v:%v#%v@did:user:%v\n", r.Name, obj, c.Name, actor)
				}
				for _, c := range r.Authorities {
					fmt.Fprintf(delegations, "did:user:%v > %v:%v#%v\n", actor, r.Name, obj, c.Name)
				}
			}
		}
	}

	manager, err := runtime.NewRuntimeManager()
	require.NoError(t, err)
	ctx := telemetry.InitRequestContext(context.Background())
	resp, err := simulator.HandleSimulate(ctx, manager, &types.SimulateRequest{
		Data: &types.SandboxData{
			PolicyDefinition: policy,
			Relationships:    inst.relationships(),
			PolicyTheorem:    fmt.Sprintf("Authorizations {\n%v}\nDelegations {\n%v}\n", authorizations, delegations),
		},
	})
	require.NoError(t, err)
	require.True(t, resp.ValidData, "invalid sandbox data: %v", resp.Errors)

	computations := make(map[string]*pc.Computation)
	for _, r := range model.Resources {
		for _, c := range append(values(r), r.Authorities...) {
			computations[computationName(c)] = c
		}
	}
	expected := make(map[string][]string)
	for name := range computations {
		expected[name] = nil
	}
	add := func(kind pc.Kind, op *types.Operation, actor string, result *types.Result) {
		require.NotEqual(t, types.ResultStatus_Error, result.Status, result.Message)
		if result.Status != types.ResultStatus_Accept {
			return
		}
		name := computationName(&pc.Computation{Kind: kind, Resource: op.Object.Resource, Name: op.Permission})
		if kind == pc.KindRelation && computations[name].Kind == pc.KindPermission {
			kind = pc.KindPermission
		}
		atom := op.Object.Id + "->" + strings.TrimPrefix(actor, "did:user:")
		expected[name] = append(expected[name], atom)
	}
	for _, res := range resp.Result.AuthorizationTheoremsResult {
		thm := res.Result.Theorem
		add(pc.KindRelation, thm.Operation, thm.Actor.Id, res.Result.Result)
	}
	for _, res := range resp.Result.DelegationTheoremsResult {
		thm := res.Result.Theorem
		add(pc.KindAuthority, thm.Operation, thm.Actor.Id, res.Result.Result)
	}
	return expected
}

// checkWithAlloy pins inst in an alloy module and checks every computation against expected.
// Returns the computations whose value differs.
func checkWithAlloy(t *testing.T, model *pc.Model, inst *instance, expected map[string][]string) []string {
	dir := t.TempDir()
	e := newEmitter(model, "")
	for _, f := range Emit(model, "") {
		require.NoError(t, os.WriteFile(filepath.Join(dir, f.Name), []byte(f.Content), 0644))
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "open %v\nopen %v_management\n\n", e.module, e.module)
	fmt.Fprintf(b, "one sig %v extends Actor {}\n", strings.Join(inst.actors, ", "))
	for _, r := range model.Resources {
		fmt.Fprintf(b, "one sig %v extends %v {}\n", strings.Join(inst.objects[r.Name], ", "), sigName(r.Name))
	}

	b.WriteString("\nfact instance {\n")
	fmt.Fprintf(b, "    Actor = %v\n", strings.Join(inst.actors, " + "))
	for _, r := range model.Resources {
		fmt.Fprintf(b, "    %v = %v\n", sigName(r.Name), strings.Join(inst.objects[r.Name], " + "))
		for _, f := range r.Fields {
			var atoms []string
			for _, tup := range inst.tuples {
				if tup.field == f {
					atoms = append(atoms, tup.object+"->"+tup.subject)
				}
			}
			fmt.Fprintf(b, "    (%v <: %v) = %v\n", sigName(r.Name), f.Name, relation(atoms))
		}
	}
	b.WriteString("}\n\n")

	names := make([]string, 0, len(expected))
	for name := range expected {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(b, "check eq_%v { %v = %v } for 10\n", name, name, relation(expected[name]))
	}

	path := filepath.Join(dir, "instance.als")
	require.NoError(t, os.WriteFile(path, []byte(b.String()), 0644))
	out, err := exec.Command("alloy", "exec", "-q", "-f", "-t", "none", "-c", "*", "-o", filepath.Join(dir, "out"), path).CombinedOutput()
	require.NoError(t, err, string(out))

	receiptBytes, err := os.ReadFile(filepath.Join(dir, "out", "receipt.json"))
	require.NoError(t, err)
	var receipt struct {
		Commands map[string]struct {
			Solution []json.RawMessage `json:"solution"`
		} `json:"commands"`
	}
	require.NoError(t, json.Unmarshal(receiptBytes, &receipt))
	require.Len(t, receipt.Commands, len(names), "alloy did not run every check:\n%v", string(out))

	var failed []string
	for name, cmd := range receipt.Commands {
		if len(cmd.Solution) > 0 {
			failed = append(failed, fmt.Sprintf("%v (expected %v)", name, relation(expected[strings.TrimPrefix(name, "eq_")])))
		}
	}
	sort.Strings(failed)
	return failed
}

func relation(atoms []string) string {
	if len(atoms) == 0 {
		return "none->none"
	}
	sort.Strings(atoms)
	return strings.Join(atoms, " + ")
}

func values(r *pc.Resource) []*pc.Computation {
	return append(append([]*pc.Computation{}, r.Relations...), r.Permissions...)
}
