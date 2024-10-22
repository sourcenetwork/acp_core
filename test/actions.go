package test

import (
	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/types"
	testutil "github.com/sourcenetwork/acp_core/test/util"
)

var DefaultTs = testutil.MustDateTimeToProto("2024-01-01 00:00:00")

type ActionState struct {
	PolicyId      string
	PolicyCreator string
}

type CreatePolicyAction struct {
	Policy   string
	Expected *types.Policy
}

func (a *CreatePolicyAction) Run(ctx *TestCtx) *types.Policy {
	req := types.CreatePolicyRequest{
		Policy:       a.Policy,
		MarshalType:  types.PolicyMarshalingType_SHORT_YAML,
		CreationTime: DefaultTs,
	}

	resp, err := ctx.Engine.CreatePolicy(ctx, &req)
	require.NoError(ctx.T, err)

	if a.Expected != nil {
		require.Equal(ctx.T, resp.Policy, a.Expected)
	}

	ctx.State.PolicyId = resp.Policy.Id
	principal, err := auth.ExtractPrincipal(ctx.Ctx)
	require.NoError(ctx.T, err)
	ctx.State.PolicyCreator = principal.Identifier()

	return resp.Policy
}

type RegisterObjectsAction struct {
	PolicyId string
	Objects  []*types.Object
}

func (a *RegisterObjectsAction) Run(ctx *TestCtx) {
	for _, obj := range a.Objects {
		req := types.RegisterObjectRequest{
			PolicyId:     a.PolicyId,
			Object:       obj,
			CreationTime: DefaultTs,
		}
		_, err := ctx.Engine.RegisterObject(ctx, &req)

		require.NoError(ctx.T, err)
	}
}

type SetRelationshipsAction struct {
	PolicyId      string
	Relationships []*types.Relationship
}

func (a *SetRelationshipsAction) Run(ctx *TestCtx) {
	for _, rel := range a.Relationships {
		req := types.SetRelationshipRequest{
			PolicyId:     a.PolicyId,
			CreationTime: DefaultTs,
			Relationship: rel,
		}
		_, err := ctx.Engine.SetRelationship(ctx, &req)
		require.NoError(ctx.T, err)
	}
}

type ArchiveObjectAction struct {
	PolicyId string
	Objects  []*types.Object
}

func (a *ArchiveObjectAction) Run(ctx *TestCtx) {
	for _, obj := range a.Objects {
		req := types.ArchiveObjectRequest{
			PolicyId: a.PolicyId,
			Object:   obj,
		}
		_, err := ctx.Engine.ArchiveObject(ctx, &req)

		require.NoError(ctx.T, err)
	}
}

type PolicySetupAction struct {
	Policy                string
	PolicyCreator         string
	ObjectsPerActor       map[string][]*types.Object
	RelationshipsPerActor map[string][]*types.Relationship
}

func (a *PolicySetupAction) Run(ctx *TestCtx) {
	ctx.SetPrincipal(a.PolicyCreator)
	a1 := CreatePolicyAction{
		Policy: a.Policy,
	}
	policy := a1.Run(ctx)

	for actor, objects := range a.ObjectsPerActor {
		ctx.SetPrincipal(actor)
		objsAction := RegisterObjectsAction{
			PolicyId: policy.Id,
			Objects:  objects,
		}
		objsAction.Run(ctx)
	}

	for actor, relationships := range a.RelationshipsPerActor {
		ctx.SetPrincipal(actor)
		relsAction := SetRelationshipsAction{
			PolicyId:      policy.Id,
			Relationships: relationships,
		}
		relsAction.Run(ctx)
	}
}
