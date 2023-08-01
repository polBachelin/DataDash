package block

type JoinGraph struct {
	Vertices map[string]*Vertex
}

type Vertex struct {
	Val   *Join
	Edges map[string]*Edge
}

type Edge struct {
	Vertex *Vertex
}

func NewGraph(files FileData) *JoinGraph {
	g := &JoinGraph{Vertices: map[string]*Vertex{}}

	for _, block := range files.Blocks {
		if len(block.Joins) == 0 {
			continue
		}
		for _, join := range block.Joins {
			g.AddVertex(block.Name)

		}
	}
}

func (graph *JoinGraph) AddVertex(key string, val *Join) {
	graph.Vertices[key] = &Vertex{Val: val, Edges: map[string]*Edge{}}
}

func (graph *JoinGraph) AddEdge(srcKey, destKey string) {
	if _, ok := graph.Vertices[srcKey]; !ok {
		return
	}
	if _, ok := graph.Vertices[destKey]; !ok {
		return
	}
	graph.Vertices[srcKey].Edges[destKey] = &Edge{Vertex: graph.Vertices[destKey]}
}

func (graph *JoinGraph) Neighbors(srcKey string) []Join {
	result := []Join{}

	for _, edge := range graph.Vertices[srcKey].Edges {
		result = append(result, *edge.Vertex.Val)
	}
	return result
}
