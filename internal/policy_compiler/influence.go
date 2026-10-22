package policy_compiler

// addInfluencers derives, for each permission, the actors which may change its value.
//
// An actor influences a permission if it can manage any relation the permission depends on,
// either on the object itself or on the objects reached through usersets and tuple to usersets.
// Influence over relations is inlined into the permissions which use them.
func addInfluencers(model *Model) {
	inf := influence{
		influencers: make(map[*Computation]*Computation),
		authorities: make(map[string]map[string]*Computation),
	}
	for _, r := range model.Resources {
		inf.authorities[r.Name] = make(map[string]*Computation)
		for _, a := range r.Authorities {
			inf.authorities[r.Name][a.Name] = a
		}
		for _, c := range r.Permissions {
			i := &Computation{Kind: KindInfluence, Resource: r.Name, Name: c.Name}
			inf.influencers[c] = i
			r.Influencers = append(r.Influencers, i)
		}
	}
	for c, i := range inf.influencers {
		i.Expr = inf.expr(c.Expr)
	}
}

type influence struct {
	influencers map[*Computation]*Computation
	authorities map[string]map[string]*Computation
}

func (inf *influence) expr(e Expr) Expr {
	switch e := e.(type) {
	case *Tuples:
		return inf.authority(e.Field)
	case *Ref:
		if e.To.Kind == KindRelation {
			return inf.expr(e.To.Expr)
		}
		return &Ref{To: inf.influencers[e.To]}
	case *Join:
		return dedupe(union(inf.authority(e.Field), &Join{Field: e.Field, Expr: inf.expr(e.Expr)}))
	case *Closure:
		var exprs []Expr
		for _, f := range e.Fields {
			exprs = append(exprs, inf.authority(f))
		}
		exprs = append(exprs, inf.expr(e.Local))
		return &Closure{Fields: e.Fields, Local: dedupe(union(exprs...))}
	case *Binary:
		return dedupe(union(inf.expr(e.Left), inf.expr(e.Right)))
	}
	return &Empty{}
}

func (inf *influence) authority(f *Field) Expr {
	return &Ref{To: inf.authorities[f.Resource][f.Relation]}
}

// dedupe removes repeated references from a union
func dedupe(e Expr) Expr {
	seen := make(map[*Computation]bool)
	var arms []Expr
	for _, arm := range unionArms(e) {
		if ref, ok := arm.(*Ref); ok {
			if seen[ref.To] {
				continue
			}
			seen[ref.To] = true
		}
		arms = append(arms, arm)
	}
	return union(arms...)
}
