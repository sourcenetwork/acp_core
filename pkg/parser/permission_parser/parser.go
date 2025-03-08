package permission_parser

import (
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
	inputStream := antlr.NewInputStream(expr)
	lexer := NewPermissionExprLexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, 0)

	errListener := parser.NewErrLsitener("Permission Expression")

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
	return result.(*types.PermissionFetchTree), report
}
