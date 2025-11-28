package policy

import (
	"testing"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/internal/specification"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	"github.com/stretchr/testify/require"
)

func TestEditPolicy_CannotRemoveResource(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	// Given Policy
	oldPol := `
name: policy
resources:
- name: file
`

	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When I attempt to remove a resource
	new := `name: policy
spec: none
`

	a := test.EditPolicyAction{
		PolicyId:    ctx.State.PolicyId,
		Policy:      new,
		ExpectedErr: ppp.ErrPreserveResource,
	}
	a.Run(ctx)
}

func TestEditPolicy_RemovingOwnerRelation_DiscretionaryTransformerRestoresIt(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	// Given Policy
	oldPol := `
name: policy
resources:
- name: file
`

	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When I attempt to remove the owner relation
	new := `
name: policy
resources:
- name: file
`

	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)
	t.Logf("pol: %v", pol)
	want := &types.Relation{
		Name: "owner",
		Doc:  "owner relations represents the object owner",
		VrTypes: []*types.Restriction{
			{
				ResourceName: "actor",
			},
		},
	}
	require.Equal(t, want, pol.GetResourceByName("file").GetRelationByName("owner"))
}

func TestEditPolicy_CannotRenameActorResource(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `actor:
  name: test
name: policy
spec: none
`

	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When I attempt to rename the actor resource
	new := `actor:
  name: new-actor-name
name: policy
spec: none
`
	a := test.EditPolicyAction{
		PolicyId:    ctx.State.PolicyId,
		Policy:      new,
		ExpectedErr: ppp.ErrPreserveResource,
	}
	a.Run(ctx)
}
func TestEditPolicy_CannotChangeSpec(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - name: read
  - name: write
spec: defra
`

	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When the I edit a policy with a new policy spec
	new := `
name: policy
resources:
- name: file
  permissions:
  - name: read
  - name: write
`

	a := test.EditPolicyAction{
		PolicyId:    ctx.State.PolicyId,
		Policy:      new,
		ExpectedErr: ppp.ErrImmutableSpec,
	}
	a.Run(ctx)
}

func TestEditPolicy_CannotEditPolicyThatDoesntExist(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	new := `
name: policy
resources:
- name: file
  permissions:
  - name: read
  - name: write
`

	a := test.EditPolicyAction{
		PolicyId:    "some-policy-id",
		Policy:      new,
		ExpectedErr: errors.ErrorType_NOT_FOUND,
	}
	a.Run(ctx)
}

func TestEditPolicy_CanAddRelation(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When the I add a new relation to an existing resource
	new := `
name: policy
resources:
- name: file
  relations:
  - name: reader
  - name: writer
`
	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	// then the new relation exists in the policy
	want := &types.Relation{
		Name:    "writer",
		VrTypes: []*types.Restriction{},
	}
	require.Equal(ctx.T, want, pol.GetResourceByName("file").GetRelationByName("writer"))
}

func TestEditPolicy_CanAddResource(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When the I add a new resource with a relation
	new := `
name: policy
resources:
- name: file
  relations:
  - name: reader
- name: group
  relations:
  - name: member
`
	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	require.Equal(ctx.T, "group", pol.GetResourceByName("group").Name)
	require.Equal(ctx.T, "member", pol.GetResourceByName("group").GetRelationByName("member").Name)
}

func TestEditPolicy_CanEditNameAndAttrs(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `description: a test policy
meta:
  key: val
  key2: val2
name: policy
spec: none
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	// When the I add a new resource with a relation
	new := `description: another test policy
meta:
  key: val2
  key2: val3
  key3: val
name: new name
spec: none
`
	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	want := &types.Policy{
		Id:          "bc7eb5a8c500111b2459a92ae23f4848537e49599df1b8d70636b5aacb47bd5f",
		Name:        "new name",
		Description: "another test policy",
		Resources:   []*types.Resource{},
		ActorResource: &types.ActorResource{
			Name: "actor",
		},
		Attributes: map[string]string{
			"key":  "val2",
			"key2": "val3",
			"key3": "val",
		},
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
	}
	require.Equal(ctx.T, want, pol)
}

func TestEditPolicy_CanEditPermissionExpr(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
  - name: writer
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	new := `
name: policy
resources:
- name: file
  permissions:
  - expr: writer
    name: read
  relations:
  - name: reader
  - name: writer
`
	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	want := "(owner + writer)"
	require.Equal(ctx.T, want, pol.GetResourceByName("file").GetPermissionByName("read").Expression)
}

func TestEditPolicy_CannotRemoveDefraPermissionsFromDefraPolicy(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - name: read
  - name: write
spec: defra
`

	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	new := `
name: policy
resources:
- name: file
spec: defra
`
	a := test.EditPolicyAction{
		PolicyId:    ctx.State.PolicyId,
		Policy:      new,
		ExpectedErr: specification.ErrDefraSpec,
	}
	a.Run(ctx)
}

func TestEditPolicy_CanAddPermission(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	new := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  - expr: reader
    name: write
  relations:
  - name: reader
`
	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	// then the new permission exists in the policy
	want := &types.Permission{
		Name:       "write",
		Expression: "(owner + reader)",
	}
	require.Equal(ctx.T, want, pol.GetResourceByName("file").GetPermissionByName("write"))
}

func TestEditPolicy_CanRemovePermission(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	new := `
name: policy
resources:
- name: file
  relations:
  - name: reader
`

	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	require.Nil(ctx.T, (pol.GetResourceByName("file").GetPermissionByName("read")))
}

func TestEditPolicy_DoesNotChangeSuppliedMetadata(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	new := `
name: policy
resources:
- name: file
  relations:
  - name: reader
`

	a := test.EditPolicyAction{
		PolicyId: ctx.State.PolicyId,
		Policy:   new,
	}
	pol := a.Run(ctx)

	require.Nil(ctx.T, (pol.GetResourceByName("file").GetPermissionByName("read")))
}

func TestEditPolicyMetadata_CanEditMetadata(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `
name: policy
resources:
- name: file
  permissions:
  - expr: reader
    name: read
  relations:
  - name: reader
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
		Metadata: &types.SuppliedMetadata{
			Blob: []byte{0, 1, 0},
		},
	}
	a1.Run(ctx)

	newMetadata := &types.SuppliedMetadata{
		Blob: []byte{1, 2, 3},
	}
	resp, err := ctx.Engine.EditPolicyMetadata(ctx, &types.EditPolicyMetadataRequest{
		PolicyId: ctx.State.PolicyId,
		Metadata: newMetadata,
	})
	require.NoError(t, err)
	require.Equal(t, newMetadata, resp.Record.Metadata.Supplied)
}

func TestEditPolicyMetadata_NonOwnerCannotEditMetadata(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `name: policy
spec: none
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	ctx.SetRootPrincipal()
	resp, err := ctx.Engine.EditPolicyMetadata(ctx, &types.EditPolicyMetadataRequest{
		PolicyId: ctx.State.PolicyId,
		Metadata: &types.SuppliedMetadata{
			Blob: []byte{1, 2, 3},
		},
	})
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
	require.Nil(t, resp)
}

func TestEditPolicy_NonOwnerCannotEditMetadata(t *testing.T) {
	ctx := test.NewTestCtx(t)
	ctx.SetPrincipal("bob")

	oldPol := `name: policy
spec: none
`
	a1 := test.CreatePolicyAction{
		Policy: oldPol,
	}
	a1.Run(ctx)

	ctx.SetRootPrincipal()
	resp, err := ctx.Engine.EditPolicy(ctx, &types.EditPolicyRequest{
		PolicyId: ctx.State.PolicyId,
		Policy:   "",
	})
	require.ErrorIs(t, err, errors.ErrorType_UNAUTHORIZED)
	require.Nil(t, resp)
}
