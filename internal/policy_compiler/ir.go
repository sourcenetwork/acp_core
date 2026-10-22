// package policy_compiler lowers a processed acp_core Policy into a relational model
// which can be emitted to a target language (eg. Alloy).
//
// Every relation and permission of a resource becomes a Computation:
// an expression which evaluates to the set of actors related to an object.
// Relation tuples are stored in Fields, one Field per relation restriction.
package policy_compiler

// Model is the intermediate representation of a Policy
type Model struct {
	Name      string
	Resources []*Resource
}

// Resource models a policy resource.
type Resource struct {
	Name string
	// Fields stores the tuples of the resource relations, one Field per relation restriction
	Fields []*Field
	// Relations evaluates each relation (including owner) into actors, expanding usersets
	Relations []*Computation
	// Permissions evaluates each permission into actors, from the effective permission expression
	Permissions []*Computation
	// Authorities evaluates, for each relation, the actors which may manage it
	Authorities []*Computation
	// Influencers evaluates, for each permission,
	// the actors which may change its value by managing some relation it depends on
	Influencers []*Computation
}

// Field models the stored tuples of a relation for one of its restrictions.
type Field struct {
	Name     string
	Resource string
	Relation string
	// Target is either ActorTarget or the name of a resource
	Target string
	// Userset is set for userset restrictions (eg. group->member),
	// and points to the expanded computation in the Target resource
	Userset *Computation
	// One is true for fields which hold exactly one tuple per object (ie. owner)
	One bool
}

// ActorTarget is the Field target for actor restrictions
const ActorTarget = "actor"

// Kind identifies what a Computation models
type Kind int

const (
	KindRelation Kind = iota
	KindPermission
	KindAuthority
	KindInfluence
)

// Computation is a named expression over a resource
// which evaluates to a set of actors for each object
type Computation struct {
	Kind     Kind
	Resource string
	// Name is the name of the relation or permission the computation refers to
	Name string
	Expr Expr
	// Source is the policy expression the computation was lowered from, if any
	Source string
}

// Expr is an expression which evaluates to a set of actors for each object of a resource
type Expr interface {
	isExpr()
}

// Tuples evaluates to the actors stored directly in Field
type Tuples struct {
	Field *Field
}

// Ref evaluates another Computation on the same object
type Ref struct {
	To *Computation
}

// Join evaluates Expr on every object reached through Field
type Join struct {
	Field *Field
	Expr  Expr
}

// Closure evaluates Local on the object and on every object
// transitively reached through Fields (ie. a recursive tuple to userset)
type Closure struct {
	Fields []*Field
	Local  Expr
}

// Op is a set operator
type Op int

const (
	Union Op = iota
	Difference
	Intersection
)

// Binary combines two expressions with a set operator
type Binary struct {
	Op    Op
	Left  Expr
	Right Expr
}

// Empty evaluates to no actors
type Empty struct{}

func (*Tuples) isExpr()  {}
func (*Ref) isExpr()     {}
func (*Join) isExpr()    {}
func (*Closure) isExpr() {}
func (*Binary) isExpr()  {}
func (*Empty) isExpr()   {}

// union combines exprs with Union, skipping Empty expressions
func union(exprs ...Expr) Expr {
	var expr Expr = &Empty{}
	for _, e := range exprs {
		switch {
		case isEmpty(e):
		case isEmpty(expr):
			expr = e
		default:
			expr = &Binary{Op: Union, Left: expr, Right: e}
		}
	}
	return expr
}

// unionArms flattens a tree of unions into its operands
func unionArms(e Expr) []Expr {
	if b, ok := e.(*Binary); ok && b.Op == Union {
		return append(unionArms(b.Left), unionArms(b.Right)...)
	}
	return []Expr{e}
}

func isEmpty(e Expr) bool {
	_, ok := e.(*Empty)
	return ok
}
