package sqlStages

import (
	"dashboard/internal/services/block"
	"fmt"
)

type MeasureTypeFunc func(string, string, string) (string, error)

var MeasureTypes = map[string]MeasureTypeFunc{
	"count":  MeasureCount,
	"number": MeasureNumber,
	"sum":    MeasureSum,
}

func MeasureCount(sql, tableName, mesName string) (string, error) {
	return fmt.Sprintf("count(%v.%v) as \"%v.%v\"", tableName, sql, tableName, mesName), nil
}

func MeasureSum(sql, tableName, mesName string) (string, error) {
	return fmt.Sprintf("sum(%v.%v) as \"%v.%v\"", tableName, sql, tableName, mesName), nil
}

func MeasureNumber(sql, tableName, mesName string) (string, error) {
	return fmt.Sprintf("%v", sql), nil
}

func GenerateMeasureSql(measure string, blockData *block.BlockData) (string, error) {
	m, err := block.GetMeasureFromBlock(blockData, measure)
	if err != nil {
		return "", fmt.Errorf("getting measure block failed")
	}
	return MeasureTypes[m.Type](m.Sql, blockData.Name, m.Name)
}

func GenerateDimensionSelect(dimension string, blockData *block.BlockData) (string, error) {

	for _, d := range blockData.Dimensions {
		if d.Name == dimension {
			return fmt.Sprintf("%v.%v as \"%s.%s\"", blockData.Name, d.Sql, blockData.Name, d.Name), nil
		}
	}
	return "", fmt.Errorf("no dimension found for block %s with name %s", blockData.Name, dimension)
}
