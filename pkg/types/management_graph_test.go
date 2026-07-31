package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManagementGraphIsWellFormedReturnsDeterministicError(t *testing.T) {
	graph := &ManagementGraph{
		ForwardEdges: map[string]*ManagerEdges{
			"zulu/source": {
				Edges: map[string]bool{
					"zulu/destination":  true,
					"alpha/destination": true,
				},
			},
			"alpha/source": {
				Edges: map[string]bool{
					"zulu/destination":  true,
					"alpha/destination": true,
				},
			},
		},
	}

	for range 100 {
		require.EqualError(t, graph.IsWellFormed(),
			"edge defined from alpha/source to alpha/destination: alpha/source not found")
	}
}
