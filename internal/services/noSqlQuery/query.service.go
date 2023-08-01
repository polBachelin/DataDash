package noSqlQuery

import (
	"dashboard/internal/services/query"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

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

func ParseQuery(q query.Query) ([]bson.M, error) {
	var stages []bson.M

	if len(q.Filters) > 0 {
		filterStages, err := BuildAStage[query.Filter](q.Filters, BuildAllFilters)
		if err != nil {
			return []bson.M{}, err
		}
		stages = append(stages, filterStages...)
	}

	if len(q.TimeDimensions) > 0 {
		timeDimensionStage, err := BuildAStage[query.TimeDimension](q.TimeDimensions, BuildAllTimeDimensions)
		if err != nil {
			return []bson.M{}, err
		}
		stages = append(stages, timeDimensionStage...)
	}

	join := query.FindBlockWithJoin(q.Dimensions)
	if join != nil {
		lookupStage := BuildLookupStage(*join)
		stages = append(stages, lookupStage)
		stages = append(stages, bson.M{"$unwind": "$" + join.Name})
	}

	groupStage := GenerateGroupStage(q.Dimensions, q.Measures, join)
	stages = append(stages, groupStage)
	stages = append(stages, generateProjectStage(q.Dimensions, q.Measures))
	stages = append(stages, generateOffsetStage(q.Offset))
	stages = append(stages, generateLimitStage(q.Limit))
	if len(q.Order.DimensionName) > 0 {
		stages = append(stages, generateOrderStage(q.Order))
	}
	log.Println(stages)
	documents := executeStages(stages, FindCollectionName(q.Dimensions, join))
	return documents, nil
}
