import enum
import typing
import dataclasses

class NodeState(enum.Enum):
    BLACK = 1
    GRAY  = 2
    WHITE = 3

@dataclasses.dataclass(unsafe_hash=True)
class Node:
    """Вершина"""
    number: int
    state: NodeState = dataclasses.field(default=NodeState.WHITE, compare=False)

class EdgeState(enum.Enum):
    TRAVERSED = 1
    FREE  = 2
    FORBIDDEN = 3

@dataclasses.dataclass
class Edge:
    """Ребро"""
    start: Node
    end: Node
    label: EdgeState = EdgeState.FREE

class Graph:
    """Граф"""
    def __init__(self, *edges: Edge):
        self._edges = list(edges)
    
    def edges(self) -> typing.List[Edge]:
        """Список рёбер"""
        return self._edges

    def nodees(self) -> typing.List[Node]:
        """Список вершин"""
        set_node = set()
        for edge in self._edges:
            set_node.add(edge.start)
            set_node.add(edge.end)
        return list(set_node)

