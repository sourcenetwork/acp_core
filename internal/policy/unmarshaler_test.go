package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func TestFullUnmarshal(t *testing.T) {
	in := `
        name: policy
        description: ok
		spec: none
        resources:
          foo:
            relations: 
              owner:
                doc: owner owns
                types:
                  - blah
                  - ok->that
                manages: 
                  - whatever
            permissions: 
              abc:
                expr: owner
                doc: abc doc
              def: 
                expr: owner + abc
        actor:
          name: actor-resource
          doc: my actor
          `
	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Relations: []*types.Relation{
					{
						Name: "owner",
						Doc:  "owner owns",
						VrTypes: []*types.Restriction{
							{
								ResourceName: "blah",
								RelationName: "",
							},
							{
								ResourceName: "ok",
								RelationName: "that",
							},
						},
						Manages: []string{
							"whatever",
						},
					},
				},
				Permissions: []*types.Permission{
					{
						Name:       "abc",
						Doc:        "abc doc",
						Expression: "owner",
					},
					{
						Name:       "def",
						Doc:        "",
						Expression: "owner + abc",
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Name: "actor-resource",
			Doc:  "my actor",
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestUnmarshalWithoutSpecDefaultsToUnkown(t *testing.T) {
	in := `
	name: policy
`
	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		Name:              "policy",
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestEmptyResourceMapsToResource(t *testing.T) {
	in := `
    resources:
      foo:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestResourceWithoutPermsOrRelsMapsToResource(t *testing.T) {
	in := `
    resources:
      foo:
        relations:
        permissions:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
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
      foo:
        relations:
          blah:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Relations: []*types.Relation{
					{
						Name: "blah",
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
      foo:
        permissions:
          blah:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
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

	unmarshaler := shortUnmarshaler{}
	_, err := unmarshaler.UnmarshalYAML(in)

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

	unmarshaler := shortUnmarshaler{}
	_, err := unmarshaler.UnmarshalYAML(in)

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

	unmarshaler := shortUnmarshaler{}
	_, err := unmarshaler.UnmarshalYAML(in)

	require.NotNil(t, err)
}

func TestRestrictionIdentifierMapsBothForms(t *testing.T) {
	in := `
    resources:
      foo:
        relations:
          blah:
            types:
              - actor
              - book->owner
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
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
	in := `
	name: test
	spec: none
	`

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	require.NoError(t, err)
	want := &types.Policy{
		Name:              "test",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Equal(t, want, out)
}

func Test_DefraSpecMapsToDefraSpecficationType(t *testing.T) {
	in := `
	name: test
	spec: defra
	`

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

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

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	require.Error(t, err)
	require.Nil(t, out)
}

func Test_EmptySpecMapsToUnknown(t *testing.T) {
	in := `
	name: test
	spec: ""
	`

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	require.NoError(t, err)
	want := &types.Policy{
		Name:              "test",
		SpecificationType: types.PolicySpecificationType_UNKNOWN_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Equal(t, want, out)
}

func TestYaml_FullUnmarshal(t *testing.T) {
	in := `name: policy
description: ok
spec: none
resources:
- name: foo
  relations: 
  - name: owner
    doc: owner owns
    types:
    - blah
    - ok->that
    manages: 
    - whatever
  permissions: 
  - name: abc
    expr: owner
    doc: abc doc
  - name: def
    expr: owner + abc
actor:
  name: actor-resource
  doc: my actor
`
	out, err := Unmarshal(in, types.PolicyMarshalingType_YAML)

	want := &types.Policy{
		Name:              "policy",
		Description:       "ok",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
				Relations: []*types.Relation{
					{
						Name: "owner",
						Doc:  "owner owns",
						VrTypes: []*types.Restriction{
							{
								ResourceName: "blah",
								RelationName: "",
							},
							{
								ResourceName: "ok",
								RelationName: "that",
							},
						},
						Manages: []string{
							"whatever",
						},
					},
				},
				Permissions: []*types.Permission{
					{
						Name:       "abc",
						Doc:        "abc doc",
						Expression: "owner",
					},
					{
						Name:       "def",
						Doc:        "",
						Expression: "owner + abc",
					},
				},
			},
		},
		ActorResource: &types.ActorResource{
			Name:      "actor-resource",
			Doc:       "my actor",
			Relations: []*types.Relation{},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestYaml_UnmarshalWithEmptyPermExpression(t *testing.T) {
	in := `name: policy
description: ok
spec: none
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
