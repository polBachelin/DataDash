package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
	"log"
)

type MeasureTypeFunc func(string, string) string

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql, tableName string) string {
	return fmt.Sprintf("count(%v.%v)", tableName, sql)
}

func GenerateMeasureSelect(measure string, blockData *block.BlockData) string {
	for _, m := range blockData.Measures {
		if m.Name == measure {
			return MeasureTypes[m.Type](m.Sql, blockData.Name)
		}
	}
	return ""
}

func GenerateDimensionSelect(dimension string, blockData *block.BlockData) string {
	for _, d := range blockData.Dimensions {
		log.Println("dimension:", d.Name)
		log.Println("dimension:", dimension)
		if d.Name == dimension {
			return fmt.Sprintf("%v.%v", blockData.Name, d.Sql)
		}
	}
	return ""
}
