package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
)

type MeasureTypeFunc func(string, string) string

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql, tableName string) string {
	return fmt.Sprintf("count(%v.%v)", tableName, sql)
}

func GenerateMeasureSql(measure string, blockData *block.BlockData) string {
	m, err := block.GetMeasureFromBlock(blockData, measure)
	if err != nil {
		return ""
	}
	return MeasureTypes[m.Type](m.Sql, blockData.Name)
}

func GenerateDimensionSelect(dimension string, blockData *block.BlockData) string {
	for _, d := range blockData.Dimensions {
		if d.Name == dimension {
			return fmt.Sprintf("%v.%v", blockData.Name, d.Sql)
		}
	}
	return ""
}
