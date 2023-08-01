package tests

import (
	blocks "dashboard/internal/services/block"
	"testing"
)

func TestBlock(t *testing.T) {
	t.Run("Get blocks", func(t *testing.T) {
		b, err := blocks.ReadBlockFile("./schema/Stories.yaml")
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if b.Blocks[0].Name != "Stories" {
			t.Errorf("Err -> \nWant %q\nGot %q", "Stories", b.Blocks[0].Name)
		}
		if b.Blocks[0].Table != "Stories" {
			t.Errorf("Err -> \nWant %q\nGot %q", "Stories", b.Blocks[0].Table)
		}
		if b.Blocks[0].Measures[0].Name != "count" {
			t.Errorf("Err -> \nWant %q\nGot %q", "count", b.Blocks[0].Measures[0].Name)
		}
		if b.Blocks[0].Dimensions[0].Name != "category" {
			t.Errorf("Err -> \nWant %q\nGot %q", "category", b.Blocks[0].Dimensions[0].Name)
		}
	})

	t.Run("Get instance", func(t *testing.T) {
		instance := blocks.GetInstance()
		if instance == nil {
			t.Fatalf("Instance is nil")
		}
		if instance.Blocks == nil {
			t.Fatalf("Blocks is nil")
		}
		if instance.Blocks[0].Blocks[0].Dimensions == nil {
			t.Fatalf("Dimensions is nil")
		}
	})

	t.Run("Get block from name", func(t *testing.T) {
		block := blocks.GetBlockFromName("Stories")
		if block.Name != "Stories" {
			t.Fatalf("Block does not have same name")
		}
	})
}
