package graph

import "expert_systems/pkg/models/edge"

type Graph []edge.Edge

func NewExampleGraph() Graph {
	// stack
	gr := make(Graph, 0)
	return gr
}
