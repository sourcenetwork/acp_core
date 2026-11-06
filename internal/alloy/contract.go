package alloy

import (
	"fmt"
	"strings"

	pc "github.com/sourcenetwork/acp_core/internal/policy_compiler"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

// contractIndent indents formulas in the emitted command
const contractIndent = "        "

// PropertyCommand returns the name of the alloy command which verifies the index-th property (from 0) of a contract
func PropertyCommand(contract *types.Contract, index int) string {
	return fmt.Sprintf("contract_%v_%v", pc.ContractIdent(contract.Name), index+1)
}

// EmitProperty renders an alloy module with a single command verifying the index-th property of a contract up to scope.
//
// An assert is checked in every state satisfying the assumptions:
// `check { assumptions implies { formula } }`, so an instance is a counterexample.
// A pred is run against them: `run { assumptions and { formula } }`, so an instance is an example.
func EmitProperty(model *pc.Model, module string, contract *types.Contract, index int, scope uint32) File {
	e := newEmitter(model, module)
	property := contract.Properties[index]
	command := PropertyCommand(contract, index)
	keyword, connective := "check", "implies"
	if property.Type == types.PropertyType_PRED {
		keyword, connective = "run", "and"
	}

	b := &strings.Builder{}
	fmt.Fprintf(b, "open %v\nopen %v_assumptions\n\n", e.module, e.module)
	for _, line := range strings.Split(contract.Description+"\n"+property.Description, "\n") {
		fmt.Fprintf(b, "// %v\n", line)
	}
	fmt.Fprintf(b, "%v %v {\n    assumptions %v {\n", keyword, command, connective)
	// the formula sits on its own lines, so a trailing comment can't swallow the closing braces
	for _, line := range strings.Split(property.Formula, "\n") {
		fmt.Fprintf(b, "%v%v\n", contractIndent, line)
	}
	fmt.Fprintf(b, "    }\n} for %v\n", scope)
	return File{Name: command + ".als", Content: b.String()}
}

// FormulaPosition maps a line and column of a module written by EmitProperty to the property's formula.
// ok is false if the position is outside the formula.
func FormulaPosition(file File, line, column int) (formulaLine, formulaColumn int, ok bool) {
	lines := strings.Split(file.Content, "\n")
	start := 0
	for i, l := range lines {
		if strings.HasPrefix(l, "    assumptions ") {
			start = i + 2 // 1-based line of the formula's first line
			break
		}
	}
	end := len(lines) - 2 // the formula is followed by the closing lines
	if start == 0 || line < start || line >= end || column <= len(contractIndent) {
		return 0, 0, false
	}
	return line - start + 1, column - len(contractIndent), true
}
