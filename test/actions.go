package test

import (
	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/auth"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

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
		Policy:      a.Policy,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}

	resp, err := ctx.Engine.CreatePolicy(ctx, &req)
	require.NoError(ctx.T, err)

	if a.Expected != nil {
		require.Equal(ctx.T, resp.Record.Policy, a.Expected)
	}

	ctx.State.PolicyId = resp.Record.Policy.Id
	principal, err := auth.ExtractPrincipal(ctx.Ctx)
	require.NoError(ctx.T, err)
	ctx.State.PolicyCreator = principal.Identifier

	return resp.Record.Policy
}

type RegisterObjectsAction struct {
	PolicyId string
	Objects  []*types.Object
}

func (a *RegisterObjectsAction) Run(ctx *TestCtx) {
	for _, obj := range a.Objects {
		req := types.RegisterObjectRequest{
			PolicyId: a.PolicyId,
			Object:   obj,
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
			Relationship: rel,
		}
		_, err := ctx.Engine.SetRelationship(ctx, &req)
		require.NoError(ctx.T, err)
	}
}

type ArchiveObjectAction struct {
	PolicyId    string
	Object      *types.Object
	Expected    *types.ArchiveObjectResponse
	ExpectedErr error
}

func (a *ArchiveObjectAction) Run(ctx *TestCtx) {
	req := types.ArchiveObjectRequest{
		PolicyId: a.PolicyId,
		Object:   a.Object,
	}
	resp, err := ctx.Engine.ArchiveObject(ctx, &req)

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		if a.Expected != nil {
			require.Equal(ctx.T, a.Expected, resp)
		}
	}
}

type UnarchiveObjectAction struct {
	PolicyId    string
	Object      *types.Object
	Expected    *types.UnarchiveObjectResponse
	ExpectedErr error
}

func (a *UnarchiveObjectAction) Run(ctx *TestCtx) {
	req := types.UnarchiveObjectRequest{
		PolicyId: a.PolicyId,
		Object:   a.Object,
	}
	resp, err := ctx.Engine.UnarchiveObject(ctx, &req)

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		if a.Expected != nil {
			require.Equal(ctx.T, a.Expected, resp)
		}
	}
}

type TransferObjectAction struct {
	PolicyId    string
	Object      *types.Object
	NewOwner    string
	Expected    *types.TransferObjectResponse
	ExpectedErr error
}

func (a *TransferObjectAction) Run(ctx *TestCtx) *types.TransferObjectResponse {
	req := types.TransferObjectRequest{
		PolicyId: a.PolicyId,
		Object:   a.Object,
		NewOwner: types.NewActor(a.NewOwner),
	}
	resp, err := ctx.Engine.TransferObject(ctx, &req)

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		if a.Expected != nil {
			require.Equal(ctx.T, a.Expected, resp)
		}
	}
	return resp
}

type AmendRegistrationAction struct {
	PolicyId    string
	Object      *types.Object
	NewOwner    string
	Expected    *types.AmendRegistrationResponse
	ExpectedErr error
}

func (a *AmendRegistrationAction) Run(ctx *TestCtx) *types.AmendRegistrationResponse {
	req := types.AmendRegistrationRequest{
		PolicyId: a.PolicyId,
		Object:   a.Object,
		NewOwner: types.NewActor(a.NewOwner),
	}
	resp, err := ctx.Engine.AmendRegistration(ctx, &req)

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		if a.Expected != nil {

			require.Equal(ctx.T, a.Expected, resp)
		}
	}
	return resp
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
