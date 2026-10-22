package policy_compiler

import "strings"

// resolveRecursion rewrites self recursive computations into closures
// and rejects any remaining recursion.
//
// A computation c is self recursive when one of its top-level union arms
// joins back into c (eg. `read = owner + reader + parent->read`).
// Its least fixpoint is the non recursive arms evaluated on the object
// and on every object transitively reached through the recursive fields.
//
// Recursion anywhere else (under a difference or intersection, or through other computations)
// has no closure equivalent, so it is reported as an error.
func (l *lowering) resolveRecursion(model *Model) {
	for _, r := range model.Resources {
		for _, c := range values(r) {
			c.Expr = closeSelfRecursion(c)
		}
	}

	state := make(map[*Computation]int) // 0: unvisited, 1: visiting, 2: done
	var visit func(c *Computation, path []*Computation)
	visit = func(c *Computation, path []*Computation) {
		switch state[c] {
		case 1:
			l.errorf("recursive definition not supported: %v: recursion is only supported through a tuple to userset or userset in a top-level union",
				formatCycle(c, path))
			return
		case 2:
			return
		}
		state[c] = 1
		for _, dep := range deps(c.Expr) {
			visit(dep, append(path, c))
		}
		state[c] = 2
	}
	for _, r := range model.Resources {
		for _, c := range values(r) {
			visit(c, nil)
		}
	}
}

func closeSelfRecursion(c *Computation) Expr {
	var fields []*Field
	var local []Expr
	for _, arm := range unionArms(c.Expr) {
		if join, ok := arm.(*Join); ok && isRef(join.Expr, c) {
			fields = append(fields, join.Field)
			continue
		}
		local = append(local, arm)
	}
	if len(fields) == 0 {
		return c.Expr
	}
	return &Closure{Fields: fields, Local: union(local...)}
}

// deps returns the computations e depends on
func deps(e Expr) []*Computation {
	switch e := e.(type) {
	case *Ref:
		return []*Computation{e.To}
	case *Join:
		return deps(e.Expr)
	case *Closure:
		return deps(e.Local)
	case *Binary:
		return append(deps(e.Left), deps(e.Right)...)
	}
	return nil
}

func isRef(e Expr, c *Computation) bool {
	ref, ok := e.(*Ref)
	return ok && ref.To == c
}

func values(r *Resource) []*Computation {
	return append(append([]*Computation{}, r.Relations...), r.Permissions...)
}

func formatCycle(c *Computation, path []*Computation) string {
	var names []string
	for i, p := range path {
		if p == c {
			for _, q := range path[i:] {
				names = append(names, q.Resource+"."+q.Name)
			}
			break
		}
	}
	names = append(names, c.Resource+"."+c.Name)
	return strings.Join(names, " -> ")
}
