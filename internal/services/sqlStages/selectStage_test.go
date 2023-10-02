package sqlStages

import (
	"dashboard/internal/services/block"
	"strings"
	"testing"
)

func TestSelectStage(t *testing.T) {
	selectStage, _ := GenerateMeasureSql("count", block.GetBlockFromName("BusinessUnits"))
	if !strings.Contains(selectStage, "count(BusinessUnits.id)") {
		t.Fatalf("Select stage is wrong : %v", selectStage)
	}
}

func TestSelectStageNumber(t *testing.T) {
	selectStage, _ := GenerateMeasureSql("revenue", block.GetBlockFromName("Projects"))
	if !strings.Contains(selectStage, "sum(amount_sold - allocated_budget)") {
		t.Fatalf("Select stage is wrong : %v", selectStage)
	}
}
