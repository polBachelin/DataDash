package query

import (
	blockService "dashboard/internal/services/block"
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
}

type MeasureTypeFunc func(sql string) bson.D

var MeasureTypes = map[string]MeasureTypeFunc{
	"count": MeasureCount,
}

func MeasureCount(sql string) bson.D {
	stage := bson.D{{Key: "$group",
		Value: bson.D{{Key: "_id", Value: sql}, {Key: "count", Value: bson.D{{Key: "$count", Value: "count"}}}},
	}}
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

func handleMeasure(block blockService.BlockData, measureName string, collectionName string) bson.D {
	measureIndex := slices.IndexFunc(block.Measures, func(data blockService.Measures) bool { return data.Name == measureName })
	if measureIndex == -1 {
		return nil
	}
	measureFunc := MeasureTypes[block.Measures[measureIndex].Type]
	return measureFunc(block.Measures[measureIndex].Sql)
}

func ParseQuery(query Query) {
	for _, measure := range query.Measures {
		n := strings.Split(measure, ".")
		block := blockService.GetBlockFromName(n[0])
		if block != nil {
			handleMeasure(*block, n[1], n[0])
		}

	}
}
