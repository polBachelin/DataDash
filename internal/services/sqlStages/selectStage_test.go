package sqlStages

import (
	"dashboard/internal/services/block"
	"strings"
	"testing"
)

func TestSelectStage(t *testing.T) {
	selectStage := GenerateMeasureSelect("count", block.GetBlockFromName("Sale"))
	if !strings.Contains(selectStage, "count(Sale.sale_id)") {
		t.Fatalf("Select stage is wrong : %v", selectStage)
	}
}
