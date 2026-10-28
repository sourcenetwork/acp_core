package policy_compiler

import (
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// Assumption is an application assumption lowered against a Model
type Assumption interface {
	isAssumption()
}

// Disjoint states that no pair of fields shares a tuple.
// Only fields which can hold the same subjects are paired.
type Disjoint struct {
	Resource  string
	Relations []string
	Pairs     [][2]*Field
}

// AtMostOne states that an object has at most one tuple across Fields
type AtMostOne struct {
	Resource string
	Relation string
	Fields   []*Field
}

// Acyclic states that following Fields from an object never leads back to it.
// Fields include objects and usersets of the same resource (eg. nested groups).
type Acyclic struct {
	Resource string
	Relation string
	Fields   []*Field
}

func (*Disjoint) isAssumption()  {}
func (*AtMostOne) isAssumption() {}
func (*Acyclic) isAssumption()   {}

// UnmarshalAssumptions parses a yaml assumptions document:
//
//	assumptions:
//	- disjoint: [doc.owner, doc.editor]
//	- at_most_one: file.parent
//	- acyclic: folder.parent
func UnmarshalAssumptions(definition string) ([]*types.SystemConstraint, error) {
	var doc struct {
		Assumptions []struct {
			Disjoint  []string `json:"disjoint"`
			AtMostOne string   `json:"at_most_one"`
			Acyclic   string   `json:"acyclic"`
		} `json:"assumptions"`
	}
	if err := yaml.UnmarshalStrict([]byte(definition), &doc); err != nil {
		return nil, errors.NewWithCause("unmarshaling assumptions", err, errors.ErrorType_BAD_INPUT)
	}

	errs := errors.NewMultiError(ErrPolicyCompiler)
	var assumptions []*types.SystemConstraint
	for i, entry := range doc.Assumptions {
		set := 0
		for _, isSet := range []bool{len(entry.Disjoint) > 0, entry.AtMostOne != "", entry.Acyclic != ""} {
			if isSet {
				set++
			}
		}
		if set != 1 {
			errs.Append(fmt.Errorf("assumption %v: expected exactly one of disjoint, at_most_one or acyclic", i+1))
			continue
		}

		switch {
		case len(entry.Disjoint) > 0:
			disjoint := &types.DisjointConstraint{}
			for _, ref := range entry.Disjoint {
				resource, relation, err := parseRelationRef(ref)
				if err != nil {
					errs.Append(fmt.Errorf("assumption %v: %w", i+1, err))
					continue
				}
				if disjoint.Resource != "" && disjoint.Resource != resource {
					errs.Append(fmt.Errorf("assumption %v: disjoint relations must belong to the same resource", i+1))
					continue
				}
				disjoint.Resource = resource
				disjoint.Relations = append(disjoint.Relations, relation)
			}
			assumptions = append(assumptions, &types.SystemConstraint{Constraint: &types.SystemConstraint_Disjoint{Disjoint: disjoint}})
		case entry.AtMostOne != "":
			resource, relation, err := parseRelationRef(entry.AtMostOne)
			if err != nil {
				errs.Append(fmt.Errorf("assumption %v: %w", i+1, err))
				continue
			}
			assumptions = append(assumptions, &types.SystemConstraint{Constraint: &types.SystemConstraint_AtMostOne{
				AtMostOne: &types.AtMostOneConstraint{Resource: resource, Relation: relation},
			}})
		case entry.Acyclic != "":
			resource, relation, err := parseRelationRef(entry.Acyclic)
			if err != nil {
				errs.Append(fmt.Errorf("assumption %v: %w", i+1, err))
				continue
			}
			assumptions = append(assumptions, &types.SystemConstraint{Constraint: &types.SystemConstraint_Acyclic{
				Acyclic: &types.AcyclicConstraint{Resource: resource, Relation: relation},
			}})
		}
	}
	if len(errs.GetErrors()) > 0 {
		return nil, errs
	}
	return assumptions, nil
}

// parseRelationRef splits a resource.relation reference
func parseRelationRef(ref string) (resource, relation string, err error) {
	resource, relation, ok := strings.Cut(ref, ".")
	if !ok || resource == "" || relation == "" {
		return "", "", fmt.Errorf("invalid relation reference %q: expected resource.relation", ref)
	}
	return resource, relation, nil
}

// LowerAssumptions resolves assumptions against the relations of a Model
func LowerAssumptions(model *Model, assumptions []*types.SystemConstraint) ([]Assumption, error) {
	errs := errors.NewMultiError(ErrPolicyCompiler)
	fieldsOf := func(resource, relation string) []*Field {
		for _, r := range model.Resources {
			if r.Name != resource {
				continue
			}
			var fields []*Field
			for _, f := range r.Fields {
				if f.Relation == relation {
					fields = append(fields, f)
				}
			}
			if len(fields) == 0 {
				errs.Append(fmt.Errorf("constraint references undefined relation %v.%v", resource, relation))
			}
			return fields
		}
		errs.Append(fmt.Errorf("constraint references undefined resource %v", resource))
		return nil
	}

	var lowered []Assumption
	for _, a := range assumptions {
		switch {
		case a.GetDisjoint() != nil:
			d := a.GetDisjoint()
			if len(d.Relations) < 2 {
				errs.Append(fmt.Errorf("disjoint %v: expected at least two relations", d.Resource))
				continue
			}
			disjoint := &Disjoint{Resource: d.Resource, Relations: d.Relations}
			fields := make([][]*Field, len(d.Relations))
			for i, relation := range d.Relations {
				fields[i] = fieldsOf(d.Resource, relation)
			}
			for i := range d.Relations {
				for j := i + 1; j < len(d.Relations); j++ {
					if len(fields[i]) == 0 || len(fields[j]) == 0 {
						continue
					}
					var pairs [][2]*Field
					for _, f1 := range fields[i] {
						for _, f2 := range fields[j] {
							if f1.Target == f2.Target && f1.Userset == f2.Userset {
								pairs = append(pairs, [2]*Field{f1, f2})
							}
						}
					}
					if len(pairs) == 0 {
						errs.Append(fmt.Errorf("disjoint %v: relations %v and %v never share a subject type",
							d.Resource, d.Relations[i], d.Relations[j]))
					}
					disjoint.Pairs = append(disjoint.Pairs, pairs...)
				}
			}
			lowered = append(lowered, disjoint)
		case a.GetAtMostOne() != nil:
			m := a.GetAtMostOne()
			lowered = append(lowered, &AtMostOne{Resource: m.Resource, Relation: m.Relation, Fields: fieldsOf(m.Resource, m.Relation)})
		case a.GetAcyclic() != nil:
			c := a.GetAcyclic()
			var fields []*Field
			for _, f := range fieldsOf(c.Resource, c.Relation) {
				// objects and usersets of the same resource both link an object to another of its kind
				if f.Target == c.Resource {
					fields = append(fields, f)
				}
			}
			if len(fields) == 0 {
				errs.Append(fmt.Errorf("acyclic %v.%v: relation never points at a %v, so it can't form cycles", c.Resource, c.Relation, c.Resource))
			}
			lowered = append(lowered, &Acyclic{Resource: c.Resource, Relation: c.Relation, Fields: fields})
		default:
			errs.Append(fmt.Errorf("empty constraint"))
		}
	}
	if len(errs.GetErrors()) > 0 {
		return nil, errs
	}
	return lowered, nil
}
