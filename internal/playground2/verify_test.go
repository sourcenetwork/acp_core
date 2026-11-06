package playground2

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/internal/telemetry"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

const folderRequirements = `
requirements:
- name: doc-view
  statement: you can see a doc if you can see its folder, or you're on the doc
  contracts:
  - name: doc-view-needs-a-role
    description: Viewing a doc requires owning it, being its viewer, or viewing its folder.
    status: accepted
    properties:
    - type: assert
      description: every doc viewer holds a role on the doc or can view its folder
      formula: |
        all d: Doc, a: Actor | a in d.doc_view implies
          (a in d.doc_owner + d.doc_viewer or a in d.folder.folder_view)
    - type: pred
      description: docs can be visible to someone other than their owner
      formula: |
        some d: Doc | some d.doc_view - d.doc_owner
  - name: only-direct-viewers
    description: Only a doc's owner and viewers can view it. False, folder viewers can too.
    properties:
    - type: pred
      description: docs can have viewers
      formula: |
        some Doc.doc_viewer
    - type: assert
      description: every doc viewer is the owner or a viewer
      formula: |
        all d: Doc, a: Actor | a in d.doc_view implies a in d.doc_owner + d.doc_viewer
  - name: folders-are-private
    description: Nobody but the owner sees a folder.
- name: visibility
  statement: some docs nobody sees
  notes:
  - "Impossible as stated: the owner always sees their doc. Asked the user."
  contracts:
  - name: invisible-docs
    description: A doc can be invisible to everyone.
    properties:
    - type: pred
      description: a doc nobody can view
      formula: |
        some d: Doc | no d.doc_view
  - name: typo
    description: A formula with an error.
    properties:
    - type: pred
      description: some docs exist
      formula: |
        some Doc
    - type: assert
      description: refers to an undefined relation
      formula: |
        all d: Doc | d.nope in Actor
`

func Test_Verify_Contracts(t *testing.T) {
	g := analyzer(t)
	requirements, err := pc.UnmarshalRequirements(folderRequirements)
	require.NoError(t, err)
	model := folderModel(t)
	model.Requirements = requirements

	resp, err := g.Verify(context.Background(), &types.VerifyRequest{SystemModel: model})
	require.NoError(t, err)
	require.Equal(t, uint32(defaultScope), resp.Scope)
	require.Len(t, resp.Results, 5)
	results := make(map[string]*types.ContractResult)
	for _, r := range resp.Results {
		results[r.Contract] = r
	}

	holds := results["doc-view-needs-a-role"]
	require.Equal(t, "doc-view", holds.Requirement)
	require.Equal(t, types.VerifyStatus_SATISFIED, holds.Status)
	require.Len(t, holds.Properties, 2)
	require.Nil(t, holds.Properties[0].World, "a satisfied assert has no world")
	require.NotEmpty(t, holds.Properties[1].World.Relationships, "a satisfied pred carries an example")

	draft := results["folders-are-private"]
	require.Equal(t, types.VerifyStatus_DRAFT, draft.Status)
	require.Empty(t, draft.Properties)

	require.Equal(t, types.VerifyStatus_VIOLATED, results["invisible-docs"].Status)

	typo := results["typo"]
	require.Equal(t, types.VerifyStatus_UNVERIFIED, typo.Status)
	require.Equal(t, types.VerifyStatus_SATISFIED, typo.Properties[0].Status)
	require.Contains(t, typo.Properties[1].Message, "at formula line 1 column 16")
	require.Contains(t, typo.Properties[1].Message, `"nope"`)

	// the counterexample names the doc and actor, and the engine agrees the actor can view the doc
	violated := results["only-direct-viewers"]
	require.Equal(t, types.VerifyStatus_VIOLATED, violated.Status)
	require.Equal(t, types.VerifyStatus_SATISFIED, violated.Properties[0].Status)
	counterexample := violated.Properties[1]
	require.Equal(t, types.VerifyStatus_VIOLATED, counterexample.Status)
	bindings := make(map[string]string)
	for _, b := range counterexample.Bindings {
		require.Len(t, b.Values, 1)
		bindings[b.Name] = b.Values[0]
	}
	resource, id, _ := strings.Cut(bindings["d"], ":")
	ctx := telemetry.InitRequestContext(context.Background())
	check, err := Check(ctx, &types.CheckRequest{
		Scenario:   &types.Scenario{SystemModel: model, World: counterexample.World},
		Object:     types.NewObject(resource, id),
		Permission: "view",
		Actor:      types.NewActor(bindings["a"]),
	})
	require.NoError(t, err)
	require.True(t, check.Allowed)
	require.NotContains(t, counterexample.World.Relationships, bindings["d"]+"#viewer@"+bindings["a"])

	out := &bytes.Buffer{}
	failed := PrintReport(out, requirements, resp, nil)
	require.Equal(t, 3, failed)
	require.Contains(t, out.String(), "1 of 4 contracts satisfied up to scope 3, 1 drafts")
	require.Contains(t, out.String(), "doc-view-needs-a-role, accepted")
	require.Contains(t, out.String(), "note: Impossible as stated")
}

func Test_Verify_OnlyDraftsNeedNoAlloy(t *testing.T) {
	g := &Analyzer{AlloyPath: "/nonexistent/alloy"}
	model := &types.SystemModel{Policy: "not even a policy", Requirements: []*types.Requirement{
		{Name: "r", Statement: "s", Contracts: []*types.Contract{{Name: "c", Description: "d"}}},
	}}
	resp, err := g.Verify(context.Background(), &types.VerifyRequest{SystemModel: model})
	require.NoError(t, err)
	require.Equal(t, types.VerifyStatus_DRAFT, resp.Results[0].Status)
}

func Test_Verify_InvalidRequests(t *testing.T) {
	g := &Analyzer{AlloyPath: "alloy"}
	_, err := g.Verify(context.Background(), &types.VerifyRequest{})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)

	_, err = g.Verify(context.Background(), &types.VerifyRequest{SystemModel: &types.SystemModel{Policy: folderPolicy}})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)

	model := &types.SystemModel{Policy: folderPolicy, Requirements: []*types.Requirement{
		{Name: "r", Statement: "s", Contracts: []*types.Contract{{Name: "c", Description: "d"}}},
	}}
	_, err = g.Verify(context.Background(), &types.VerifyRequest{SystemModel: model, Config: &types.VerifyConfig{Scope: maxScope + 1}})
	require.ErrorIs(t, err, errors.ErrorType_BAD_INPUT)
}

func Test_UnmarshalRequirements_Invalid(t *testing.T) {
	contract := func(extra string) string {
		return "requirements:\n- name: r\n  statement: s\n  contracts:\n  - name: c\n    description: d\n" + extra
	}
	for name, doc := range map[string]string{
		"no statement":          "requirements:\n- name: r\n",
		"bad requirement name":  "requirements:\n- name: 1r\n  statement: s\n",
		"no description":        "requirements:\n- name: r\n  statement: s\n  contracts:\n  - name: c\n",
		"bad status":            contract("    status: confirmed\n"),
		"missing property type": contract("    properties:\n    - description: p\n      formula: some Doc\n"),
		"no property desc":      contract("    properties:\n    - type: pred\n      formula: some Doc\n"),
		"no formula":            contract("    properties:\n    - type: pred\n      description: p\n"),
		"clashing contracts":    contract("  - name: c\n    description: d\n"),
		"unknown key":           contract("    checks: []\n"),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := pc.UnmarshalRequirements(doc)
			require.Error(t, err)
		})
	}

	reqs, err := pc.UnmarshalRequirements(contract(""))
	require.NoError(t, err)
	require.Equal(t, types.ContractStatus_PROPOSED, reqs[0].Contracts[0].Status, "proposed is the default")
}
