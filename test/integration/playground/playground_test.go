package playground

import (
	"encoding/json"
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
					"did:example:alice",
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
		t.Logf("%v", test.MustProtoToJson(response))
		require.True(t, response.Result.Ok)
	}
}

func Test_ShinzoPolicy_Ok(t *testing.T) {
	in := `{
  "policyDefinition": "name: shinzo\nresources:\n  primitive:\n\t\trelations:\n\t\t\tadmin:\n\t\t\t\tmanages:\n\t\t\t\t- writer\n\t\t\t\t- reader\n\t\t\t\t- banned\n\t\t\t\ttypes:\n\t\t\t\t- actor\n\t\t\t\t- group->administrator\n\t\t\twriter:\n\t\t\t\ttypes:\n\t\t\t\t- actor\n\t\t\t\t- group->member\n\t\t\treader:\n\t\t\t\ttypes:\n\t\t\t\t- actor\n\t\t\t\t- group->member\n\t\t\tsubscriber:\n\t\t\t\ttypes:\n\t\t\t\t- actor\n\t\t\tbanned:\n\t\t\t\ttypes:\n\t\t\t\t- actor\n\t\tpermissions:\n\t\t\tread:\n\t\t\t\texpr: (writer + reader) - banned\n\t\t\tupdate:\n\t\t\t\texpr: writer - banned\n\t\t\tdelete:\n\t\t\t\texpr: owner\n\t\t\tquery:\n\t\t\t\texpr: (subscriber) - banned\n  view:\n    relations:\n      admin:\n        manages:\n        - writer\n        - reader\n        - subscriber\n        - banned\n        types:\n        - actor\n        - group->administrator\n      creator:\n        types:\n        - actor\n      writer:\n        types:\n        - actor\n        - group->member\n      reader:\n        types:\n        - actor\n        - group->member\n      banned:\n        types:\n        - actor\n      subscriber:\n        types:\n        - actor\n      parent:\n        types:\n        - primitive\n        - view\n    permissions:\n      read:\n        expr: writer + reader + parent->read + subscriber - banned\n      update:\n        expr: ((writer + parent->update) & parent->read) - banned\n      query:\n        expr: (subscriber) - banned\n      delete:\n        expr: owner\n  group:  \n    relations:\n      owner:\n        manages:\n        - admin\n        - guest\n        - blocked\n        types:\n        - actor\n      admin:\n        manages:\n        - guest\n        - blocked\n        types:\n        - actor\n        - group->administrator\n      guest:\n        types:\n        - actor\n      blocked:\n        types:\n        - actor\n    permissions:\n      member:\n        expr: guest - blocked\n      administrator:\n        expr: (owner + admin) - blocked",
  "policyTheorem": "Authorizations {\n  // @@@ blocks tests\n  primitive:blocks#read@did:user:randomUser\n  !primitive:blocks#update@did:user:randomUser\n  !primitive:blocks#query@did:user:randomUser\n\n  primitive:blocks#read@did:user:aHost\n  !primitive:blocks#update@did:user:aHost\n  !primitive:blocks#query@did:user:aHost\n\n  primitive:blocks#read@did:user:anIndexer\n  primitive:blocks#update@did:user:anIndexer\n  !primitive:blocks#query@did:user:anIndexer\n\n  !primitive:blocks#read@did:user:duncan \n  !primitive:blocks#update@did:user:duncan\n  !primitive:blocks#query@did:user:duncan\n  !primitive:blocks#read@did:user:addo\n  !primitive:blocks#update@did:user:addo\n  !primitive:blocks#query@did:user:addo\n  !primitive:blocks#read@did:user:quinn\n  !primitive:blocks#update@did:user:quinn\n  !primitive:blocks#query@did:user:quinn\n  !primitive:blocks#read@did:user:daniel\n  !primitive:blocks#update@did:user:daniel\n  !primitive:blocks#query@did:user:daniel\n\n  !primitive:blocks#read@did:user:shinzohub\n  !primitive:blocks#update@did:user:shinzohub\n  !primitive:blocks#query@did:user:shinzohub\n\n  !primitive:blocks#read@did:user:unregisteredUser\n  !primitive:blocks#update@did:user:unregisteredUser\n  !primitive:blocks#query@did:user:unregisteredUser\n\n  !primitive:blocks#read@did:user:aBlockedIndexer\n  !primitive:blocks#update@did:user:aBlockedIndexer\n  !primitive:blocks#query@did:user:aBlockedIndexer\n\n  !primitive:blocks#read@did:user:aBannedIndexer\n  !primitive:blocks#update@did:user:aBannedIndexer\n  !primitive:blocks#query@did:user:aBannedIndexer\n\n  !primitive:blocks#read@did:user:aBlockedHost\n  !primitive:blocks#update@did:user:aBlockedHost\n  !primitive:blocks#query@did:user:aBlockedHost\n\n  // @@@ Secondary data feed (a data feed derived directly from primitives) tests\n  view:datafeedA#read@did:user:subscriber\n  !view:datafeedA#update@did:user:subscriber\n  view:datafeedA#query@did:user:subscriber\n  !view:datafeedA#creator@did:user:subscriber\n\n  !view:datafeedA#read@did:user:subscriberToB\n  !view:datafeedA#update@did:user:subscriberToB\n  !view:datafeedA#query@did:user:subscriberToB\n  !view:datafeedA#creator@did:user:subscriberToB\n\n  !view:datafeedA#read@did:user:anonsubscriber\n  !view:datafeedA#update@did:user:anonsubscriber\n  !view:datafeedA#query@did:user:anonsubscriber\n  !view:datafeedA#creator@did:user:anonsubscriber\n\n  view:datafeedA#read@did:user:aHost\n  view:datafeedA#update@did:user:aHost\n  !view:datafeedA#query@did:user:aHost\n  !view:datafeedA#creator@did:user:aHost\n\n  !view:datafeedA#read@did:user:creator\n  !view:datafeedA#update@did:user:creator\n  !view:datafeedA#query@did:user:creator\n  view:datafeedA#creator@did:user:creator\n\n  !view:datafeedA#read@did:user:creatorOfB\n  !view:datafeedA#update@did:user:creatorOfB\n  !view:datafeedA#query@did:user:creatorOfB\n  !view:datafeedA#creator@did:user:creatorOfB\n\n  !view:datafeedA#read@did:user:addo\n  !view:datafeedA#update@did:user:addo\n  !view:datafeedA#query@did:user:addo\n  !view:datafeedA#creator@did:user:addo\n\n  !view:datafeedA#read@did:user:shinzohub\n  !view:datafeedA#update@did:user:shinzohub\n  !view:datafeedA#query@did:user:shinzohub\n  !view:datafeedA#creator@did:user:shinzohub\n\n  !view:datafeedA#read@did:user:aBannedHost\n  !view:datafeedA#update@did:user:aBannedHost\n  !view:datafeedA#query@did:user:aBannedHost\n  !view:datafeedA#creator@did:user:aBannedHost\n\n  !view:datafeedA#read@did:user:aBlockedHost\n  !view:datafeedA#update@did:user:aBlockedHost\n  !view:datafeedA#query@did:user:aBlockedHost\n  !view:datafeedA#creator@did:user:aBlockedHost\n\n  view:datafeedA#read@did:user:anIndexer\n  view:datafeedA#update@did:user:anIndexer\n  !view:datafeedA#query@did:user:anIndexer\n  !view:datafeedA#creator@did:user:anIndexer\n\n  // @@@ Tertiary data feed (a datafeed derived directly from another data feed) tests\n  view:datafeedB#read@did:user:subscriber\n  !view:datafeedB#update@did:user:subscriber\n  !view:datafeedB#query@did:user:subscriber\n  !view:datafeedB#creator@did:user:subscriber\n\n  view:datafeedB#read@did:user:subscriberToB\n  !view:datafeedB#update@did:user:subscriberToB\n  view:datafeedB#query@did:user:subscriberToB\n  !view:datadeefB#creator@did:user:subscriberToB\n\n  !view:datafeedB#read@did:user:anonsubscriber\n  !view:datafeedB#update@did:user:anonsubscriber\n  !view:datafeedB#query@did:user:anonsubscriber\n  !view:datafeedB#creator@did:user:anonsubscriber\n\n  // Members of the host group inherit write permissions on datafeedB as they have it for datafeedA and datafeedA is datafeedB's parent\n  view:datafeedB#read@did:user:aHost\n  view:datafeedB#update@did:user:aHost\n  !view:datafeedB#query@did:user:aHost\n  !view:datafeedB#creator@did:user:aHost\n\n  !view:datafeedB#read@did:user:creator\n  !view:datafeedB#update@did:user:creator\n  !view:datafeedB#query@did:user:creator\n  !view:datafeedB#creator@did:user:creator\n\n  !view:datafeedB#read@did:user:creatorOfB\n  !view:datafeedB#update@did:user:creatorOfB\n  !view:datafeedB#query@did:user:creatorOfB\n  view:datafeedB#creator@did:user:creatorOfB\n\n  !view:datafeedB#read@did:user:addo\n  !view:datafeedB#update@did:user:addo\n  !view:datafeedB#query@did:user:addo\n  !view:datafeedB#creator@did:user:addo\n\n  !view:datafeedB#read@did:user:shinzohub\n  !view:datafeedB#update@did:user:shinzohub\n  !view:datafeedB#query@did:user:shinzohub\n  !view:datafeedB#creator@did:user:shinzohub\n\n  !view:datafeedB#read@did:user:aBannedHost\n  !view:datafeedB#update@did:user:aBannedHost\n  !view:datafeedB#query@did:user:aBannedHost\n  !view:datafeedB#creator@did:user:aBannedHost\n\n  !view:datafeedB#read@did:user:aBlockedHost\n  !view:datafeedB#update@did:user:aBlockedHost\n  !view:datafeedB#query@did:user:aBlockedHost\n  !view:datafeedB#creator@did:user:aBlockedHost\n\n  view:datafeedB#read@did:user:anIndexer\n  view:datafeedB#update@did:user:anIndexer\n  !view:datafeedB#query@did:user:anIndexer\n  !view:datafeedB#creator@did:user:anIndexer\n\n}\n\nDelegations {\n  // Shinzo Team group administration tests\n  !did:user:quinn > group:indexer#guest\n  !did:user:quinn > group:indexer#admin\n  !did:user:quinn > group:indexer#owner\n  did:user:duncan > group:indexer#guest\n  !did:user:duncan > group:indexer#admin\n  !did:user:duncan > group:indexer#owner\n  did:user:addo > group:indexer#guest\n  !did:user:addo > group:indexer#admin\n  !did:user:addo > group:indexer#owner\n\n  !did:user:quinn > group:host#guest\n  !did:user:quinn > group:host#admin\n  !did:user:quinn > group:host#owner\n  did:user:duncan > group:host#guest\n  !did:user:duncan > group:host#admin\n  !did:user:duncan > group:host#owner\n  did:user:addo > group:host#guest\n  !did:user:addo > group:host#admin\n  !did:user:addo > group:host#owner\n\n  !did:user:quinn > group:shinzoteam#guest\n  !did:user:quinn > group:shinzoteam#admin\n  !did:user:quinn > group:shinzoteam#owner\n  did:user:duncan > group:shinzoteam#guest\n  !did:user:duncan > group:shinzoteam#admin\n  !did:user:duncan > group:shinzoteam#owner\n  did:user:addo > group:shinzoteam#guest\n  did:user:addo > group:shinzoteam#admin\n  did:user:addo > group:shinzoteam#owner \n\n  // Shinzo team file administration \n  !did:user:addo > primitive:blocks#reader\n  !did:user:addo > primitive:blocks#writer\n  !did:user:addo > primitive:blocks#admin\n  !did:user:addo > primitive:blocks#owner\n\n  did:user:addo > view:datafeedA#reader\n  did:user:addo > view:datafeedA#writer\n  !did:user:addo > view:datafeedA#creator\n  !did:user:addo > view:datafeedA#admin\n  !did:user:addo > view:datafeedA#owner\n\n  did:user:addo > view:datafeedB#reader\n  did:user:addo > view:datafeedB#writer\n  !did:user:addo > view:datafeedB#creator\n  !did:user:addo > view:datafeedB#admin\n  !did:user:addo > view:datafeedB#owner\n\n  !did:user:duncan > primitive:blocks#reader\n  !did:user:duncan > primitive:blocks#writer\n  !did:user:duncan > primitive:blocks#admin\n  !did:user:duncan > primitive:blocks#owner\n\n  did:user:duncan > view:datafeedA#reader\n  did:user:duncan > view:datafeedA#writer\n  !did:user:duncan > view:datafeedA#creator\n  !did:user:duncan > view:datafeedA#admin\n  !did:user:duncan > view:datafeedA#owner\n\n  did:user:duncan > view:datafeedB#reader\n  did:user:duncan > view:datafeedB#writer\n  !did:user:duncan > view:datafeedB#creator\n  !did:user:duncan > view:datafeedB#admin\n  !did:user:duncan > view:datafeedB#owner\n\n  !did:user:quinn > primitive:blocks#reader\n  !did:user:quinn > primitive:blocks#writer\n  !did:user:quinn > primitive:blocks#admin\n  !did:user:quinn > primitive:blocks#owner\n\n  !did:user:quinn > view:datafeedA#reader\n  !did:user:quinn > view:datafeedA#writer\n  !did:user:quinn > view:datafeedA#creator\n  !did:user:quinn > view:datafeedA#admin\n  !did:user:quinn > view:datafeedA#owner\n\n  // Shinzo hub access administration tests\n  did:user:shinzohub > group:indexer#guest\n  !did:user:shinzohub > group:indexer#admin\n  !did:user:shinzohub > group:indexer#owner\n\n  did:user:shinzohub > group:host#guest\n  !did:user:shinzohub > group:host#admin\n  !did:user:shinzohub > group:host#owner\n\n  did:user:shinzohub > primitive:blocks#reader\n  did:user:shinzohub > primitive:blocks#writer\n  !did:user:shinzohub > primitive:blocks#admin\n  !did:user:shinzohub > primitive:blocks#owner\n\n  did:user:shinzohub > view:datafeedA#reader\n  did:user:shinzohub > view:datafeedA#writer\n  did:user:shinzohub > view:datafeedA#subscriber\n  !did:user:shinzohub > view:datafeedA#creator\n  !did:user:shinzohub > view:datafeedA#admin\n  !did:user:shinzohub > view:datafeedA#owner\n\n  did:user:shinzohub > view:datafeedB#reader\n  did:user:shinzohub > view:datafeedB#writer\n  did:user:shinzohub > view:datafeedB#subscriber\n  !did:user:shinzohub > view:datafeedB#creator\n  !did:user:shinzohub > view:datafeedB#admin\n  !did:user:shinzohub > view:datafeedB#owner\n}",
  "relationships": "primitive:blocks#owner@did:user:sourcehub\nprimitive:blocks#admin@did:user:shinzohub\nprimitive:blocks#writer@group:indexer#member\nprimitive:blocks#reader@group:host#member\nprimitive:blocks#reader@did:user:randomUser\nprimitive:blocks#banned@did:user:aBannedIndexer\nprimitive:blocks#subscriber@did:user:subscriber\n\nview:datafeedA#owner@did:user:sourcehub\nview:datafeedA#admin@did:user:shinzohub\nview:datafeedA#admin@group:shinzoteam#administrator\nview:datafeedA#creator@did:user:creator\nview:datafeedA#subscriber@did:user:subscriber\nview:datafeedA#writer@group:host#member\nview:datafeedA#banned@did:user:aBannedHost\nview:datafeedA#parent@primitive:blocks\n\nview:datafeedB#owner@did:user:sourcehub\nview:datafeedB#admin@did:user:shinzohub\nview:datafeedB#admin@group:shinzoteam#administrator\nview:datafeedB#creator@did:user:creatorOfB\nview:datafeedB#subscriber@did:user:subscriberToB\nview:datafeedB#parent@view:datafeedA\n\ngroup:shinzoteam#owner@did:user:addo \ngroup:shinzoteam#admin@did:user:duncan\ngroup:shinzoteam#guest@did:user:quinn\ngroup:shinzoteam#guest@did:user:daniel\n\ngroup:indexer#owner@did:user:sourcehub\ngroup:indexer#admin@did:user:shinzohub\ngroup:indexer#admin@group:shinzoteam#administrator\ngroup:indexer#guest@did:user:anIndexer\ngroup:indexer#guest@did:user:aBlockedIndexer\ngroup:indexer#blocked@did:user:aBlockedIndexer\ngroup:indexer#guest@did:user:aBannedIndexer\n\ngroup:host#owner@did:user:sourcehub\ngroup:host#admin@did:user:shinzohub\ngroup:host#admin@group:shinzoteam#administrator\ngroup:host#guest@did:user:aHost\ngroup:host#guest@did:user:aBlockedHost\ngroup:host#blocked@did:user:aBlockedHost"
}`
	d := make(map[string]string)
	err := json.Unmarshal([]byte(in), &d)
	require.NoError(t, err)

	ctx := test.NewTestCtx(t)
	action := NewAndSet{
		Data: &types.SandboxData{
			PolicyDefinition: d["policyDefinition"],
			PolicyTheorem:    d["policyTheorem"],
			Relationships:    d["relationships"],
		},
	}
	handle := action.Run(ctx)

	result, err := ctx.Playground.VerifyTheorems(ctx, &types.VerifyTheoremsRequest{
		Handle: handle,
	})
	require.NoError(t, err)
	_ = result

	require.True(t, result.Result.Ok)
	require.Zero(t, result.Result.Failures)
}

func Test_DOTExplainCheck_Example(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a := NewAndSet{
		Data: sandbox.Samples[0].Data,
	}
	handle := a.Run(ctx)

	resp, err := ctx.Playground.DOTExplainCheck(ctx, &types.DOTExplainCheckRequest{
		Handle:     handle,
		Object:     types.NewObject("file", "def"),
		Permission: "read",
		Actor:      types.NewActor("did:user:eve"),
	})
	require.NoError(t, err)
	t.Log(resp.DotGraph)
	require.NotNil(t, resp)
}

func Test_ExplainCheck_Example(t *testing.T) {
	ctx := test.NewTestCtx(t)

	a := NewAndSet{
		Data: sandbox.Samples[0].Data,
	}
	handle := a.Run(ctx)

	resp, err := ctx.Playground.ExplainCheck(ctx, &types.ExplainCheckRequest{
		Handle:     handle,
		Object:     types.NewObject("file", "def"),
		Permission: "read",
		Actor:      types.NewActor("did:user:eve"),
	})
	require.NoError(t, err)
	t.Log(resp.Graph.String())
	require.NotNil(t, resp.Graph)
}
