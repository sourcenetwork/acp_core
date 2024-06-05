package parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func ParseRelationship(relationship string) (*types.Relationship, *ParseErrors) {
	input := antlr.NewInputStream(relationship)
	lexer := NewTestSuiteLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := errListener{}

	parser := NewTestSuiteParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errListener)
	tree := parser.Relationship()

	listener := relationshipListener{}
	antlr.ParseTreeWalkerDefault.Walk(&listener, tree)

	err := errListener.GetError()
	if err != nil {
		return nil, err
	}

	rel := listener.GetRelationship()
	return rel, nil
}

func ParseRelationships(relationshipSet string) ([]*types.Relationship, error) {
	return nil, nil
}

func ParseTestSuite(relationshipSet string) ([]*types.Relationship, error) {
	return nil, nil
}
