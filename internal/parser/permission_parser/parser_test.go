package permission_parser

import (
	"testing"

	"github.com/sourcenetwork/acp_core/pkg/types"

	"github.com/stretchr/testify/require"
)

// Maybe use a LISP-like / Reverse Polish Notation tree representation to ease asserting
func TestPermissionParser(t *testing.T) {
	perm := "owner + (that + resource->rel) - bar"

	tree, report := Parse(perm)

	wantJson := `
	{
  "combNode": {
    "left": {
      "combNode": {
        "left": {
          "operation": {
            "cu": {
              "relation": "owner"
            }
          }
        },
        "combinator": "UNION",
        "right": {
          "combNode": {
            "left": {
              "operation": {
                "cu": {
                  "relation": "that"
                }
              }
            },
            "combinator": "UNION",
            "right": {
              "operation": {
                "ttu": {
                  "resource": "resource",
                  "relation": "rel"
                }
              }
            }
          }
        }
      }
    },
    "combinator": "DIFFERENCE",
    "right": {
      "operation": {
        "cu": {
          "relation": "bar"
        }
      }
    }
  }
}
	`
	want := &types.PermissionFetchTree{}
	err := want.UnmarshalJSON(wantJson)
	require.NoError(t, err)
	require.False(t, report.HasError())
	require.Equal(t, want, tree)
}
