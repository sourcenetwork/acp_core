package parser

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// parserCaller represents is a clouser which calls a method within a parser
type parserCaller func(parser *TheoremParser) antlr.ParseTree

// ParseRelationship attempts to extract a single Relationship from input
func ParseRelationship(input string) (*types.Relationship, *ParserReport) {
	obj, report := ParseRelationshipWithLocation(input)
	if report.HasError() {
		return nil, report
	}
	return obj.Obj, report
}

// ParseRelationship attempts to extract a single Relationship from input.
// Returns location information about parsed Relationship.
func ParseRelationshipWithLocation(input string) (*LocatedObject[*types.Relationship], *ParserReport) {
	result, msgs := parseAndVisit(input, func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_document()
	})
	report := NewParserReport("relationship parse report", msgs...)
	if report.HasError() {
		return nil, report
	}
	obj := result.(LocatedObject[*types.Relationship])
	return &obj, report
}

// ParseRelationships greedly parses relationships in input.
// Consumes all of the input stream
func ParseRelationships(input string) ([]*types.Relationship, *ParserReport) {
	rels, report := ParseRelationshipsWithLocation(input)
	if report.HasError() {
		return nil, report
	}
	return utils.MapSlice(rels, func(o LocatedObject[*types.Relationship]) *types.Relationship { return o.Obj }), report
}

// ParseRelationshipsWithLocation greedily parses the input for all relationships it can find
// and returns a located object for each.
func ParseRelationshipsWithLocation(relationshipSet string) ([]LocatedObject[*types.Relationship], *ParserReport) {
	result, msgs := parseAndVisit(relationshipSet, func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_set()
	})
	report := NewParserReport("relationship set parse report", msgs...)
	if report.HasError() {
		return nil, report
	}
	return result.([]LocatedObject[*types.Relationship]), report
}

// ParsePolicyTheorem greedily consumes the input and returns a PolicyTheorem
func ParsePolicyTheorem(input string) (*LocatedPolicyTheorem, *ParserReport) {
	result, msgs := parseAndVisit(input, func(p *TheoremParser) antlr.ParseTree {
		return p.Policy_thorem()
	})
	report := NewParserReport("policy theorem parse report", msgs...)
	if report.HasError() {
		return nil, report
	}

	return result.(*LocatedPolicyTheorem), report
}

// parserAndVisit handles the boilerplate to parse an input stream,
// parse one production rule as given by caller and visits the resulting tree
// using the custom visitor.
//
// Return visitor result or error
func parseAndVisit(input string, caller parserCaller) (any, []*types.LocatedMessage) {
	inputStream := antlr.NewInputStream(input)
	lexer := NewTheoremLexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := errListener{}

	parser := NewTheoremParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&errListener)

	tree := caller(parser)

	visitor := theoremVisitorImpl{}
	result := visitor.Visit(tree)
	if result == nil {
		return nil, errListener.GetMessages()
	}
	return result, errListener.GetMessages()
}
