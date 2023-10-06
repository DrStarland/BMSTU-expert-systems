import models
from typing import List, Deque, Optional
from collections import deque 
import utils



class AlgorithmDFS:
    def __init__(self, graph: models.Graph):
        self._graph: models.Graph = graph


    def search(
        self,
        source: models.Node,
        target: models.Node
    ) ->  Optional[List[models.Node]]:
        stack: Deque[models.Node] = deque()
        open_vertexes: List[models.Node] = self._graph.vertexes()
        close_vertexes: List[models.Node] = list()
        
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
    Edge(Node(1), Node(2)),
    Edge(Node(2), Node(3)),
    Edge(Node(3), Node(4)),

    Edge(Node(3), Node(6)),
    #Edge(Node(6), Node(3)),

    Edge(Node(5), Node(6)),
    Edge(Node(6), Node(7)),
    Edge(Node(6), Node(8)),
)

alg = AlgorithmDFS(graph)
source = Node(1)
target = Node(8)

path = alg.search(source, target)

utils.print_path(path)
utils.show_graph(graph, path, node_color_map={
    source.number: 'red',
    target.number: 'red',
})