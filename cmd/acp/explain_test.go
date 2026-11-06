package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sourcenetwork/acp_core/pkg/types"
)

func Test_PrintGraph_FoldsUnions(t *testing.T) {
	authorized := &types.SearchResult{Authorized: true, Explored: true, Exhausted: true}
	explored := &types.SearchResult{Explored: true, Exhausted: true}
	node := func(id string, kind types.CheckExplainNodeType, text, detail string, result *types.SearchResult) *types.CheckExplainNode {
		return &types.CheckExplainNode{Id: id, NodeType: kind, Text: text, Detail: detail, Result: result}
	}
	union, userset, diff := types.CheckExplainNodeType_UNION_NODE, types.CheckExplainNodeType_USERSET_NODE, types.CheckExplainNodeType_DIFF_NODE
	graph := &types.CheckExplainGraph{
		RootNodeId: "1", // wrong on purpose, as zanzi's often is
		Nodes: []*types.CheckExplainNode{
			node("0", userset, "(doc:roadmap, view)", "reason: Request", authorized),
			node("1", union, "+", "", authorized),
			node("2", userset, "(doc:roadmap, owner)", "reason: Rule Computed Userset: owner", explored),
			node("3", userset, "(actor:did:user:carol)", "reason: Rule This", explored),
			node("4", diff, "-", "", authorized),
			node("5", userset, "(folder:specs, view)", "reason: Rule Tuple To Userset: folder view", authorized),
			node("6", userset, "(actor:did:user:bob)", "reason: Rule This", authorized),
			node("7", userset, "(doc:roadmap, banned)", "reason: Rule Computed Userset: banned", explored),
			node("8", userset, "(doc:roadmap, viewer)", "reason: Rule Computed Userset: viewer", nil),
		},
		Edges: []*types.CheckExplainEdge{
			{SourceNodeId: "0", DestNodeId: "1"},
			{SourceNodeId: "1", DestNodeId: "2"},
			{SourceNodeId: "2", DestNodeId: "3"},
			{SourceNodeId: "1", DestNodeId: "4"},
			{SourceNodeId: "4", DestNodeId: "5"},
			{SourceNodeId: "5", DestNodeId: "6"},
			{SourceNodeId: "4", DestNodeId: "7"},
			{SourceNodeId: "1", DestNodeId: "8"},
		},
	}
	out := &bytes.Buffer{}
	printGraph(out, graph)
	require.Equal(t, `✓ (doc:roadmap, view)
    (doc:roadmap, owner)  [computed userset: owner]
      (actor:did:user:carol)
  ✓ - (the first term, minus the others)
    ✓ (folder:specs, view)  [tuple to userset: folder view]
      ✓ (actor:did:user:bob)
      (doc:roadmap, banned)  [computed userset: banned]
    (doc:roadmap, viewer)  [computed userset: viewer]  (not explored)
`, out.String())
}
