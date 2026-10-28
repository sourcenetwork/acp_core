package alloy

import (
	"fmt"
	"strings"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
)

// EmitGenerate renders an alloy module with a single run command, name, which generates worlds:
// states satisfying the assumptions, with up to scope objects of each resource and exactly actors actors,
// every one of which takes part in some relationship.
func EmitGenerate(model *pc.Model, module, name string, scope, actors uint32) File {
	e := newEmitter(model, module)

	// actors related to any object, through any relation
	var related []string
	for _, r := range model.Resources {
		for _, f := range r.Fields {
			if f.Target == pc.ActorTarget {
				related = append(related, "univ."+e.qualified(f))
			}
		}
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "open %v\nopen %v_assumptions\n\n", e.module, e.module)
	fmt.Fprintf(b, "run %v {\n    assumptions\n    Actor in %v\n} for %v but exactly %v Actor\n",
		name, strings.Join(related, " + "), scope, actors)
	return File{Name: name + ".als", Content: b.String()}
}
