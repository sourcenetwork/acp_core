package policy

import (
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

const (
	V1_0 string = "1.0"
)

func Unmarshal(pol string, t types.PolicyMarshalingType) (*types.Policy, error) {
	var policy *types.Policy
	var err error

	switch t {
	case types.PolicyMarshalingType_SHORT_YAML:
		unmarshaler := shortUnmarshaler{}
		policy, err = unmarshaler.UnmarshalYAML(pol)
	case types.PolicyMarshalingType_SHORT_JSON:
		unmarshaler := shortUnmarshaler{}
		policy, err = unmarshaler.UnmarshalJSON(pol)
	default:
		err = ErrUnknownMarshalingType
	}
	if err != nil {
		return policy, fmt.Errorf("%w: %w", ErrUnmarshaling, err)
	}

	return policy, nil
}

// shortUnmarshaler is a container type for unmarshaling
// short policy definitions into acp's Policy type.
type shortUnmarshaler struct{}

const typeDivider string = "->"

// Unmarshal a YAML serialized PolicyShort definition
func (u *shortUnmarshaler) UnmarshalYAML(pol string) (*types.Policy, error) {
	// remove trailing
	pol = strings.ReplaceAll(pol, "\t", "    ")
	pol = strings.Trim(pol, "\n")
	// Strict returns error if any key is duplicated
	polBytes, err := yaml.YAMLToJSONStrict([]byte(pol))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidShortPolicy, err)
	}

	return u.UnmarshalJSON(string(polBytes))
}

// Unmarshal a JSON serialized PolicyShort definition
func (u *shortUnmarshaler) UnmarshalJSON(pol string) (*types.Policy, error) {
	polShort := types.PolicyShort{}

	err := json.Unmarshal([]byte(pol), &polShort)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidShortPolicy, err)
	}

	return u.mapPolShort(&polShort)
}

func (u *shortUnmarshaler) mapPolShort(pol *types.PolicyShort) (*types.Policy, error) {
	resources := make([]*types.Resource, 0, len(pol.Resources))
	for name, resource := range pol.Resources {
		mapped := u.mapResource(name, resource)
		resources = append(resources, mapped)
	}

	spec, err := u.mapSpec(pol.Spec)
	if err != nil {
		return nil, err
	}

	policy := &types.Policy{
		Name:              pol.Name,
		Description:       pol.Description,
		Attributes:        pol.Meta,
		Resources:         resources,
		ActorResource:     pol.Actor,
		SpecificationType: spec,
	}

	// sort to ensure unmarshaling tests are not flaky
	sorted := ppp.SortTransformer{}
	sortedPol, _ := sorted.Transform(*policy) // SortTransformer does not error
	return &sortedPol, nil
}

func (u *shortUnmarshaler) mapResource(name string, resource *types.ResourceShort) *types.Resource {
	if resource == nil {
		return &types.Resource{
			Name: name,
		}
	}

	perms := make([]*types.Permission, 0, len(resource.Permissions))
	for name, perm := range resource.Permissions {
		mapped := u.mapPermission(name, perm)
		perms = append(perms, mapped)
	}

	rels := make([]*types.Relation, 0, len(resource.Relations))
	for name, rel := range resource.Relations {
		mapped := u.mapRelation(name, rel)
		rels = append(rels, mapped)
	}

	return &types.Resource{
		Name:        name,
		Doc:         resource.Doc,
		Permissions: perms,
		Relations:   rels,
	}
}

func (u *shortUnmarshaler) mapRelation(name string, rel *types.RelationShort) *types.Relation {
	if rel == nil {
		return &types.Relation{
			Name: name,
		}
	}

	vrTypes := utils.MapSlice(rel.Types, func(typeStr string) *types.Restriction {
		return u.mapType(typeStr)
	})
	return &types.Relation{
		Name:    name,
		Doc:     rel.Doc,
		Manages: rel.Manages,
		VrTypes: vrTypes,
	}
}

func (u *shortUnmarshaler) mapType(typeStr string) *types.Restriction {
	resource, rel, _ := strings.Cut(typeStr, typeDivider)
	return &types.Restriction{
		ResourceName: resource,
		RelationName: rel,
	}
}

func (u *shortUnmarshaler) mapPermission(name string, entry *types.PermissionShort) *types.Permission {
	perm := &types.Permission{
		Name: name,
	}
	if entry != nil {
		perm.Doc = entry.Doc
		perm.Expression = entry.Expr
	}
	return perm
}

func (u *shortUnmarshaler) mapSpec(spec string) (types.PolicySpecificationType, error) {
	switch strings.ToLower(spec) {
	case "defra":
		return types.PolicySpecificationType_DEFRA_SPEC, nil
	case "none":
		return types.PolicySpecificationType_NO_SPEC, nil
	case "":
		return types.PolicySpecificationType_UNKNOWN_SPEC, nil
	default:
		return types.PolicySpecificationType_UNKNOWN_SPEC, errors.Wrap("invalid specification", errors.ErrorType_BAD_INPUT)
	}
}
