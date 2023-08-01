package sqlStages

import "dashboard/internal/services/block"

type MeasureTypeFunc func(string) string

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql string) string {
	return "count(" + sql + ")"
}

func GenerateMeasureSelect(measure string, blockData *block.BlockData) string {
	for _, m := range blockData.Measures {
		if m.Name == measure {
			return MeasureTypes[m.Type](m.Sql)
		}
	}
	return ""
}

func GenerateDimensionSelect(dimension string, blockData *block.BlockData) string {
	for _, d := range blockData.Dimensions {
		if d.Name == dimension {
			return d.Sql
		}
	}
	return ""
}
