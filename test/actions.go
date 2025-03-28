package test

import (
	"github.com/stretchr/testify/require"

	prototypes "github.com/cosmos/gogoproto/types"
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
	Metadata *types.SuppliedMetadata
}

func (a *CreatePolicyAction) Run(ctx *TestCtx) *types.Policy {
	req := types.CreatePolicyRequest{
		Policy:      a.Policy,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
		Metadata:    a.Metadata,
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
	PolicyId     string
	Object       *types.Object
	NewOwner     string
	NewTimestamp *prototypes.Timestamp
	Expected     *types.AmendRegistrationResponse
	ExpectedErr  error
}

func (a *AmendRegistrationAction) Run(ctx *TestCtx) *types.AmendRegistrationResponse {
	req := types.AmendRegistrationRequest{
		PolicyId:      a.PolicyId,
		Object:        a.Object,
		NewOwner:      types.NewActor(a.NewOwner),
		NewCreationTs: a.NewTimestamp,
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

type RevealRegistrationAction struct {
	PolicyId    string
	Object      *types.Object
	Ts          *prototypes.Timestamp
	Metadata    *types.SuppliedMetadata
	Expected    *types.RelationshipRecord
	ExpectedErr error
}

func (a *RevealRegistrationAction) Run(ctx *TestCtx) *types.RevealRegistrationResponse {
	req := types.RevealRegistrationRequest{
		PolicyId:   a.PolicyId,
		Object:     a.Object,
		CreationTs: a.Ts,
		Metadata:   a.Metadata,
	}
	resp, err := ctx.Engine.RevealRegistration(ctx, &req)

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		if a.Expected != nil {
			require.Equal(ctx.T, a.Expected, resp.Record)
		}
	}
	return resp
}

type EditPolicyAction struct {
	PolicyId         string
	Policy           string
	ExpectedErr      error
	Expected         *types.Policy
	ExpectedMetadata *types.RecordMetadata
}

func (a *EditPolicyAction) Run(ctx *TestCtx) *types.Policy {
	req := types.EditPolicyRequest{
		PolicyId:    a.PolicyId,
		Policy:      a.Policy,
		MarshalType: types.PolicyMarshalingType_SHORT_YAML,
	}

	resp, err := ctx.Engine.EditPolicy(ctx, &req)
	getResp, getErr := ctx.Engine.GetPolicy(ctx, &types.GetPolicyRequest{
		Id: a.PolicyId,
	})

	if a.ExpectedErr != nil {
		require.ErrorIs(ctx.T, err, a.ExpectedErr)
	} else {
		require.NoError(ctx.T, err)
		require.NoError(ctx.T, getErr)
	}

	if a.Expected != nil {
		require.Equal(ctx.T, a.Expected, resp.Record.Policy)
		require.Equal(ctx.T, a.Expected, getResp.Record.Policy)
	}
	if a.ExpectedMetadata != nil {
		require.Equal(ctx.T, a.ExpectedMetadata, resp.Record.Metadata)
		require.Equal(ctx.T, a.ExpectedMetadata, getResp.Record.Metadata)
	}

	if resp != nil {
		return resp.Record.Policy
	}
	return nil
}
