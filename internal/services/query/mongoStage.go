package query

import (
	"context"
	"dashboard/internal/database"
	blockService "dashboard/internal/services/block"
	"log"

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

func handleMeasure(block blockService.BlockData, measureName string) bson.M {
	measureIndex := slices.IndexFunc(block.Measures, func(data blockService.Measures) bool { return data.Name == measureName })
	if measureIndex == -1 {
		return nil
	}
	measureFunc := MeasureTypes[block.Measures[measureIndex].Type]
	return measureFunc(block.Measures[measureIndex].Sql, block.Dimensions[0])
}

func executeStage(stage bson.M, collectionName string) []bson.M {
	collection := database.GetCollection(collectionName)
	res, err := collection.Aggregate(context.TODO(), []bson.M{stage})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	document := []bson.M{}
	err = res.All(context.TODO(), &document)
	if err != nil {
		return nil
	}
	return document
}
