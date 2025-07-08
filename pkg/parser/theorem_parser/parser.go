package theorem_parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/sourcenetwork/acp_core/pkg/parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"github.com/sourcenetwork/acp_core/pkg/utils"
)

// parserCaller represents is a clouser which calls a method within a parser
type parserCaller func(parser *TheoremParser) antlr.ParseTree

// ParseRelationship attempts to extract a single Relationship from input
func ParseRelationship(input string) (*types.Relationship, *parser.ParserReport) {
	obj, report := ParseRelationshipWithLocation(input)
	if report.HasError() {
		return nil, report
	}
	return obj.Obj, report
}

// ParseRelationship attempts to extract a single Relationship from input.
// Returns location information about parsed Relationship.
func ParseRelationshipWithLocation(input string) (parser.LocatedObject[*types.Relationship], *parser.ParserReport) {
	result, report := parseAndVisit(input, "relationship document", func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_document()
	})
	if report.HasError() {
		return parser.LocatedObject[*types.Relationship]{}, report
	}
	return result.(parser.LocatedObject[*types.Relationship]), report
}

// ParseRelationships greedly parses relationships in input.
// Consumes all of the input stream
func ParseRelationships(input string) ([]*types.Relationship, *parser.ParserReport) {
	rels, report := ParseRelationshipsWithLocation(input)
	if report.HasError() {
		return nil, report
	}
	return utils.MapSlice(rels, func(o parser.LocatedObject[*types.Relationship]) *types.Relationship { return o.Obj }), report
}

// ParseRelationshipsWithLocation greedily parses the input for all relationships it can find
// and returns a located object for each.
func ParseRelationshipsWithLocation(relationshipSet string) ([]parser.LocatedObject[*types.Relationship], *parser.ParserReport) {
	result, report := parseAndVisit(relationshipSet, "relationship set", func(p *TheoremParser) antlr.ParseTree {
		return p.Relationship_set()
	})
	if report.HasError() {
		return nil, report
	}
	return result.([]parser.LocatedObject[*types.Relationship]), report
}

// ParsePolicyTheorem greedily consumes the input and returns a PolicyTheorem
func ParsePolicyTheorem(input string) (*parser.LocatedPolicyTheorem, *parser.ParserReport) {
	result, report := parseAndVisit(input, "policy theorem", func(p *TheoremParser) antlr.ParseTree {
		return p.Policy_thorem()
	})
	if report.HasError() {
		return nil, report
	}

	return result.(*parser.LocatedPolicyTheorem), report
}

// parserAndVisit handles the boilerplate to parse an input stream,
// parse one production rule as given by caller and visits the resulting tree
// using the custom visitor.
//
// Return visitor result or error
func parseAndVisit(input string, productionName string, caller parserCaller) (any, *parser.ParserReport) {
	endPosition := types.BufferPosition{
		Line:   1,
		Column: 1,
	}
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		lineCount := uint64(len(lines))
		endPosition.Line = lineCount
		lastLine := lines[lineCount-1]
		endPosition.Column = uint64(len(lastLine))
	}

	inputStream := antlr.NewInputStream(input)
	lexer := NewTheoremLexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := parser.NewErrLsitener(productionName, endPosition)

	parser := NewTheoremParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errListener)
	tree := caller(parser)

	report := errListener.GetReport()
	if report.HasError() {
		return nil, report
	}

	visitor := theoremVisitorImpl{}
	result := visitor.Visit(tree)
	return result, report
}
