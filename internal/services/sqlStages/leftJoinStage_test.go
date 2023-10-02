package sqlStages

import (
	"dashboard/internal/services/block"
	"dashboard/pkg/utils"
	"testing"
)

func TestLeftJoin(t *testing.T) {
	path := []string{"Departments", "BusinessUnits"}
	all, err := block.ReadAllBlocks(utils.GetEnvVar("SCHEMA_PATH", "./schema"))
	if err != nil {
		t.Fatalf("Failed to read schema: %v", err)
	}
	graph := block.NewGraph(all)
	t.Run("Check two level join", func(t *testing.T) {
		clause := GenerateJoinClause(path, graph)
		if clause != " LEFT JOIN departments as Departments ON Departments.business_unit_id = BusinessUnits.id" {
			t.Fatalf("Wrong clause generated")
		}
	})

}
