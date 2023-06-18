package query

import (
	blockService "dashboard/internal/services/block"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type MeasureTypeFunc func(sql string, dimension blockService.Dimensions) bson.M

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql string, dimension blockService.Dimensions) bson.M {
	stage := bson.M{"$group": bson.M{"_id": "$" + dimension.Sql, "count": bson.M{"$sum": 1}}}
	return stage
}

func handleMeasure(block blockService.BlockData, measureName string) bson.M {
	measureIndex := slices.IndexFunc(block.Measures, func(data blockService.Measures) bool { return data.Name == measureName })
	if measureIndex == -1 {
		return nil
	}
	measureFunc := MeasureTypes[block.Measures[measureIndex].Type]
	return measureFunc(block.Measures[measureIndex].Sql, block.Dimensions[0])
}
