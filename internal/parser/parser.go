package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

func ParseRelationship(relationship string) (*types.Relationship, error) {
	relationships, err := ParseRelationships(relationship)
	if err != nil {
		return nil, err
	}
	if len(relationships) == 0 {
		return nil, nil //probably sould return an error instead
	}
	return relationships[0], nil
}

func ParseRelationships(relationshipSet string) ([]*types.Relationship, error) {
	rels, err := ParseRelationshipsWithPosition(relationshipSet)
	if err != nil {
		return nil, err
	}

	return utils.MapSlice(rels, func(o IndexedObject[*types.Relationship]) *types.Relationship { return o.Obj }), nil
}

func ParseRelationshipsWithPosition(relationshipSet string) ([]IndexedObject[*types.Relationship], error) {
	input := antlr.NewInputStream(relationshipSet)
	lexer := NewTheoremLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := errListener{}

	parser := NewTheoremParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errListener)
	tree := parser.Relationship_set()

	err := errListener.GetError()
	if err != nil {
		return nil, err
	}

	visitor := newTheoremVisitor()
	result := visitor.Visit(tree)
	if result == nil {
		return nil, nil
	}
	return result.([]IndexedObject[*types.Relationship]), nil
}

func ParsePolicyTheorem(policyTheorem string) (*IndexedPolicyTheorem, error) {
	input := antlr.NewInputStream(policyTheorem)
	lexer := NewTheoremLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := errListener{}

	parser := NewTheoremParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errListener)
	tree := parser.Policy_thorem()

	err := errListener.GetError()
	if err != nil {
		return nil, err
	}

	visitor := newTheoremVisitor()
	result := visitor.Visit(tree)
	if result == nil {
		return nil, nil
	}
	return result.(*IndexedPolicyTheorem), nil
}
