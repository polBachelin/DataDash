package block

import (
	"dashboard/pkg/utils"
	"log"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGraph(t *testing.T) {
	path := utils.GetEnvVar("SCHEMA_PATH", "./example_schema/sale_db")
	all, err := ReadAllBlocks(path)
	if err != nil {
		t.Fatalf("Failed to read schemas: %v", err)
	}
	graph := NewGraph(all)
	t.Run("Test graph build", func(t *testing.T) {
		saleNeighbors := graph.Neighbors("Sale")
		if slices.ContainsFunc(saleNeighbors, func(data BlockData) bool { return data.Name == "Product" }) != true {
			t.Fatalf("Product is not a neighbor of Sale")
		}
		graph.printGraph()
	})
	t.Run("Test graph FindJoinPath", func(t *testing.T) {
		joinPath, found := graph.FindJoinPath(graph.Vertices["Sale"], "Status_name")
		log.Println(joinPath)
		if found == false {
			t.Fatalf("Join Path not found: %v", joinPath)
		}
	})
}
