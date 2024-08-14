package playground

import (
	"testing"

	_ "github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	_ "github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/playground"
	"github.com/sourcenetwork/acp_core/pkg/types"
	_ "github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/test"
	_ "github.com/sourcenetwork/acp_core/test"
)

var noopTheorem = `
Authorizations {}
Delegations {}
`

var setupData = &playground.SandboxData{
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
		Req: &playground.NewSandboxRequest{
			Name:        "test",
			Description: "test",
		},
		Expected: &playground.NewSandboxResponse{
			Record: &playground.SandboxRecord{
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
		Req: &playground.NewSandboxRequest{
			Name:        "",
			Description: "test",
		},
		Expected: &playground.NewSandboxResponse{
			Record: &playground.SandboxRecord{
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
		Req: &playground.NewSandboxRequest{
			Name:        "test",
			Description: "",
		},
		Expected: &playground.NewSandboxResponse{
			Record: &playground.SandboxRecord{
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
		Data: &playground.SandboxData{
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
		Data: &playground.SandboxData{
			PolicyDefinition: `name: test`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	handle := a1.Run(ctx)

	a := VerifyTheorems{
		Req: &playground.VerifyTheoremsRequest{
			Handle: handle,
		},
		Expected: &playground.VerifyTheoremsResponse{
			Result: &types.AnnotatedPolicyTheoremResult{
				Theorem: &types.PolicyTheorem{
					AuthorizationTheorems: make([]*types.AuthorizationTheorem, 0),
					DelegationTheorems:    make([]*types.DelegationTheorem, 0),
				},
			},
		},
	}
	a.Run(ctx)
}

func Test_Evaluate_UninitializedSandboxCannotBeEvaluated(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a1 := NewSandbox{
		Req: &playground.NewSandboxRequest{},
	}
	handle := a1.Run(ctx)

	a := VerifyTheorems{
		Req: &playground.VerifyTheoremsRequest{
			Handle: handle.Record.Handle,
		},
		ExpectedErr: errors.ErrorType_OPERATION_FORBIDDEN,
	}
	a.Run(ctx)
}

func Test_ListSandboxes_ReturnsExistingSandboxes(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a := NewAndSet{
		Data: &playground.SandboxData{
			PolicyDefinition: `name: test1`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	a.Run(ctx)

	a = NewAndSet{
		Data: &playground.SandboxData{
			PolicyDefinition: `name: test2`,
			Relationships:    ``,
			PolicyTheorem:    noopTheorem,
		},
	}
	a.Run(ctx)

	action := ListSandboxes{
		Req: &playground.ListSandboxesRequest{},
		Expected: &playground.ListSandboxesResponse{
			Records: []*playground.SandboxRecord{
				{
					Handle:      1,
					Name:        "test1",
					Description: "",
					Data:        nil,
					Scratchpad:  nil,
					Ctx:         nil,
					Initialized: false,
				},
				{
					Handle:      2,
					Name:        "test2",
					Description: "",
					Data:        nil,
					Scratchpad:  nil,
					Ctx:         nil,
					Initialized: false,
				},
			},
		},
	}
	action.Run(ctx)
}

func Test_SetState_SettingValidStateReturnsOk(t *testing.T) {
	ctx := test.NewTestCtx(t)

	new := NewSandbox{
		Req: &playground.NewSandboxRequest{
			Name:        "test",
			Description: "",
		},
	}
	resp := new.Run(ctx)

	a := SetState{
		Req: &playground.SetStateRequest{
			Handle: resp.Record.Handle,
			Data: &playground.SandboxData{
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
		Req: &playground.GetCatalogueRequest{
			Handle: handle,
		},
		Expected: &playground.GetCatalogueResponse{
			Catalogue: &types.PolicyCatalogue{
				ActorResourceName: "actor",
				ResourceCatalogue: map[string]*types.ResourceCatalogue{
					"file": &types.ResourceCatalogue{
						Permissions: []string{
							"read",
							"write",
							"_can_manage_owner",
							"_can_manage_reader",
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

func Test_SetState_(t *testing.T) {}

func Test_VerifyTheorem_(t *testing.T) {}
