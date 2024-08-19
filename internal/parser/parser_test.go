package parser

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestParseRelationship_EmptyString_ReturnsError(t *testing.T) {
	//relationship, err := ParseRelationship("")

	//require.Error(t, err)
	//require.Nil(t, relationship)
}

func TestParseRelationship_RelationshipPlusJunk_ReturnsError(t *testing.T) {
	relationship, err := ParseRelationship("file:abc#owner@did:example:bob aksldfskwer#junk")

	require.Error(t, err)
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

func TestParseRelationships_RelationshipSetWithTrailingData_Errors(t *testing.T) {
	relationships := `
	resource:abc#relation@did:example:bob
	resource:abc#relation@resource:thing
	resource:abc#relation@resource:userset#member

	abc1234
	`
	rels, err := ParseRelationships(relationships)

	require.Error(t, err)
	require.Nil(t, rels)
}

func TestParseRelationships_EmptySuiteGetsParsed(t *testing.T) {
	relationships := ""

	rels, err := ParseRelationships(relationships)

	require.Nil(t, err)
	require.Len(t, rels, 0)
}

func TestPolicyTheorem_ParsesCorrectly(t *testing.T) {
	theorem := `
	Authorizations {
      note:abc#owner@did:example:bob
      !note:abc#owner@did:example:alice //this is a comment which extensd until the end of the line
	}

	Delegations {
	  did:ex:bob > note:abc#read
	  ! did:ex:bob > note:abc#read
	}
	`

	thm, err := ParsePolicyTheorem(theorem)

	require.Nil(t, err)
	want := &types.PolicyTheorem{
		AuthorizationTheorems: []*types.AuthorizationTheorem{
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
	require.Equal(t, want, thm.ToPolicyTheorem())
}

func TestPolicyTheorem_TheoremWithTrailingInput_Errors(t *testing.T) {
	theorem := `
	Authorizations { }

	Delegations { }

	trailniig-data-abc
	`

	thm, err := ParsePolicyTheorem(theorem)
	require.Error(t, err)
	require.Nil(t, thm)
}
