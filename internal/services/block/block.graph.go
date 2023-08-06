package block

import (
	"fmt"
	"log"
)

type JoinGraph struct {
	Vertices map[string]*Vertex
}

type Vertex struct {
	Val   *BlockData
	Edges map[string]*Edge
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
		log.Println("NOT BUILDING EDGE because of srcKey :", srcKey)
		return
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		log.Println("NOT BUILDING EDGE because of destKey :", destKey)
		return
	}
	graph.Vertices[srcKey].Edges[destKey] = &Edge{Vertex: graph.Vertices[destKey]}
}

func (graph *JoinGraph) Neighbors(srcKey string) []BlockData {
	result := []BlockData{}

	for _, edge := range graph.Vertices[srcKey].Edges {
		result = append(result, *edge.Vertex.Val)
	}
	return result
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
				fmt.Printf("      %s -> %s\n", edgeName, edge.Vertex.Val.Name)
			}
		}
		fmt.Println()
	}
}
