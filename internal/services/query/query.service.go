package query

import (
	blockService "dashboard/internal/services/block"
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

func getStringsWithBlockName(blockName string, arr []string) []string {
	res := make([]string, 0)
	for _, v := range arr {
		if strings.HasPrefix(v, blockName) {
			res = append(res, v)
		}
	}
	return res
}

func getBlockQueriesFromQuery(query Query) []BlockQuery {
	blockQueries := make([]BlockQuery, 0)
	wholeLength := len(query.Dimensions) + len(query.Measures)
	currentLength := 0
	i := 0

	for currentLength != wholeLength {
		blockName := strings.Split(query.Dimensions[i], ".")[0]
		dimensionInQuery := getStringsWithBlockName(blockName, query.Dimensions)
		measuresInQuery := getStringsWithBlockName(blockName, query.Measures)
		blockQueries = append(blockQueries, buildBlockQuery(dimensionInQuery, measuresInQuery, blockName))
		currentLength += len(dimensionInQuery) + len(measuresInQuery)
		i++
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

func ParseQuery(query Query) (QueryResult, error) {
	var res QueryResult

	blockQueries := getBlockQueriesFromQuery(query)
	err := checkJoinFromQueries(blockQueries)
	if err != nil {
		return QueryResult{}, err
	}
	for _, blockQuery := range blockQueries {
		block := blockService.GetBlockFromName(blockQuery.Name)
		if block == nil {
			return QueryResult{}, errors.New("No block file found for name " + blockQuery.Name)
		}
	}
	// if block != nil {
	// 	measureStage := handleMeasure(*block, n[1])
	// 	documents := executeStage(measureStage, n[0])
	// 	resData := buildResData(documents, n[0], n[1])
	// 	res.Data = append(res.Data, resData...)
	// }
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
