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
			join, err := blockService.GetBlockJoinFromName(n[0], *lastBlock)
			if err != nil {
				return bson.M{}, errors.New("could not find dimension and join from block :" + block.Name)
			}
			ids[join.Name] = "$" + join.LocalField
		}
		lastName = n[0]
	}
	return bson.M{"$group": ids}, nil
}

func ParseQuery(query Query) (QueryResult, error) {
	var res QueryResult
	var stages []bson.M

	blockQueries := GetBlockQueriesFromQuery(query)
	err := checkJoinFromQueries(blockQueries)
	if err != nil {
		return QueryResult{}, err
	}
	groupStage, err := BuildGroupStageFromDimensions(query.Dimensions)
	if err != nil {
		return QueryResult{}, err
	}
	stages = append(stages, groupStage)
	if len(query.Filters) > 0 {
		filterStages, err := BuildAllFilters(query.Filters)
		if err != nil {
			return QueryResult{}, err
		}
		stages = append(stages, filterStages...)
	}
	if len(query.TimeDimensions) > 0 {
		timeDimensionStage, err := BuildAllTimeDimensions(query.TimeDimensions)
		if err != nil {
			return QueryResult{}, err
		}
		stages = append(stages, timeDimensionStage...)
	}
	return res, nil
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
