package parser

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestParseRelationship_EmptyString(t *testing.T) {
	relationship, err := ParseRelationship("")

	require.Nil(t, err)
	require.Nil(t, relationship)
}

func TestParseRelationship_DIDActorRelationship(t *testing.T) {
	relationship, err := ParseRelationship("test:abc#rel@did:example:bob")

	require.Nil(t, err)
	require.Equal(t, relationship, types.NewActorRelationship("test", "abc", "rel", "did:example:bob"))
}

func TestParseRelationship_UsersetRelationship(t *testing.T) {
	relationship, err := ParseRelationship("test:abc#rel@group:blah#member")

	require.Nil(t, err)
	require.Equal(t, relationship, types.NewActorSetRelationship("test", "abc", "rel", "group", "blah", "member"))
}

func TestParseRelationship_ObjectSubject(t *testing.T) {
	relationship, err := ParseRelationship("test:abc#rel@file:test")

	require.Nil(t, err)
	require.Equal(t, relationship, types.NewRelationship("test", "abc", "rel", "file", "test"))
}

func TestParseRelationship_EscapedObjectParsedCorrectly(t *testing.T) {
	relationship, err := ParseRelationship(`test:"abc:#@1234"#rel@file:test`)

	require.Nil(t, err)
	require.Equal(t, relationship, types.NewRelationship("test", `abc:#@1234`, "rel", "file", "test"))
}

func TestParseRelationship_InvalidRelationshipError(t *testing.T) {
	relationship, err := ParseRelationship(`test:invalid#id#rel@file:test`)

	require.NotNil(t, err)
	require.Nil(t, relationship)
}

func TestParseRelationships(t *testing.T) {
	relationships := `
	resource:abc#relation@did:example:bob
	resource:abc#relation@resource:thing
	resource:abc#relation@resource:userset#member
	`

	rels, err := ParseRelationships(relationships)

	require.Nil(t, err)
	want := []*types.Relationship{
		types.NewActorRelationship("resource", "abc", "relation", "did:example:bob"),
		types.NewRelationship("resource", "abc", "relation", "resource", "thing"),
		types.NewActorSetRelationship("resource", "abc", "relation", "resource", "userset", "member"),
	}
	require.Equal(t, want, rels)
}
func TestParseTestSuite_EmptySuiteGetsParsed(t *testing.T) {
	relationships := ""

	rels, err := ParseRelationships(relationships)

	require.Nil(t, err)
	want := make([]*types.Relationship, 0)
	require.Equal(t, want, rels)
}

func TestPolicyTheorem_ParsesCorrectly(t *testing.T) {
	theorem := `
	Authorizations {
      note:abc#owner@did:example:bob
      !note:abc#owner@did:example:alice
	}

	Delegations {
	  did:ex:bob > note:abc#read
	  ! did:ex:bob > note:abc#read
	}
	`

	thm, err := ParsePolicyTheorem(theorem)

	require.Nil(t, err)
	want := &types.PolicyTheorem{
		AuthorizationThereoms: []*types.AuthorizationTheorem{
			{
				Operation:  types.NewOperation(types.NewObject("note", "abc"), "owner"),
				Actor:      types.NewActor("did:example:bob"),
				AssertTrue: true,
			},
			{
				Operation:  types.NewOperation(types.NewObject("note", "abc"), "owner"),
				Actor:      types.NewActor("did:example:alice"),
				AssertTrue: false,
			},
		},
		DelegationTheorems: []*types.DelegationTheorem{
			{
				Actor:      types.NewActor("did:ex:bob"),
				Operation:  types.NewOperation(types.NewObject("note", "abc"), "read"),
				AssertTrue: true,
			},
			{
				Actor:      types.NewActor("did:ex:bob"),
				Operation:  types.NewOperation(types.NewObject("note", "abc"), "read"),
				AssertTrue: false,
			},
		},
	}
	require.Equal(t, want, thm)
}
