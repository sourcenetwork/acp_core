package policy_compiler

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/sourcenetwork/acp_core/pkg/errors"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

var identifier = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)

// UnmarshalRequirements parses a yaml contracts document:
//
//	requirements:
//	- name: read-access
//	  statement: actors can read if they're readers or writers, but not banned
//	  notes:
//	  - Assumed writers means the doc's writers, not its folder's.
//	  contracts:
//	  - name: unbanned-readers-can-read
//	    description: Readers can read the doc, unless they're banned.
//	    status: accepted                # or proposed, the default
//	    properties:
//	    - type: assert
//	      description: every unbanned reader of a doc can read it
//	      formula: |
//	        all d: Doc | d.doc_reader - d.doc_banned in d.doc_read
//	    - type: pred
//	      description: unbanned readers can exist
//	      formula: |
//	        some d: Doc | some d.doc_reader - d.doc_banned
//
// Formulas should be block scalars (|): alloy's `x: T` declarations aren't valid in plain yaml strings.
func UnmarshalRequirements(definition string) ([]*types.Requirement, error) {
	var doc struct {
		Requirements []struct {
			Name      string   `json:"name"`
			Statement string   `json:"statement"`
			Notes     []string `json:"notes"`
			Contracts []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Status      string `json:"status"`
				Properties  []struct {
					Type        string `json:"type"`
					Description string `json:"description"`
					Formula     string `json:"formula"`
				} `json:"properties"`
			} `json:"contracts"`
		} `json:"requirements"`
	}
	if err := yaml.UnmarshalStrict([]byte(definition), &doc); err != nil {
		return nil, errors.NewWithCause("unmarshaling contracts", err, errors.ErrorType_BAD_INPUT)
	}

	errs := errors.NewMultiError(ErrPolicyCompiler)
	var requirements []*types.Requirement
	for _, r := range doc.Requirements {
		requirement := &types.Requirement{Name: r.Name, Statement: strings.TrimSpace(r.Statement), Notes: r.Notes}
		for _, c := range r.Contracts {
			contract := &types.Contract{Name: c.Name, Description: strings.TrimSpace(c.Description)}
			switch c.Status {
			case "", "proposed":
				contract.Status = types.ContractStatus_PROPOSED
			case "accepted":
				contract.Status = types.ContractStatus_ACCEPTED
			default:
				errs.Append(fmt.Errorf("contract %v: status must be proposed or accepted, got %q", c.Name, c.Status))
			}
			for i, p := range c.Properties {
				property := &types.Property{Description: strings.TrimSpace(p.Description), Formula: strings.TrimSpace(p.Formula)}
				switch p.Type {
				case "assert":
					property.Type = types.PropertyType_ASSERT
				case "pred":
					property.Type = types.PropertyType_PRED
				default:
					errs.Append(fmt.Errorf("contract %v: property %v: type must be assert or pred, got %q", c.Name, i+1, p.Type))
				}
				contract.Properties = append(contract.Properties, property)
			}
			requirement.Contracts = append(requirement.Contracts, contract)
		}
		requirements = append(requirements, requirement)
	}
	if len(errs.GetErrors()) > 0 {
		return nil, errs
	}
	if err := ValidateRequirements(requirements); err != nil {
		return nil, err
	}
	return requirements, nil
}

// ValidateRequirements checks that requirements and contracts have unique, well formed names,
// and that every statement, description and formula is present. Formulas themselves are only validated by Alloy.
func ValidateRequirements(requirements []*types.Requirement) error {
	errs := errors.NewMultiError(ErrPolicyCompiler)
	// names must stay distinct as alloy identifiers, where - becomes _
	requirementNames := make(map[string]string)
	contractNames := make(map[string]string)
	checkName := func(kind, name string, seen map[string]string) {
		if !identifier.MatchString(name) {
			errs.Append(fmt.Errorf("%v %q: invalid name: expected letters, digits, - or _, starting with a letter", kind, name))
			return
		}
		if other, ok := seen[ContractIdent(name)]; ok {
			errs.Append(fmt.Errorf("%v %q: name clashes with %q", kind, name, other))
		}
		seen[ContractIdent(name)] = name
	}

	for _, r := range requirements {
		checkName("requirement", r.Name, requirementNames)
		if strings.TrimSpace(r.Statement) == "" {
			errs.Append(fmt.Errorf("requirement %v: statement is required: the user's words", r.Name))
		}
		for _, c := range r.Contracts {
			checkName("contract", c.Name, contractNames)
			if strings.TrimSpace(c.Description) == "" {
				errs.Append(fmt.Errorf("contract %v: description is required", c.Name))
			}
			for i, p := range c.Properties {
				if strings.TrimSpace(p.Description) == "" {
					errs.Append(fmt.Errorf("contract %v: property %v: description is required", c.Name, i+1))
				}
				if strings.TrimSpace(p.Formula) == "" {
					errs.Append(fmt.Errorf("contract %v: property %v: formula is required (a contract without properties is a draft)", c.Name, i+1))
				}
			}
		}
	}
	if len(errs.GetErrors()) > 0 {
		return errs
	}
	return nil
}

// ContractIdent returns the alloy identifier for a name: doc-edit -> doc_edit
func ContractIdent(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
