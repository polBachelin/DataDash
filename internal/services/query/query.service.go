package query

import (
	"context"
	"dashboard/internal/database"
	blockService "dashboard/internal/services/block"
	"fmt"
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
	stage := bson.M{"$group": bson.M{"_id": "$" + dimension.Sql, "count": bson.M{"$sum": 1}}}
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

func handleCleanup(block blockService.BlockData) []bson.M {
	return []bson.M{}
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

type BlockQuery struct {
	Measure    []string
	Dimensions []string
	Name       string
}

func buildBlockQuery(dimensions []string, measure []string) BlockQuery {
	n := strings.Split(measure, ".")
	blockQuery := BlockQuery{Name: n[0], Measure: n[1], Dimensions: []string{}}
	for _, dim := range dimensions {
		t := strings.Split(dim, ".")
		if t[0] == n[0] {
			blockQuery.Dimensions = append(blockQuery.Dimensions, t[1])
		}
	}
	return blockQuery
}

func getStringsWithBlockName(blockName string, arr []string) []string {
	res := make([]string, 0)
	for _, v := range arr {
		if strings.HasPrefix(v, blockName) {
			res = append(res, v)
		}
	}
	return res
}

func ParseQuery(query Query) QueryResult {
	var res QueryResult
	blockQueries := make([]BlockQuery, 0)
	wholeLength := len(query.Dimensions) + len(query.Measures)
	currentLength := 0
	i := 0

	for currentLength != wholeLength {
		blockName := strings.Split(query.Dimensions[i], ".")[0]
		dimensionInQuery := getStringsWithBlockName(blockName, query.Dimensions)
		measuresInQuery := getStringsWithBlockName(blockName, query.Measures)
		blockQueries = append(blockQueries, buildBlockQuery(dimensionInQuery, measuresInQuery))
		currentLength += len(dimensionInQuery) + len(measuresInQuery)
		i++
	}
	// block := blockService.GetBlockFromName(n[0])
	// if block != nil {
	// 	measureStage := handleMeasure(*block, n[1])
	// 	documents := executeStage(measureStage, n[0])
	// 	resData := buildResData(documents, n[0], n[1])
	// 	res.Data = append(res.Data, resData...)
	// }
	return res
}

// Name needs to contain [CUBE_NAME, MEASURE_NAME]
func buildResData(documents []bson.M, blockName string, measureName string) []ResultData {
	resData := make([]ResultData, 0)
	var data ResultData

	for _, doc := range documents {
		data.Name = blockName
		data.Dimension = fmt.Sprintf("%v", doc["_id"])
		data.Measure = fmt.Sprintf("%v", doc[measureName])
		resData = append(resData, data)
		log.Println(doc)
		log.Println(doc["count"])
	}
	return resData
}
