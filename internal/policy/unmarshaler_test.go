package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func TestFullUnmarshal(t *testing.T) {
	in := `
description: ok
name: policy
resources:
- name: foo
  permissions:
  - doc: abc doc
    name: abc
  - expr: abc
    name: def
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name:      "foo",
				Relations: []*types.Relation{},
				Permissions: []*types.Permission{
					{
						Name:       "abc",
						Doc:        "abc doc",
						Expression: "",
					},
					{
						Name:       "def",
						Doc:        "",
						Expression: "abc",
					},
				},
			},
		},
		ActorResource: nil,
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestUnmarshalWithoutSpecDefaultsToNone(t *testing.T) {
	in := `name: policy
`
	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		Name:              "policy",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestEmptyResourceMapsToResource(t *testing.T) {
	in := `
resources:
- name: foo
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name:        "foo",
				Relations:   []*types.Relation{},
				Permissions: []*types.Permission{},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestResourceWithoutPermsOrRelsMapsToResource(t *testing.T) {
	in := `
resources:
- name: foo
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name:        "foo",
				Permissions: []*types.Permission{},
				Relations:   []*types.Relation{},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestEmptyRelationMapsToRelation(t *testing.T) {
	in := `
resources:
- name: foo
  relations:
  - name: blah
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Relations: []*types.Relation{
					{
						Name:    "blah",
						VrTypes: []*types.Restriction{},
					},
				},
				Permissions: []*types.Permission{},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestEmptyPermissionMapsToPermission(t *testing.T) {
	// NOTE The purpose of this test is to assert that the values are
	// correctly unmarshaled.
	// Therefore, even though a permission requires an expression,
	// it's ok because the validation will happen elsewhere.
	// Asserting the type unmarhsals correctly means that the validator -
	// as opposed to the unmarshaler - will error out leading to better error msgs.
	in := `
resources:
- name: foo
  permissions:
  - name: blah
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Permissions: []*types.Permission{
					{
						Name: "blah",
					},
				},
				Relations: []*types.Relation{},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestDuplicatedResourceErrors(t *testing.T) {
	in := `
    resources:
      foo:
      foo:
    `

	_, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NotNil(t, err)
}

func TestDuplicatedPermissionErrors(t *testing.T) {
	in := `
    resources:
      foo:
        permissions:
          read:
          read:
    `

	_, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NotNil(t, err)
}

func TestDuplicatedRelationErrors(t *testing.T) {
	in := `
    resources:
      foo:
        relations:
          reader:
          reader:
    `

	_, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NotNil(t, err)
}

func TestRestrictionIdentifierMapsBothForms(t *testing.T) {
	in := `
resources:
- name: foo
  relations:
  - name: blah
    types:
    - actor
    - book->owner
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Relations: []*types.Relation{
					{
						Name: "blah",
						VrTypes: []*types.Restriction{
							{
								ResourceName: "actor",
								RelationName: "",
							},
							{
								ResourceName: "book",
								RelationName: "owner",
							},
						},
					},
				},
				Permissions: []*types.Permission{},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func Test_NoneSpecMapsToNoneSpecficationType(t *testing.T) {
	in := `name: test
spec: none
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NoError(t, err)
	want := &types.Policy{
		Name:              "test",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Equal(t, want, out)
}

func Test_DefraSpecMapsToDefraSpecficationType(t *testing.T) {
	in := `name: test
spec: defra
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NoError(t, err)
	want := &types.Policy{
		Name:              "test",
		SpecificationType: types.PolicySpecificationType_DEFRA_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Equal(t, want, out)
}

func Test_GibberingSpecMapsErrors(t *testing.T) {
	in := `
	name: test
	spec: gibberish-1234
	`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.Error(t, err)
	require.Nil(t, out)
}

func Test_EmptySpecMapsToNone(t *testing.T) {
	in := `name: test
spec: ""
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	require.NoError(t, err)
	want := &types.Policy{
		Name:              "test",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Equal(t, want, out)
}

func TestYaml_FullUnmarshal(t *testing.T) {
	in := `
actor:
  doc: my actor
  name: actor-resource
description: ok
name: policy
resources:
- name: foo
  permissions:
  - doc: abc doc
    name: abc
  - expr: abc
    name: def
`

	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name:      "foo",
				Relations: []*types.Relation{},
				Permissions: []*types.Permission{
					{
						Name:       "abc",
						Doc:        "abc doc",
						Expression: "",
					},
					{
						Name:       "def",
						Doc:        "",
						Expression: "abc",
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Relations: []*types.Relation{},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestYaml_UnmarshalWithEmptyPermExpression(t *testing.T) {
	in := `
description: ok
name: policy
resources:
- name: foo
  permissions:
  - name: abc
`
	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name:      "foo",
				Relations: []*types.Relation{},
				Permissions: []*types.Permission{
					{
						Name: "abc",
					},
				},
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}
