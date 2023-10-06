package edge

import (
	"expert_systems/pkg/models/enums"
	"expert_systems/pkg/models/node"
)

type Edge struct {
	Start *node.Node
	End   *node.Node
	Label enums.EdgeLabelEnum
}
