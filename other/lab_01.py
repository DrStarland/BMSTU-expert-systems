import models
from typing import List, Deque, Optional
from collections import deque 
import utils



class AlgorithmDFS:
    def __init__(self, graph: models.Graph):
        self._graph: models.Graph = graph


    def search(
        self,
        source: models.Vertex,
        target: models.Vertex
    ) ->  Optional[List[models.Vertex]]:
        stack: Deque[models.Vertex] = deque()
        open_vertexes: List[models.Vertex] = self._graph.vertexes()
        close_vertexes: List[models.Vertex] = list()
        
        stack.append(source)
        while stack:
            source = stack[-1] # stack.peek()
            edge = self._search_edge(source)
            if edge:
                edge.label = models.EdgeState.TRAVERSED
                stack.append(edge.end)

                if edge.end == target:
                    return list(stack)
            else:
                stack.pop()
                open_vertexes.remove(source)
                close_vertexes.append(source)

        return None

    def _search_edge(
        self,
        source: models.Graph
    ) -> Optional[models.Edge]:
        for edge in self._graph.edges():
            if edge.start == source \
                and edge.label == models.EdgeState.FREE:
                return edge
        return None

from models import *

graph = Graph(
    Edge(Vertex(1), Vertex(2)),
    Edge(Vertex(2), Vertex(3)),
    Edge(Vertex(3), Vertex(4)),

    Edge(Vertex(3), Vertex(6)),
    #Edge(Vertex(6), Vertex(3)),

    Edge(Vertex(5), Vertex(6)),
    Edge(Vertex(6), Vertex(7)),
    Edge(Vertex(6), Vertex(8)),
)

alg = AlgorithmDFS(graph)
source = Vertex(1)
target = Vertex(8)

path = alg.search(source, target)

utils.print_path(path)
utils.show_graph(graph, path, node_color_map={
    source.number: 'red',
    target.number: 'red',
})