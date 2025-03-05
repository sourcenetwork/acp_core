package permission_parser

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/sourcenetwork/acp_core/internal/parser"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Parse(expr string) (*types.PermissionFetchTree, *parser.ParserReport) {
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
