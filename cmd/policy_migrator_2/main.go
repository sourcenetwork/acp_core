package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strconv"
	"strings"

	"os"

	"github.com/sourcenetwork/acp_core/internal/policy"
	"github.com/sourcenetwork/acp_core/pkg/parser/permission_parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
	"sigs.k8s.io/yaml"
)

func main() {
	files := os.Args[1:len(os.Args)]

	fs := token.NewFileSet()
	for _, file := range files {

		parsed, err := parser.ParseFile(fs, file, nil, parser.SkipObjectResolution|parser.ParseComments)
		if err != nil {
			log.Fatalf("parsing file: %v: %v", file, err)
		}
		visitor := stringLiteralVisitor{
			File: file,
		}
		ast.Walk(&visitor, parsed)

		outFile, err := os.Create(file)
		if err != nil {
			log.Fatalf("writing file: %v: %v", file, err)
		}
		defer outFile.Close()

		printer.Fprint(outFile, fs, parsed)
	}
}

type stringLiteralVisitor struct {
	File string
}

func (v *stringLiteralVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BasicLit:
		if n.Kind != token.STRING {
			return v
		}
		var err error
		literal := ""
		doubleQuote := false
		if n.Value[0] == '"' {
			literal, err = strconv.Unquote(n.Value)
			if err != nil {
				panic(err)
			}
			doubleQuote = true
		} else { //multiline string
			literal, err = strconv.Unquote(n.Value)
			if err != nil {
				panic(err)
			}
		}
		// if string is just whitespace, do nothing
		if strings.TrimSpace(literal) == "" {
			return v
		}

		if literal == "null" {
			return v
		}

		pol, err := policy.Unmarshal(literal, types.PolicyMarshalingType_YAML)
		if err != nil {
			return v
		}

		// skip if policy had no resources - could be due to a json match or something similar
		if len(pol.Resources) == 0 {
			return v
		}

		for _, resource := range pol.Resources {
			resource.Relations = utils.FilterSlice(resource.Relations, func(r *types.Relation) bool { return r.Name != "owner" })
			for _, p := range resource.Permissions {
				tree, err := permission_parser.Parse(p.Expression)
				if err != nil {
					return v
				}
				tree = removeOwnerFromTree(tree)
				if tree == nil {
					p.Expression = ""
				} else {
					p.Expression = tree.IntoPermissionExpr()
				}

			}
		}

		yamlPol := mapPolicyToYaml(pol)
		bytes, err := yaml.Marshal(yamlPol)
		if err != nil {
			panic(err)
		}
		str := string(bytes)
		if doubleQuote {
			str = strings.TrimSpace(str)
			str = strconv.Quote(str)
		} else {
			str = "`\n" + str + "`"
		}

		n.Value = string(str)
		log.Printf("update file %v: len %v", v.File, len(literal))
	}
	return v
}

func mapPolicyToYaml(in *types.Policy) *types.PolicyYaml {
	out := &types.PolicyYaml{
		Name:        in.Name,
		Description: in.Description,
		Meta:        in.Attributes,
		Spec:        mapSpecTypeToString(in.SpecificationType),
		Resources:   utils.MapSlice(in.Resources, mapResourceToYaml),
		Actor:       mapActorResourceToYaml(in.ActorResource),
	}

	return out
}

func mapResourceToYaml(r *types.Resource) *types.ResourceYaml {
	return &types.ResourceYaml{
		Name:        r.Name,
		Description: r.Doc,
		Permissions: utils.MapSlice(r.Permissions, mapPermissionToYaml),
		Relations:   utils.MapSlice(r.Relations, mapRelationToYaml),
	}
}

func mapRelationToYaml(rel *types.Relation) *types.RelationYaml {
	return &types.RelationYaml{
		Name:    rel.Name,
		Doc:     rel.Doc,
		Manages: rel.Manages,
		Types:   utils.MapSlice(rel.VrTypes, mapRestrictionToString),
	}
}

func mapPermissionToYaml(p *types.Permission) *types.PermissionYaml {
	return &types.PermissionYaml{
		Name: p.Name,
		Doc:  p.Doc,
		Expr: p.Expression,
	}
}

func mapActorResourceToYaml(in *types.ActorResource) *types.ActorResourceYaml {
	if in == nil {
		return nil
	}

	return &types.ActorResourceYaml{
		Name:      in.Name,
		Doc:       in.Doc,
		Relations: utils.MapSlice(in.Relations, mapRelationToYaml),
	}
}

func mapRestrictionToString(in *types.Restriction) string {
	if in.RelationName != "" {
		return in.ResourceName + "->" + in.RelationName
	}
	return in.ResourceName
}

func mapSpecTypeToString(specType types.PolicySpecificationType) string {
	switch specType {
	case types.PolicySpecificationType_DEFRA_SPEC:
		return "defra"
	case types.PolicySpecificationType_NO_SPEC:
		return ""
	default:
		return "unknown"
	}
}

// remove owner from tree if left or right nodes are CU
func removeOwnerFromTree(tree *types.PermissionFetchTree) *types.PermissionFetchTree {
	switch term := tree.Term.(type) {
	case *types.PermissionFetchTree_CombNode:
		if term.CombNode.Combinator == types.Combinator_UNION {
			if isOwnerCu(term.CombNode.Left) {
				return removeOwnerFromTree(term.CombNode.Right)
			} else if isOwnerCu(term.CombNode.Right) {
				return removeOwnerFromTree(term.CombNode.Left)
			}
		}
		term.CombNode.Left = removeOwnerFromTree(term.CombNode.Left)
		term.CombNode.Right = removeOwnerFromTree(term.CombNode.Right)
		return tree
	case *types.PermissionFetchTree_Operation:
		if isOwnerCu(tree) {
			return nil
		}
		return tree
	default:
		return tree
	}
}

func isOwnerCu(t *types.PermissionFetchTree) bool {
	op := t.GetOperation()
	if op == nil {
		return false
	}
	cu := op.GetCu()
	if cu == nil {
		return false
	}
	return cu.Relation == "owner"
}
