package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// parserCaller represents is a clouser which calls a method within a parser
type parserCaller func(parser *TheoremParser) antlr.ParseTree

// ParseRelationship attempts to extract a single Relationship from input
func ParseRelationship(input string) (*types.Relationship, error) {
	obj, err := ParseRelationshipWithLocation(input)
	if err != nil {
		return nil, errors.Wrap("parse relationship failed", err)
	}
	return obj.Obj, nil
}

// ParseRelationship attempts to extract a single Relationship from input.
// Returns location information about parsed Relationship.
func ParseRelationshipWithLocation(input string) (*LocatedObject[*types.Relationship], error) {
	result, err := parseAndVisit(input, func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_document()
	})
	if err != nil {
		return nil, errors.Wrap("parse relationship failed", err)
	}
	obj := result.(LocatedObject[*types.Relationship])
	return &obj, nil
}

// ParseRelationships greedly parses relationships in input.
// Consumes all of the input stream
func ParseRelationships(input string) ([]*types.Relationship, error) {
	rels, err := ParseRelationshipsWithLocation(input)
	if err != nil {
		return nil, errors.Wrap("parse relationships failed", err)
	}
	return utils.MapSlice(rels, func(o LocatedObject[*types.Relationship]) *types.Relationship { return o.Obj }), nil
}

// ParseRelationshipsWithLocation greedily parses the input for all relationships it can find
// and returns a located object for each.
func ParseRelationshipsWithLocation(relationshipSet string) ([]LocatedObject[*types.Relationship], error) {
	result, err := parseAndVisit(relationshipSet, func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_set()
	})
	if err != nil {
		return nil, errors.Wrap("parse relationships failed", err)
	}
	return result.([]LocatedObject[*types.Relationship]), nil
}

// ParsePolicyTheorem greedily consumes the input and returns a PolicyTheorem
func ParsePolicyTheorem(input string) (*LocatedPolicyTheorem, error) {
	result, err := parseAndVisit(input, func(p *TheoremParser) antlr.ParseTree {
		return p.Policy_thorem()
	})
	if err != nil {
		return nil, errors.Wrap("parse policy theorem failed", err)
	}
	return result.(*LocatedPolicyTheorem), nil
}

// parserAndVisit handles the boilerplate to parse an input stream,
// parse one production rule as given by caller and visits the resulting tree
// using the custom visitor.
//
// Return visitor result or error
func parseAndVisit(input string, caller parserCaller) (any, error) {
	inputStream := antlr.NewInputStream(input)
	lexer := NewTheoremLexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := errListener{}

	parser := NewTheoremParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errListener)

	tree := caller(parser)
	err := errListener.GetError()
	if err != nil {
		return nil, err
	}

	visitor := theoremVisitorImpl{}
	result := visitor.Visit(tree)
	return result, nil
}
