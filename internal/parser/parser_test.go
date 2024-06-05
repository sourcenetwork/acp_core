package parser

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/stretchr/testify/require"
)

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

func TestParseRelationships_(t *testing.T) {}

func TestParseTestSuite_SuiteGetsParsed(t *testing.T) {

}

func TestParseTestSuite_EmptySuiteGetsParsed(t *testing.T) {

}
