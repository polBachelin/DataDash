package block

import (
	"fmt"
	"log"
	"math"
)

type JoinGraph struct {
	Vertices map[string]*Vertex
}

type Vertex struct {
	Val   *BlockData
	Edges map[string]*Edge
	Prev  *Vertex
}

type Edge struct {
	Vertex *Vertex
}

func NewGraph(files []*FileData) *JoinGraph {
	g := &JoinGraph{Vertices: map[string]*Vertex{}}

	for _, f := range files {
		for _, block := range f.Blocks {
			if _, ok := g.Vertices[block.Name]; !ok {
				g.AddVertex(block.Name, &block)
			}
			for _, edge := range block.Joins {
				if _, ok := g.Vertices[edge.Name]; !ok {
					edgeBlock := GetBlockFromName(edge.Name)
					g.AddVertex(edge.Name, edgeBlock)
				}
				g.AddEdge(block.Name, edge.Name)
			}
		}
	}
	return g
}

func (graph *JoinGraph) AddVertex(key string, val *BlockData) {
	graph.Vertices[key] = &Vertex{Val: val, Edges: map[string]*Edge{}}
}

func (graph *JoinGraph) AddEdge(srcKey, destKey string) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		log.Println("NOT BUILDING EDGE because srcKey does not exist :", srcKey)
		return
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		log.Println("NOT BUILDING EDGE because destKey does not exist :", destKey)
		return
	}
	graph.Vertices[srcKey].Edges[destKey] = &Edge{Vertex: graph.Vertices[destKey]}
}

func (graph *JoinGraph) Neighbors(srcKey string) []BlockData {
	result := []BlockData{}

	for _, edge := range graph.Vertices[srcKey].Edges {
		if edge != nil && edge.Vertex != nil && edge.Vertex.Val != nil {
			result = append(result, *edge.Vertex.Val)
		}
	}
	return result
}

func (graph *JoinGraph) FindJoinPath(startVertex *Vertex, targetVertexName string) ([]string, bool) {
	visited := make(map[string]bool)
	parentVertex := startVertex
	path, found := graph.DfsWithPath(parentVertex, targetVertexName, visited)
	for !found {
		parentVertex = graph.FindVertexWithEdge(parentVertex.Val.Name)
		if parentVertex == nil {
			log.Println("No vertex found")
			return nil, false
		}
		var p []string
		p, found = graph.DfsWithPath(parentVertex, targetVertexName, make(map[string]bool))
		path = append(path, p...)
	}
	return path, found
}

func (graph *JoinGraph) FindVertexWithEdge(targetVertexName string) *Vertex {
	for _, vertex := range graph.Vertices {
		if _, found := vertex.Edges[targetVertexName]; found {
			return vertex
		}
	}
	return nil
}

func (graph *JoinGraph) DfsWithPath(currentVertex *Vertex, targetVertexName string, visited map[string]bool) ([]string, bool) {
	visited[currentVertex.Val.Name] = true

	if currentVertex.Val.Name == targetVertexName {
		log.Printf("Found target vertex : %s", currentVertex.Val.Name)
		return []string{currentVertex.Val.Name}, true
	}

	for edgeName, edge := range currentVertex.Edges {
		if !visited[edgeName] {
			if path, found := graph.DfsWithPath(edge.Vertex, targetVertexName, visited); found {
				return append(path, currentVertex.Val.Name), true
			}
		}
	}
	return nil, false
}

func (graph *JoinGraph) ShortestPath(startName string, targetName string) []string {
	distances := graph.Dijkstra(startName)

	// Check if there is a path from startVertex to endVertex
	if distances[targetName] == math.Inf(1) {
		log.Printf("No path from startVertex")
		return nil // No path found
	}

	path := []string{targetName}
	currentVertexName := targetName
	for currentVertexName != startName {
		currentVertex := graph.Vertices[currentVertexName]
		prevVertex := currentVertex.Prev
		path = append([]string{prevVertex.Val.Name}, path...)
		currentVertexName = prevVertex.Val.Name
	}

	return path
}

func (graph *JoinGraph) Dijkstra(startVertexName string) map[string]float64 {
	// Initialize distances and visited map
	distances := make(map[string]float64)
	visited := make(map[string]bool)

	// Initialize distances to infinity and visited to false
	for vertexName := range graph.Vertices {
		distances[vertexName] = math.Inf(1)
		visited[vertexName] = false
	}

	// Set the distance of the start vertex to 0
	distances[startVertexName] = 0

	// Loop through all vertices
	for len(visited) < len(graph.Vertices) {
		// Find the vertex with the smallest distance
		minDist := math.Inf(1)
		currentVertexName := ""
		for vertexName, dist := range distances {
			if !visited[vertexName] && dist < minDist {
				minDist = dist
				currentVertexName = vertexName
			}
		}

		// Mark the current vertex as visited
		visited[currentVertexName] = true

		// Update distances to adjacent vertices
		currentVertex := graph.Vertices[currentVertexName]
		for edgeName, edge := range currentVertex.Edges {
			if !visited[edgeName] {
				newDist := distances[currentVertexName] + 1 // You can replace this with your own weight function
				if newDist < distances[edgeName] {
					distances[edgeName] = newDist
					edge.Vertex.Prev = currentVertex
				}
			}
		}
	}

	return distances
}

func (graph *JoinGraph) printGraph() {
	for vertexName, vertex := range graph.Vertices {
		fmt.Printf("Vertex: %s\n", vertexName)
		fmt.Printf("   Block: %v\n", vertex.Val)
		if len(vertex.Edges) == 0 {
			fmt.Println("   Edges: []")
		} else {
			fmt.Println("   Edges:")
			for edgeName, edge := range vertex.Edges {
				if edge == nil || edge.Vertex == nil || edge.Vertex.Val == nil {
					fmt.Printf(" %s -> nil\n", edgeName)
				} else {
					fmt.Printf("      %s -> %s\n", edgeName, edge.Vertex.Val.Name)
				}
			}
		}
		fmt.Println()
	}
}
