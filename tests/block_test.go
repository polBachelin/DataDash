package tests

import (
	blocks "dashboard/internal/services/block"
	"testing"
)

func TestBlock(t *testing.T) {
	t.Run("Get blocks", func(t *testing.T) {
		b, err := blocks.ReadFile("/home/polo/Projects/cube_remake/dashboard_api/schema/Stories.yaml")
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if b.Blocks[0].Name != "Stories" {
			t.Errorf("Err -> \nWant %q\nGot %q", "Stories", b.Blocks[0].Name)
		}
		if b.Blocks[0].Sql != "SELECT * FROM stories" {
			t.Errorf("Err -> \nWant %q\nGot %q", "SELECT * FROM stories", b.Blocks[0].Sql)
		}
		if b.Blocks[0].Measures[0].Name != "count" {
			t.Errorf("Err -> \nWant %q\nGot %q", "count", b.Blocks[0].Measures[0].Name)
		}
		if b.Blocks[0].Dimensions[0].Name != "category" {
			t.Errorf("Err -> \nWant %q\nGot %q", "category", b.Blocks[0].Dimensions[0].Name)
		}

	})
}
