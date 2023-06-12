package query

import (
	blockService "dashboard/internal/services/block"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

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
