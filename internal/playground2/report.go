package playground2

import (
	"fmt"
	"io"
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

// PrintReport writes a plain text report of a Verify response for the requirements it verified,
// calling save (if not nil) with each world and the name to save it as.
// It returns the number of contracts which aren't satisfied, drafts excluded.
func PrintReport(w io.Writer, requirements []*types.Requirement, resp *types.VerifyResponse, save func(name string, world *types.World)) int {
	contracts := make(map[string]*types.Contract)
	notes := make(map[string][]string)
	for _, r := range requirements {
		notes[r.Name] = r.Notes
		for _, c := range r.Contracts {
			contracts[c.Name] = c
		}
	}

	failed, drafts := 0, 0
	requirement := ""
	for _, result := range resp.Results {
		if result.Requirement != requirement {
			requirement = result.Requirement
			fmt.Fprintf(w, "\n== %v\n", requirement)
			for _, note := range notes[requirement] {
				fmt.Fprintf(w, "note: %v\n", note)
			}
		}
		contract := contracts[result.Contract]
		accepted := ""
		if contract.Status == types.ContractStatus_ACCEPTED {
			accepted = ", accepted"
		}
		fmt.Fprintf(w, "%-10v %v%v: %v\n", result.Status, result.Contract, accepted, contract.Description)
		switch result.Status {
		case types.VerifyStatus_DRAFT:
			drafts++
			continue
		case types.VerifyStatus_SATISFIED:
		default:
			failed++
		}

		for i, p := range result.Properties {
			property := contract.Properties[i]
			kind := strings.ToLower(property.Type.String())
			fmt.Fprintf(w, "    %v. %-10v %v (%v)", i+1, p.Status, property.Description, kind)
			switch {
			case p.Status == types.VerifyStatus_UNVERIFIED:
				fmt.Fprintf(w, ": %v\n", p.Message)
			case p.Status == types.VerifyStatus_SATISFIED && property.Type == types.PropertyType_ASSERT:
				fmt.Fprintf(w, ": holds in every world up to scope %v\n", resp.Scope)
			case p.Status == types.VerifyStatus_SATISFIED:
				fmt.Fprintf(w, ": example\n")
			case property.Type == types.PropertyType_ASSERT:
				fmt.Fprintf(w, ": counterexample\n")
			default:
				fmt.Fprintf(w, ": holds in no world up to scope %v\n", resp.Scope)
			}
			if p.World == nil {
				continue
			}
			for _, b := range p.Bindings {
				fmt.Fprintf(w, "         %v = %v\n", b.Name, strings.Join(b.Values, ", "))
			}
			for _, line := range strings.Split(strings.TrimSpace(p.World.Relationships), "\n") {
				fmt.Fprintf(w, "         %v\n", line)
			}
			if save != nil {
				save(fmt.Sprintf("%v-%v", result.Contract, i+1), p.World)
			}
		}
	}

	verified := len(resp.Results) - drafts
	fmt.Fprintf(w, "\n%v of %v contracts satisfied up to scope %v", verified-failed, verified, resp.Scope)
	if drafts > 0 {
		fmt.Fprintf(w, ", %v drafts", drafts)
	}
	fmt.Fprintln(w)
	return failed
}
