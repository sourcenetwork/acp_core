package policy_compiler

import (
	"fmt"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/internal/policy/ppp"
	"github.com/sourcenetwork/acp_core/pkg/errors"
	parser "github.com/sourcenetwork/acp_core/pkg/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var ErrPolicyCompiler = errors.New("policy compiler", errors.ErrorType_BAD_INPUT)

// Compile unmarshals and processes a policy definition, lowering it into a Model
func Compile(definition string, marshalType types.PolicyMarshalingType) (*Model, error) {
	pol, err := policy.Unmarshal(definition, marshalType)
	if err != nil {
		return nil, errors.NewWithCause("unmarshaling policy", err, errors.ErrorType_BAD_INPUT)
	}

	// sequence is irrelevant since the policy is never persisted
	pipeline := ppp.CreatePolicyPipelineFactory(0, pol.SpecificationType)
	result, err := pipeline.Process(pol)
	if err != nil {
		return nil, err
	}

	return Lower(&result.Policy)
}

// Lower builds a Model from a processed Policy.
//
// Lowering uses the effective permission expressions and the management rules,
// which means pol must have gone through the policy processing pipeline.
func Lower(pol *types.Policy) (*Model, error) {
	l := &lowering{
		actorResource: pol.ActorResource.Name,
		fields:        make(map[string]map[string][]*Field),
		values:        make(map[string]map[string]*Computation),
		errs:          errors.NewMultiError(ErrPolicyCompiler),
	}
	model := &Model{Name: pol.Name}

	// computations are declared upfront so fields and expressions can reference them
	for _, res := range pol.Resources {
		r := &Resource{Name: res.Name}
		for _, rel := range relations(res) {
			r.Relations = append(r.Relations, l.declare(KindRelation, res.Name, rel.Name))
		}
		for _, perm := range res.Permissions {
			r.Permissions = append(r.Permissions, l.declare(KindPermission, res.Name, perm.Name))
		}
		model.Resources = append(model.Resources, r)
	}

	for i, res := range pol.Resources {
		r := model.Resources[i]
		for _, rel := range relations(res) {
			r.Fields = append(r.Fields, l.lowerFields(res.Name, rel)...)
		}
	}

	for i, res := range pol.Resources {
		r := model.Resources[i]
		for j, rel := range relations(res) {
			r.Relations[j].Expr = l.relationExpr(res.Name, rel.Name)
		}
		for j, perm := range res.Permissions {
			r.Permissions[j].Expr = l.parseExpr(res.Name, perm.EffectiveExpression)
			r.Permissions[j].Source = perm.EffectiveExpression
		}
		for _, rule := range res.ManagementRules {
			r.Authorities = append(r.Authorities, &Computation{
				Kind:     KindAuthority,
				Resource: res.Name,
				Name:     rule.Relation,
				Expr:     l.parseExpr(res.Name, rule.Expression),
				Source:   rule.Expression,
			})
		}
	}
	if len(l.errs.GetErrors()) > 0 {
		return nil, l.errs
	}

	l.resolveRecursion(model)
	if len(l.errs.GetErrors()) > 0 {
		return nil, l.errs
	}

	addInfluencers(model)
	return model, nil
}

type lowering struct {
	actorResource string
	// fields indexes fields by resource and relation
	fields map[string]map[string][]*Field
	// values indexes relation and permission computations by resource and name
	values map[string]map[string]*Computation
	errs   *errors.MultiError
}

// relations returns the resource relations, including owner
func relations(res *types.Resource) []*types.Relation {
	return append([]*types.Relation{res.Owner}, res.Relations...)
}

func (l *lowering) declare(kind Kind, resource, name string) *Computation {
	c := &Computation{Kind: kind, Resource: resource, Name: name}
	if l.values[resource] == nil {
		l.values[resource] = make(map[string]*Computation)
	}
	l.values[resource][name] = c
	return c
}

// lowerFields creates one Field per relation restriction.
//
// Fields are named after the relation, restrictions are appended to the name
// only when needed to disambiguate between fields of the same relation.
func (l *lowering) lowerFields(resource string, rel *types.Relation) []*Field {
	var fields []*Field
	for _, restriction := range rel.VrTypes {
		f := &Field{
			Name:     rel.Name,
			Resource: resource,
			Relation: rel.Name,
			Target:   restriction.ResourceName,
			One:      rel.Name == ppp.OwnerRelationName,
		}
		if f.Target == l.actorResource {
			f.Target = ActorTarget
		} else if l.values[f.Target] == nil {
			l.errorf("resource %v: relation %v: restriction references undefined resource %v", resource, rel.Name, f.Target)
			continue
		}

		if restriction.RelationName != "" {
			f.Userset = l.values[f.Target][restriction.RelationName]
			if f.Userset == nil {
				l.errorf("resource %v: relation %v: restriction references undefined relation %v->%v",
					resource, rel.Name, restriction.ResourceName, restriction.RelationName)
				continue
			}
		}

		if f.Target != ActorTarget && len(rel.VrTypes) > 1 {
			f.Name += "_" + f.Target
			if restriction.RelationName != "" {
				f.Name += "_" + restriction.RelationName
			}
		}
		fields = append(fields, f)
	}

	if l.fields[resource] == nil {
		l.fields[resource] = make(map[string][]*Field)
	}
	for _, f := range fields {
		for _, others := range l.fields[resource] {
			for _, other := range others {
				if other.Name == f.Name {
					l.errorf("resource %v: relation %v: field name %v clashes with relation %v", resource, rel.Name, f.Name, other.Relation)
				}
			}
		}
	}
	l.fields[resource][rel.Name] = fields
	return fields
}

// relationExpr expands a relation into the actors stored in its fields,
// following usersets. Fields which store objects contribute no actors.
func (l *lowering) relationExpr(resource, relation string) Expr {
	var exprs []Expr
	for _, f := range l.fields[resource][relation] {
		switch {
		case f.Target == ActorTarget:
			exprs = append(exprs, &Tuples{Field: f})
		case f.Userset != nil:
			exprs = append(exprs, &Join{Field: f, Expr: &Ref{To: f.Userset}})
		}
	}
	return union(exprs...)
}

func (l *lowering) parseExpr(resource, expr string) Expr {
	tree, err := parser.Parse(expr)
	if err != nil {
		l.errorf("resource %v: parsing expression %q: %w", resource, expr, err)
		return &Empty{}
	}
	return l.lowerTree(resource, tree)
}

func (l *lowering) lowerTree(resource string, tree *types.PermissionFetchTree) Expr {
	if comb := tree.GetCombNode(); comb != nil {
		ops := map[types.Combinator]Op{
			types.Combinator_UNION:        Union,
			types.Combinator_DIFFERENCE:   Difference,
			types.Combinator_INTERSECTION: Intersection,
		}
		op, ok := ops[comb.Combinator]
		if !ok {
			l.errorf("resource %v: unknown combinator %v", resource, comb.Combinator)
			return &Empty{}
		}
		return &Binary{Op: op, Left: l.lowerTree(resource, comb.Left), Right: l.lowerTree(resource, comb.Right)}
	}

	operation := tree.GetOperation()
	switch {
	case operation.GetCu() != nil:
		name := operation.GetCu().Relation
		c := l.values[resource][name]
		if c == nil {
			l.errorf("resource %v: expression references undefined relation %v", resource, name)
			return &Empty{}
		}
		return &Ref{To: c}
	case operation.GetTtu() != nil:
		return l.lowerTTU(resource, operation.GetTtu())
	default:
		l.errorf("resource %v: unsupported expression operation %v", resource, operation)
		return &Empty{}
	}
}

func (l *lowering) lowerTTU(resource string, ttu *types.TupleToUsersetNode) Expr {
	fields, ok := l.fields[resource][ttu.LookupRelation]
	if !ok {
		l.errorf("resource %v: tuple to userset references undefined relation %v", resource, ttu.LookupRelation)
		return &Empty{}
	}
	if len(fields) != 1 || fields[0].Target == ActorTarget || fields[0].Userset != nil {
		l.errorf("resource %v: tuple to userset %v->%v: relation %v must have exactly one resource type restriction",
			resource, ttu.LookupRelation, ttu.ComputedUsersetRelation, ttu.LookupRelation)
		return &Empty{}
	}
	f := fields[0]
	c := l.values[f.Target][ttu.ComputedUsersetRelation]
	if c == nil {
		l.errorf("resource %v: tuple to userset %v->%v: resource %v has no relation %v",
			resource, ttu.LookupRelation, ttu.ComputedUsersetRelation, f.Target, ttu.ComputedUsersetRelation)
		return &Empty{}
	}
	return &Join{Field: f, Expr: &Ref{To: c}}
}

func (l *lowering) errorf(format string, args ...any) {
	l.errs.Append(fmt.Errorf(format, args...))
}
