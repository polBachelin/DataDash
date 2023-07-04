package query

import (
	blockService "dashboard/internal/services/block"
	"dashboard/pkg/utils"
	"errors"
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

func checkJoinFromQueries(blockQueries []BlockQuery) error {
	for i, blockQuery := range blockQueries {
		block := blockService.GetBlockFromName(blockQuery.Name)
		for _, nextBlockQuery := range blockQueries[i+1:] {
			check := slices.IndexFunc(block.Joins, func(join blockService.Join) bool { return join.Name == nextBlockQuery.Name })
			if check == -1 {
				return errors.New("No join between " + block.Name + " and " + nextBlockQuery.Name)
			}
		}
	}
	return nil
}

// This function return : joinParent, joinChild, error
func findJoinParent(blockQueries []BlockQuery) (int, int, error) {
	for i, blockQuery := range blockQueries {
		block := blockService.GetBlockFromName(blockQuery.Name)
		if block == nil {
			return -1, -1, errors.New("no block file found for name " + blockQuery.Name)
		}
		joinFound := slices.IndexFunc(block.Joins, func(data blockService.Join) bool { return data.Name == blockQuery.Name })
		if joinFound != -1 {
			return i, joinFound, nil
		}
	}
	return -1, -1, errors.New("no join parent found in block queries")
}

// TODO: segfault when dimension does not exist in block
func BuildGroupStageFromDimensions(dimensions []string) (bson.M, error) {
	ids := make(bson.M)
	lastName := dimensions[0]

	for _, dimension := range dimensions {
		n := strings.Split(dimension, ".")
		lastName = strings.Split(lastName, ".")[0]
		block := blockService.GetBlockFromName(n[0])
		if block == nil {
			return bson.M{}, errors.New("no block found")
		}
		if lastName == n[0] {
			dimIndex := slices.IndexFunc(block.Dimensions, func(data blockService.Dimensions) bool { return data.Name == n[1] })
			ids[n[1]] = "$" + block.Dimensions[dimIndex].Sql
		} else {
			//if it was not found in dimensions check joins
			lastBlock := blockService.GetBlockFromName(lastName)
			j, err := blockService.GetBlockJoinFromName(n[0], *lastBlock)
			if err != nil {
				return bson.M{}, errors.New("could not find dimension and join from block :" + block.Name)
			}
			ids[j.Name] = "$" + j.LocalField
		}
		lastName = n[0]
	}
	return bson.M{"$group": bson.M{"_id": ids}}, nil
}

func BuildLookupStage(join blockService.Join) bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from":         join.Name,
			"localField":   "_id." + join.Name,
			"foreignField": join.ForeignField,
			"as":           join.Name,
		}}
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

// TODO: currently only works if there is a join in the request, need to fix that
func ParseQuery(query Query) ([]bson.M, error) {
	var stages []bson.M

	// blockQueries := GetBlockQueriesFromQuery(query)
	join := FindBlockWithJoin(query.Dimensions)
	if join != nil {
		lookupStage := BuildLookupStage(*join)
		stages = append(stages, lookupStage)
	}

	groupStage, err := BuildGroupStageFromDimensions(query.Dimensions)
	if err != nil {
		return []bson.M{}, err
	}
	filterStages, err := BuildAStage[Filter](query.Filters, BuildAllFilters)
	if err != nil {
		return []bson.M{}, err
	}
	timeDimensionStage, err := BuildAStage[TimeDimension](query.TimeDimensions, BuildAllTimeDimensions)
	if err != nil {
		return []bson.M{}, err
	}
	stages = append(stages, filterStages...)
	stages = append(stages, timeDimensionStage...)
	stages = append(stages, groupStage)
	//TODO for each measure there should be one group stage ?
	measureStage := BuildGroupStageForMeasures(query, join)
	stages = append(stages, measureStage)
	log.Println(stages)
	documents := executeStages(stages, "Stories")
	log.Println(documents)
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
