package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func BuildTimeDimension(timeDimension TimeDimension) (bson.M, error) {
	member := strings.Split(timeDimension.Dimension, ".")[1]
	return bson.M{"$match": bson.M{member: bson.M{"$gte": timeDimension.DateRange[0], "$lte": timeDimension.DateRange[1]}}}, nil
}

// TODO check if I should enable the fact to add time dimensions for a join
// TODO not using granularity need to see how I can take advantage of that
func BuildAllTimeDimensions(timeDimensions []TimeDimension) ([]bson.M, error) {
	timeDimensionStages := make([]bson.M, len(timeDimensions))

	for _, d := range timeDimensions {
		stage, _ := BuildTimeDimension(d)
		timeDimensionStages = append(timeDimensionStages, stage)
	}
	return timeDimensionStages, nil
}
