package edge

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/vertex"
)

type Edge struct {
	Start *vertex.Vertex
	End   *vertex.Vertex
	Label enums.EdgeLabelEnum
}
