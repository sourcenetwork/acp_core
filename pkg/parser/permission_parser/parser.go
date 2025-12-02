package permission_parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/pkg/parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// Parser transforms a permission expressiong into a PermissionFetchTree and a parser report
//
// Invalid expressiosn fail to parse and the report will contain errros
func Parse(expr string) (*types.PermissionFetchTree, error) {
	tree, report := ParseWithReport(expr)
	if report.HasError() {
		return nil, report.ToMultiError(ErrPermissionParser)
	}
	return tree, nil
}

// Parser transforms a permission expressiong into a PermissionFetchTree and a parser report
//
// Invalid expressiosn fail to parse and the report will contain errros
func ParseWithReport(expr string) (*types.PermissionFetchTree, *parser.ParserReport) {
	endPosition := types.BufferPosition{
		Line:   1,
		Column: uint64(len(expr)),
	}
	lines := strings.Split(expr, "\n")
	if len(lines) > 0 {
		lineCount := uint64(len(lines))
		endPosition.Line = lineCount
		lastLine := lines[lineCount-1]
		endPosition.Column = uint64(len(lastLine))
	}

	errListener := parser.NewErrLsitener("Permission Expression", endPosition)

	inputStream := antlr.NewInputStream(expr)
	lexer := NewPermissionExprLexer(inputStream)
	lexer.AddErrorListener(errListener)

	stream := antlr.NewCommonTokenStream(lexer, 0)

	parser := NewPermissionExprParser(stream)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errListener)

	tree := parser.Expr()

	report := errListener.GetReport()
	if report.HasError() {
		return nil, report
	}

	v := visitor{}
	result := v.Visit(tree)
	if result == nil {
		result = &types.PermissionFetchTree{}
	}
	return result.(*types.PermissionFetchTree), report
}
