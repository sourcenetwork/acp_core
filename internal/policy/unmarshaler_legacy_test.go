package policy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func TestLegacy_FullUnmarshal(t *testing.T) {
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

func TestLegacy_UnmarshalWithoutSpecDefaultsToNone(t *testing.T) {
	in := `name: policy`
	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		Name:              "policy",
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources:         []*types.Resource{},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestLegacy_EmptyResourceMapsToResource(t *testing.T) {
	in := `
    resources:
      foo:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
		Resources: []*types.Resource{
			{
				Name: "foo",
			},
		},
	}
	require.Nil(t, err)
	require.Equal(t, want, out)
}

func TestLegacy_ResourceWithoutPermsOrRelsMapsToResource(t *testing.T) {
	in := `
    resources:
      foo:
        relations:
        permissions:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

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

func TestLegacy_EmptyRelationMapsToRelation(t *testing.T) {
	in := `
    resources:
      foo:
        relations:
          blah:
    `

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	want := &types.Policy{
		SpecificationType: types.PolicySpecificationType_NO_SPEC,
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

func TestLegacy_EmptyPermissionMapsToPermission(t *testing.T) {
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

func TestLegacy_DuplicatedResourceErrors(t *testing.T) {
	in := `
    resources:
      foo:
      foo:
    `

	unmarshaler := shortUnmarshaler{}
	_, err := unmarshaler.UnmarshalYAML(in)

	require.NotNil(t, err)
}

func TestLegacy_DuplicatedPermissionErrors(t *testing.T) {
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

func TestLegacy_DuplicatedRelationErrors(t *testing.T) {
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

func TestLegacy_RestrictionIdentifierMapsBothForms(t *testing.T) {
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

func TestLegacy__NoneSpecMapsToNoneSpecficationType(t *testing.T) {
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

func TestLegacy__DefraSpecMapsToDefraSpecficationType(t *testing.T) {
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

func TestLegacy__GibberingSpecMapsErrors(t *testing.T) {
	in := `
	name: test
	spec: gibberish-1234
	`

	unmarshaler := shortUnmarshaler{}
	out, err := unmarshaler.UnmarshalYAML(in)

	require.Error(t, err)
	require.Nil(t, out)
}

func TestLegacy__EmptySpecMapsToNone(t *testing.T) {
	in := `
	name: test
	spec: ""
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
