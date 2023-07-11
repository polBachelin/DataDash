package query

import (
	blockService "dashboard/internal/services/block"
	"log"

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

	groupStage := GenerateGroupStage(query.Dimensions, query.Measures, join)
	stages = append(stages, groupStage)
	stages = append(stages, generateProjectStage(query.Dimensions, query.Measures))
	stages = append(stages, generateOffsetStage(query.Offset))
	stages = append(stages, generateLimitStage(query.Limit))
	if len(query.Order.DimensionName) > 0 {
		stages = append(stages, generateOrderStage(query.Order))
	}
	log.Println(stages)
	documents := executeStages(stages, FindCollectionName(query.Dimensions, join))
	return documents, nil
}
