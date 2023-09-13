package edge

import "expert_systems/pkg/models/vertex"

type Edge struct {
	Start *vertex.Vertex
	End   *vertex.Vertex
	Label int
}
