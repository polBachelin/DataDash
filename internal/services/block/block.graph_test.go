package block

import (
	"dashboard/pkg/utils"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGraph(t *testing.T) {
	t.Run("Test graph build", func(t *testing.T) {
		path := utils.GetEnvVar("SCHEMA_PATH", "./example_schema/sale_db")
		all, err := ReadAllBlocks(path)
		if err != nil {
			t.Fatalf("Failed to read schemas: %v", err)
		}
		graph := NewGraph(all)
		saleNeighbors := graph.Neighbors("Sale")
		if slices.ContainsFunc(saleNeighbors, func(data BlockData) bool { return data.Name == "Product" }) != true {
			t.Fatalf("Product is not a neighbor of Sale")
		}
	})
}
