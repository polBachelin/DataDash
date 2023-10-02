package block

import (
	"dashboard/pkg/utils"
	"log"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGraph(t *testing.T) {
	path := utils.GetEnvVar("SCHEMA_PATH", "./schema")
	all, err := ReadAllBlocks(path)
	if err != nil {
		t.Fatalf("Failed to read schemas: %v", err)
	}
	graph := NewGraph(all)
	graph.printGraph()
	t.Run("Test graph build", func(t *testing.T) {
		saleNeighbors := graph.Neighbors("Departments")
		if slices.ContainsFunc(saleNeighbors, func(data BlockData) bool { return data.Name == "BusinessUnits" }) != true {
			t.Fatalf("Product is not a neighbor of Sale")
		}
		graph.printGraph()
	})
	t.Run("Test graph FindJoinPath", func(t *testing.T) {
		joinPath, found := graph.FindJoinPath(graph.Vertices["Projects"], "BusinessUnits")
		log.Println(joinPath)
		resPath := []string{"BusinessUnits", "Departments", "Projects"}
		if found == false || !slices.Equal(joinPath, resPath) {
			t.Fatalf("Join Path not found: %v", joinPath)
		}
	})
	t.Run("Test graph Other key", func(t *testing.T) {
		joinPath, found := graph.FindJoinPath(graph.Vertices["Departments"], "BusinessUnits")
		log.Println(joinPath)
		resPath := []string{"Departments", "BusinessUnits"}
		if found == false || !slices.Equal(joinPath, resPath) {
			t.Fatalf("Join Path not found: %v", joinPath)
		}
	})

	t.Run("Shortest path", func(t *testing.T) {
		path := graph.ShortestPath("Departments", "BusinessUnits")
		if path == nil || !slices.Equal(path, []string{"Departments", "BusinessUnits"}) {
			t.Fatalf("Incorrect path: %v", path)
		}
	})

	t.Run("Other side shortest path", func(t *testing.T) {
		path := graph.ShortestPath("BusinessUnits", "Departments")
		log.Println(path)
		if path == nil || !slices.Equal(path, []string{"BusinessUnits", "Departments"}) {
			t.Fatalf("Incorrect path: %v", path)
		}
	})

}
