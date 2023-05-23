package query

import (
	"context"
	"dashboard/internal/database"
	blockService "dashboard/internal/services/block"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type Query struct {
	Measures       []string        `json:"measures"`
	Dimensions     []string        `json:"dimensions"`
	Filters        []Filter        `json:"filters"`
	TimeDimensions []TimeDimension `json:"time_dimensions"`
	Limit          int             `json:"limit"`
	Offset         int             `json:"offset"`
	Order          Order           `json:"order"`
}

type Filter struct {
	Member   string   `json:"member"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type TimeDimension struct {
	Dimension   string   `json:"dimension"`
	DateRange   []string `json:"date_range"`
	Granularity string   `json:"granularity"`
}

type Order struct {
	DimensionName  []string `json:"dimension_name"`
	DimensionOrder []string `json:"dimension_order"`
	MeasureName    []string `json:"measure_name"`
	MeasureOrder   []string `json:"measure_order"`
}

type QueryResult struct {
	Data []ResultData `json:"data"`
}

type ResultData struct {
	Name          string `json:"name"`
	MeasureType   string `json:"type"`
	Measure       string `json:"result"`
	Dimension     string `json:"dimension"`
	DimensionType string `json:"dimension_type"`
}

type MeasureTypeFunc func(sql string, dimension blockService.Dimensions) bson.M

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql string, dimension blockService.Dimensions) bson.M {
	stage := bson.M{}
	// stage := bson.D{{Key: "$group",
	// 	Value: bson.D{{Key: "_id", Value: "$" + dimension.Sql}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}},
	// }}
	log.Println(stage)
	return stage
}

// Name is always under the format CUBE_NAME.MEMBER_NAME
func checkName(name string) bool {
	s := strings.Split(name, ".")
	blockInstance := blockService.GetInstance()
	for _, fd := range blockInstance.Blocks {
		for _, block := range fd.Blocks {
			if block.Name == s[0] {
				return true
			}
		}
	}
	return false
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

func ParseQuery(query Query) []QueryResult {
	//res := make([]QueryResult, 0)

	for _, measure := range query.Measures {
		//// Retrieving measure name from CUBE_MEMBER.MEMBER_NAME convention
		n := strings.Split(measure, ".")
		block := blockService.GetBlockFromName(n[0])
		if block != nil {
			measureStage := handleMeasure(*block, n[1])
			documents := executeStage(measureStage, n[0])
			log.Println(documents)
		}
	}
	return []QueryResult{}
}
