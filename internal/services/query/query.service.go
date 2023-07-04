package query

import (
	blockService "dashboard/internal/services/block"
	"dashboard/pkg/utils"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
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

func getStringsWithBlockName(blockName string, arr *[]string) []string {
	res := make([]string, 0)
	for i, v := range *arr {
		if strings.HasPrefix(v, blockName) {
			res = append(res, v)
			*arr = utils.Remove(*arr, i)
		}
	}
	return res
}

func GetBlockQueriesFromQuery(query Query) []BlockQuery {
	blockQueries := make([]BlockQuery, 0)

	for len(query.Dimensions) > 0 {
		blockName := strings.Split(query.Dimensions[0], ".")[0]
		dimensionInQuery := getStringsWithBlockName(blockName, &query.Dimensions)
		measuresInQuery := getStringsWithBlockName(blockName, &query.Measures)
		blockQueries = append(blockQueries, buildBlockQuery(dimensionInQuery, measuresInQuery, blockName))
	}
	return blockQueries
}

func BuildGroupStageForMeasures(query Query, join *blockService.Join) bson.M {
	blockQueries := GetBlockQueriesFromQuery(query)
	d := bson.M{}
	for _, blockQuery := range blockQueries {
		for _, dimension := range blockQuery.Dimensions {
			if blockQuery.Name != join.Name {
				d[dimension] = "$_id." + dimension
			} else {
				d[dimension] = "$" + join.Name + "." + dimension
			}
		}
	}
	//TODO: only handling count measure for now, need to find out how to add multiple measures to the group stage in mongoDB
	return bson.M{"$group": bson.M{"_id": d, "count": bson.M{"$sum": 1}}}
}

func FindBlockWithJoin(dimensions []string) *blockService.Join {
	for i, dimension := range dimensions {
		block := blockService.GetBlockFromName(getBlockName(dimension))
		for _, join := range block.Joins {
			if hasBlockName(dimensions[i+1:], join.Name) {
				return &join
			}
		}
	}
	return nil
}

func GenerateGroupStage(dimensions []string, join *blockService.Join) bson.M {
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
	return bson.M{"$group": bson.M{"_id": groupStage}}
}

func ParseQuery(query Query) ([]bson.M, error) {
	var stages []bson.M

	if len(query.Filters) > 0 {
		filterStages, err := BuildAStage[Filter](query.Filters, BuildAllFilters)
		if err != nil {
			return []bson.M{}, err
		}
		stages = append(stages, filterStages...)
	}

	if len(query.TimeDimensions) > 0 {
		timeDimensionStage, err := BuildAStage[TimeDimension](query.TimeDimensions, BuildAllTimeDimensions)
		if err != nil {
			return []bson.M{}, err
		}
		stages = append(stages, timeDimensionStage...)
	}

	join := FindBlockWithJoin(query.Dimensions)
	if join != nil {
		lookupStage := BuildLookupStage(*join)
		stages = append(stages, lookupStage)
		stages = append(stages, bson.M{"$unwind": "$" + join.Name})
	}

	groupStage := GenerateGroupStage(query.Dimensions, join)
	stages = append(stages, groupStage)
	stages = append(stages, generateProjectStage(query.Dimensions, query.Measures))
	stages = append(stages, generateOffsetStage(query.Offset))
	stages = append(stages, generateLimitStage(query.Limit))
	if len(query.Order.DimensionName) > 0 {
		stages = append(stages, generateOrderStage(query.Order))
	}
	log.Println(stages)
	documents := executeStages(stages, "Stories")
	return documents, nil
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
