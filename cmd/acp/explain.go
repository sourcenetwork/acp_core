package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

// printGraph renders the engine's search graph as an indented tree, marking nodes on an authorized path with ✓
// and nodes the search skipped with (not explored).
// Unions are folded into their parent, as they only nest the terms of a `+`. Differences and intersections are kept,
// since they change the meaning of their children. The tree starts from the nodes without incoming edges:
// zanzi's root node id is unreliable.
func printGraph(w io.Writer, graph *types.CheckExplainGraph) {
	if graph == nil {
		return
	}
	nodes := make(map[string]*types.CheckExplainNode)
	for _, n := range graph.Nodes {
		nodes[n.Id] = n
	}
	children := make(map[string][]string)
	hasParent := make(map[string]bool)
	for _, e := range graph.Edges {
		children[e.SourceNodeId] = append(children[e.SourceNodeId], e.DestNodeId)
		hasParent[e.DestNodeId] = true
	}

	var visit func(id string, depth int)
	visit = func(id string, depth int) {
		n := nodes[id]
		if n == nil {
			return
		}
		if n.NodeType != types.CheckExplainNodeType_UNION_NODE {
			mark := "  "
			if n.Result != nil && n.Result.Authorized {
				mark = "✓ "
			}
			skipped := ""
			if n.Result == nil || !n.Result.Explored {
				skipped = "  (not explored)" // the search had already found an answer
			}
			fmt.Fprintf(w, "%v%v%v%v%v\n", strings.Repeat("  ", depth), mark, nodeLabel(n), nodeReason(n), skipped)
			depth++
		}
		for _, child := range children[id] {
			visit(child, depth)
		}
	}
	for _, n := range graph.Nodes {
		if !hasParent[n.Id] {
			visit(n.Id, 0)
		}
	}
}

func nodeLabel(n *types.CheckExplainNode) string {
	switch n.NodeType {
	case types.CheckExplainNodeType_DIFF_NODE:
		return "- (the first term, minus the others)"
	case types.CheckExplainNodeType_INTERSECTION_NODE:
		return "& (in every term)"
	}
	return n.Text
}

// nodeReason renders a node's detail, eg. `reason: Rule Tuple To Userset: folder view` as `  [tuple to userset: folder view]`
func nodeReason(n *types.CheckExplainNode) string {
	reason := strings.TrimPrefix(strings.TrimPrefix(n.Detail, "reason: "), "Rule ")
	switch reason {
	case "", "This", "Request":
		return ""
	}
	return "  [" + strings.ToLower(reason) + "]"
}
