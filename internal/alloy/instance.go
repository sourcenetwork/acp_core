package alloy

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
)

// Instance is an Alloy instance translated into acp_core relationships
type Instance struct {
	// Command is the Alloy command which produced the instance
	Command string
	// Relationships are in the PolicyTheoremDSL syntax, owner relationships first
	Relationships []string
	// Skolems names the objects and actors an instance is about (eg. a counterexample's witness),
	// mapping the skolem name to acp ids
	Skolems map[string][]string
	// Objects and Actors list every acp id in the instance
	Objects []string
	Actors  []string
}

// actorNames are given to actor atoms in order, falling back to actorN
var actorNames = []string{"alice", "bob", "carol", "dave", "erin", "frank", "grace", "heidi", "ivan", "judy"}

// xmlInstance mirrors the instance format written by `alloy exec -t xml`
type xmlInstance struct {
	Command string     `xml:"command,attr"`
	Sigs    []xmlSig   `xml:"sig"`
	Fields  []xmlField `xml:"field"`
	Skolems []xmlField `xml:"skolem"`
}

type xmlSig struct {
	Label string    `xml:"label,attr"`
	ID    string    `xml:"ID,attr"`
	Atoms []xmlAtom `xml:"atom"`
}

type xmlField struct {
	Label    string     `xml:"label,attr"`
	ParentID string     `xml:"parentID,attr"`
	Tuples   []xmlTuple `xml:"tuple"`
}

type xmlTuple struct {
	Atoms []xmlAtom `xml:"atom"`
}

type xmlAtom struct {
	Label string `xml:"label,attr"`
}

// ParseInstance translates the first instance of an Alloy XML solution
// into relationships of the policy model was compiled from.
//
// Objects are named after their resource and atom index (file0),
// actors get names in atom order (did:user:alice).
func ParseInstance(model *pc.Model, data []byte) (*Instance, error) {
	var doc struct {
		Instances []xmlInstance `xml:"instance"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing alloy instance: %w", err)
	}
	if len(doc.Instances) == 0 {
		return nil, fmt.Errorf("parsing alloy instance: no instance found")
	}
	in := doc.Instances[0]

	resources := make(map[string]*pc.Resource)
	for _, r := range model.Resources {
		resources[sigName(r.Name)] = r
	}

	inst := &Instance{Command: in.Command, Skolems: make(map[string][]string)}
	ids := make(map[string]string) // atom label -> acp id
	sigResources := make(map[string]*pc.Resource)
	for _, sig := range in.Sigs {
		name := unqualify(sig.Label)
		if name == "Actor" {
			for _, atom := range sig.Atoms {
				id := "did:user:" + actorName(atomIndex(atom.Label))
				ids[atom.Label] = id
				inst.Actors = append(inst.Actors, id)
			}
			continue
		}
		r, ok := resources[name]
		if !ok {
			continue
		}
		sigResources[sig.ID] = r
		for _, atom := range sig.Atoms {
			id := fmt.Sprintf("%v:%v%v", r.Name, r.Name, atomIndex(atom.Label))
			ids[atom.Label] = id
			inst.Objects = append(inst.Objects, id)
		}
	}

	var owners, others []string
	for _, field := range in.Fields {
		r, ok := sigResources[field.ParentID]
		if !ok {
			continue
		}
		f := lookupField(r, field.Label)
		if f == nil {
			return nil, fmt.Errorf("resource %v has no field %v: was the instance generated from this policy?", r.Name, field.Label)
		}
		for _, tuple := range field.Tuples {
			if len(tuple.Atoms) != 2 {
				return nil, fmt.Errorf("field %v: expected binary tuples", field.Label)
			}
			subject := ids[tuple.Atoms[1].Label]
			if f.Userset != nil {
				subject += "#" + f.Userset.Name
			}
			line := fmt.Sprintf("%v#%v@%v", ids[tuple.Atoms[0].Label], f.Relation, subject)
			if f.One {
				owners = append(owners, line)
			} else {
				others = append(others, line)
			}
		}
	}
	sort.Strings(owners)
	sort.Strings(others)
	inst.Relationships = append(owners, others...)

	for _, skolem := range in.Skolems {
		var values []string
		for _, tuple := range skolem.Tuples {
			var atoms []string
			for _, atom := range tuple.Atoms {
				if id, ok := ids[atom.Label]; ok {
					atoms = append(atoms, id)
				}
			}
			values = append(values, strings.Join(atoms, "->"))
		}
		inst.Skolems[skolem.Label] = values
	}
	return inst, nil
}

func lookupField(r *pc.Resource, name string) *pc.Field {
	for _, f := range r.Fields {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// unqualify strips the module from a sig label: policy/File -> File
func unqualify(label string) string {
	return label[strings.LastIndex(label, "/")+1:]
}

// atomIndex returns the index of an atom label: policy/File$3 -> 3
func atomIndex(label string) string {
	return label[strings.LastIndex(label, "$")+1:]
}

func actorName(index string) string {
	var i int
	if _, err := fmt.Sscan(index, &i); err == nil && i < len(actorNames) {
		return actorNames[i]
	}
	return "actor" + index
}
