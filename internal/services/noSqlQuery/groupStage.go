package noSqlQuery

import (
	blockService "dashboard/internal/services/block"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type MeasureTypeFunc func() bson.M

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount() bson.M {
	return bson.M{"$sum": 1}
}

func BuildGroupStage(block blockService.BlockData, joinChildIndex int, blockQuery BlockQuery) bson.M {
	ids := make(bson.M)

	log.Println(block.Dimensions)
	for _, dimension := range blockQuery.Dimensions {
		dimensionIndex := slices.IndexFunc(block.Dimensions, func(data blockService.Dimensions) bool { return data.Name == dimension })
		log.Println("Dimension : " + dimension)
		if dimensionIndex != -1 {
			ids[dimension] = "$" + block.Dimensions[dimensionIndex].Sql
		} else {
			ids[block.Joins[joinChildIndex].Name] = "$" + block.Joins[joinChildIndex].LocalField
		}
	}
	return bson.M{"$group": ids}
}

func AddMeasureToGroupStage(measures []string) bson.M {
	//TODO: when measures have been figured out, add a loop here
	measureFunc := MeasureTypes[getMemberName(measures[0])]
	return measureFunc()
}

func GenerateGroupStage(dimensions, measures []string, join *blockService.Join) bson.M {
	groupStage := bson.M{}
	for _, dimension := range dimensions {
		memberName := getMemberName(dimension)
		blockName := getBlockName(dimension)
		if join != nil && blockName == join.Name {
			groupStage[memberName] = "$" + join.Name + "." + memberName
		} else {
			groupStage[memberName] = "$" + memberName
		}
	}
	measureStage := AddMeasureToGroupStage(measures)
	return bson.M{"$group": bson.M{"_id": groupStage, getMemberName(measures[0]): measureStage}}
}
