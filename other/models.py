import enum
import typing
import dataclasses

class VertexState(enum.Enum):
    BLACK = 1
    GRAY  = 2
    WHITE = 3

@dataclasses.dataclass(unsafe_hash=True)
class Vertex:
    """Вершина"""
    number: int
    state: VertexState = dataclasses.field(default=VertexState.WHITE, compare=False)

class EdgeState(enum.Enum):
    TRAVERSED = 1
    FREE  = 2
    FORBIDDEN = 3

@dataclasses.dataclass
class Edge:
    """Ребро"""
    start: Vertex
    end: Vertex
    label: EdgeState = EdgeState.FREE

class Graph:
    """Граф"""
    def __init__(self, *edges: Edge):
        self._edges = list(edges)
    
    def edges(self) -> typing.List[Edge]:
        """Список рёбер"""
        return self._edges

    def vertexes(self) -> typing.List[Vertex]:
        """Список вершин"""
        set_vertex = set()
        for edge in self._edges:
            set_vertex.add(edge.start)
            set_vertex.add(edge.end)
        return list(set_vertex)

