package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/types"
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
	return result.([]*types.Relationship), nil
}

func ParsePolicyTheorem(policyTheorem string) (*types.PolicyTheorem, error) {
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
	return result.(*types.PolicyTheorem), nil
}
