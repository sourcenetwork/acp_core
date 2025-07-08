package playground

import (
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/internal/sandbox"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	_ "github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	_ "github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	_ "github.com/sourcenetwork/acp_core/test"
)

var noopTheorem = `
Authorizations {}
Delegations {}
`

var setupData = &types.SandboxData{
	PolicyDefinition: `
				name: test
				resources:
				  file:
				    relations:
					  owner:
					    types:
						  - actor
					  reader:
					    types:
						  - actor
				    permissions:
					  read:
					    expr: owner + reader
					  write:
					    expr: owner
				`,
	Relationships: `
				file:readme#owner@did:example:bob
				file:readme#reader@did:example:alice
				`,
	PolicyTheorem: `
				Authorizations {
				  file:readme#read@did:example:bob
				  file:readme#write@did:example:bob

				  !file:readme#write@did:example:alice
				  file:readme#read@did:example:alice
				}
				Delegations {}
				ImpliedRelations {}
				`,
}

func Test_NewSandbox_ReturnsHandle(t *testing.T) {
	ctx := test.NewTestCtx(t)
	a := NewSandbox{
		Req: &types.NewSandboxRequest{
			Name:        "test",
			Description: "test",
		},
		Expected: &types.NewSandboxResponse{
			Record: &types.SandboxRecord{
				Handle:      1,
				Name:        "test",
				Description: "test",
				Initialized: false,
			},
		},
	}
	a.Run(ctx)
}

func Test_NewSandbox_UnamedSandbox_ReturnsSandboxWithHandleAsName(t *testing.T) {
	ctx := test.NewTestCtx(t)
	a := NewSandbox{
		Req: &types.NewSandboxRequest{
			Name:        "",
			Description: "test",
		},
		Expected: &types.NewSandboxResponse{
			Record: &types.SandboxRecord{
				Handle:      1,
				Name:        "1",
				Description: "test",
				Initialized: false,
			},
		},
	}
	a.Run(ctx)
}

func Test_NewSandbox_CanCreateSandboxWithoutDescription(t *testing.T) {
	ctx := test.NewTestCtx(t)
	a := NewSandbox{
		Req: &types.NewSandboxRequest{
			Name:        "test",
			Description: "",
		},
		Expected: &types.NewSandboxResponse{
			Record: &types.SandboxRecord{
				Handle:      1,
				Name:        "test",
				Description: "",
				Initialized: false,
			},
		},
	}
	a.Run(ctx)
}

func Test_SetState_EmptyTheoremErrors(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewAndSet{
		Data: &types.SandboxData{
			PolicyDefinition: `name: test`,
			Relationships:    ``,
			PolicyTheorem:    "",
		},
		Assertions: []Assertion{
			HasTheoremError("mismatched input"),
		},
	}
	a1.Run(ctx)
}

func Test_Evaluate_SandboxWithEmptyTheoremOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewAndSet{
		Data: &types.SandboxData{
			PolicyDefinition: `name: test`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	handle := a1.Run(ctx)

	a := VerifyTheorems{
		Handle: handle,
		Expected: &types.VerifyTheoremsResponse{
			Result: &types.AnnotatedPolicyTheoremResult{
				Theorem: &types.PolicyTheorem{
					AuthorizationTheorems: make([]*types.AuthorizationTheorem, 0),
					DelegationTheorems:    make([]*types.DelegationTheorem, 0),
				},
				Ok:           true,
				TheoremCount: 0,
				Failures:     0,
			},
		},
	}
	a.Run(ctx)
}

func Test_Evaluate_UninitializedSandboxCannotBeEvaluated(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewSandbox{
		Req: &types.NewSandboxRequest{},
	}
	handle := a1.Run(ctx)

	a := VerifyTheorems{
		Handle:      handle.Record.Handle,
		ExpectedErr: errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}

func Test_ListSandboxes_ReturnsExistingSandboxes(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a := NewAndSet{
		Data: &types.SandboxData{
			PolicyDefinition: `name: test1`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	a.Run(ctx)

	a = NewAndSet{
		Data: &types.SandboxData{
			PolicyDefinition: `name: test2`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	a.Run(ctx)

	action := ListSandboxes{
		Req:         &types.ListSandboxesRequest{},
		ExpectedLen: 2,
	}
	action.Run(ctx)
}

func Test_SetState_SettingValidStateReturnsOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	new := NewSandbox{
		Req: &types.NewSandboxRequest{
			Name:        "test",
			Description: "",
		},
	}
	resp := new.Run(ctx)

	a := SetState{
		Req: &types.SetStateRequest{
			Handle: resp.Record.Handle,
			Data: &types.SandboxData{
				PolicyDefinition: `
                name: test
                resources:
                  file:
                    relations:
                      owner:
                        types:
                          - actor
                      reader:
                        types:
                          - actor
                    permissions:
                      read:
                        expr: owner + reader
                      write:
                        expr: owner`,
				Relationships: `
				file:readme#owner@did:example:bob
				file:readme#reader@did:example:alice
				`,
				PolicyTheorem: `
				Authorizations {
				  file:readme#read@did:example:bob
				  file:readme#write@did:example:bob

				  !file:readme#write@did:example:alice
				  file:readme#read@did:example:alice
				}
				Delegations {}
				ImpliedRelations {}
				`,
			},
		},
	}
	a.Run(ctx)
}

func Test_GetCatalogue_ReturnsSandboxCatalogue(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewAndSet{
		Data: setupData,
	}
	handle := a1.Run(ctx)

	a := GetCatalogue{
		Req: &types.GetCatalogueRequest{
			Handle: handle,
		},
		Expected: &types.GetCatalogueResponse{
			Catalogue: &types.PolicyCatalogue{
				ActorResourceName: "actor",
				ResourceCatalogue: map[string]*types.ResourceCatalogue{
					"file": {
						Permissions: []string{
							"_can_manage_owner",
							"_can_manage_reader",
							"read",
							"write",
						},
						Relations: []string{
							"owner",
							"reader",
						},
						ObjectIds: []string{
							"readme",
						},
					},
				},
				Actors: []string{
					"did:example:bob",
				},
			},
		},
	}
	a.Run(ctx)
}

func Test_Simulate(t *testing.T) {
	ctx := test.NewTestCtx(t)

	data := types.SandboxData{
		PolicyDefinition: `
        name: test
        resources:
          file:
            relations:
              owner:
                types:
                  - actor
              reader:
                types:
                  - actor
            permissions:
              read:
                expr: owner + reader
              write:
                expr: owner
        `,
		Relationships: `
		file:abc#owner@did:ex:bob
		file:abc#reader@did:ex:alice
		`,
		PolicyTheorem: `
		Authorizations {
		  file:abc#read@did:ex:bob
		  file:abc#write@did:ex:bob
		  file:abc#read@did:ex:alice
		  file:abc#unknown_relation@did:ex:alice
		  ! file:abc#write@did:ex:alice
		}
		Delegations {
		  did:ex:bob > file:abc#reader
		  did:ex:alice > file:abc#reader
		}
		`,
	}

	resp, err := ctx.Playground.Simulate(ctx, &types.SimulateRequest{
		Data: &data,
	})
	require.NoError(t, err)
	require.True(t, resp.ValidData)
}

func Test_GetSandbox_ReturnsSandbox(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewSandbox{
		Req: &types.NewSandboxRequest{},
	}
	a1.Run(ctx)

	a1 = NewSandbox{
		Req: &types.NewSandboxRequest{},
	}
	a1.Run(ctx)

	resp, err := ctx.Playground.GetSandbox(ctx, &types.GetSandboxRequest{Handle: 1})

	want := &types.SandboxRecord{
		Handle:      1,
		Name:        "1",
		Description: "",
		Data:        nil,
		Scratchpad:  nil,
		Ctx:         nil,
		Initialized: false,
	}
	require.NoError(t, err)
	require.Equal(t, want, resp.Record)

	resp, err = ctx.Playground.GetSandbox(ctx, &types.GetSandboxRequest{Handle: 2})

	want = &types.SandboxRecord{
		Handle:      2,
		Name:        "2",
		Description: "",
		Data:        nil,
		Scratchpad:  nil,
		Ctx:         nil,
		Initialized: false,
	}
	require.NoError(t, err)
	require.Equal(t, want, resp.Record)
}

func Test_GetPlaygroundSamples_ReturnSamples(t *testing.T) {
	ctx := test.NewTestCtx(t)

	resp, err := ctx.Playground.GetSampleSandboxes(ctx, &types.GetSampleSandboxesRequest{})
	require.NoError(t, err)

	want := &types.GetSampleSandboxesResponse{
		Samples: sandbox.Samples,
	}
	require.Equal(t, want, resp)

	for _, data := range want.Samples {
		action := NewAndSet{
			Data: data.Data,
		}
		handle := action.Run(ctx)

		a2 := VerifyTheorems{
			Handle: handle,
		}
		response := a2.Run(ctx)
		require.True(t, response.Result.Ok)
	}
}
